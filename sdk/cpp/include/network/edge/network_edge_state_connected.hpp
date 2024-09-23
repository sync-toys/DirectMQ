#pragma once

#include <memory>
#include <stdexcept>

#include "../../protocol/decoder.hpp"
#include "../../protocol/messages/unsubscribe.hpp"
#include "../constants.hpp"
#include "../participant.hpp"
#include "network_edge_state_manager.hpp"
#include "state_name.hpp"
#include "utilities.hpp"

namespace directmq::network::edge {
class NetworkEdgeStateConnected : public NetworkEdgeState {
   private:
    void exchangeAllNodeSubscriptions() {
        for (auto& topic : edge->globalNetwork->getAllSubscribedTopics()) {
            protocol::messages::SubscribeMessage subscription{
                .frame =
                    protocol::messages::DataFrame{
                        .ttl = edge->globalNetwork->config.hostTTL,
                        .traversed = {edge->globalNetwork->config.hostID},
                    },
                .topic = *topic,
            };

            auto result = edge->encoder->subscribe(subscription, *edge->portal);

            if (result.error != nullptr) {
                edge->setDisconnectingState(
                    "Failed to exchange subscriptions: " +
                    std::string(result.error));
                return;
            }
        }
    }

    void updateSubscriptions(
        const std::list<std::shared_ptr<const std::string>>& oldTopics,
        const protocol::messages::DataFrame& frame) {
        auto updatedTopLevelSubscriptions =
            edge->bridgedNodeSubscriptions.getOnlyTopLevelSubscribedTopics();
        auto diff = topics::getDeduplicatedOverlappingTopicsDiff(
            oldTopics, updatedTopLevelSubscriptions);

        for (auto& topic : diff.added) {
            protocol::messages::SubscribeMessage subscription{
                .frame = frame,
                .topic = *topic,
            };

            edge->globalNetwork->subscribed(subscription);
        }

        for (auto& topic : diff.removed) {
            protocol::messages::UnsubscribeMessage unsubscription{
                .frame = frame,
                .topic = *topic,
            };

            edge->globalNetwork->unsubscribed(unsubscription);
        }
    }

    int* findSubscriptionUsingTopicPattern(const std::string& topicPattern) {
        auto subscriptions = edge->bridgedNodeSubscriptions.getSubscriptions();

        for (auto& subscription : subscriptions) {
            if (*subscription.topicPattern == topicPattern) {
                return &subscription.id;
            }
        }

        return nullptr;
    }

   public:
    NetworkEdgeStateConnected(NetworkEdgeStateManager* edge)
        : NetworkEdgeState(edge) {}

    virtual EdgeStateName getStateName() const {
        return EdgeStateName::CONNECTED;
    };

    void onSet() override {
        edge->globalNetwork->getDiagnostics()->onConnectionEstablished(
            edge->edgeInfo.bridgedNodeID, *edge->portal);

        exchangeAllNodeSubscriptions();
    };

    /**
     * NetworkParticipant interface implementation
     */

    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return edge->bridgedNodeSubscriptions.getOnlyTopLevelSubscribedTopics();
    }

    bool willHandleTopic(const std::string& topic) const override {
        return edge->bridgedNodeSubscriptions.willHandleTopic(topic);
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        return edge->bridgedNodeSubscriptions.alreadyHandlesPattern(pattern);
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        (void)frame;
        throw std::runtime_error(
            "this method should not be used, use NetworkEdge::isOriginOfFrame "
            "instead");
    }

    bool handlePublish(
        const protocol::messages::PublishMessage& publication) override {
        if (edge->isOriginOfFrame(publication.frame)) {
            return false;
        }

        bool triggered = edge->bridgedNodeSubscriptions
                             .getTriggeredSubscriptions(publication.topic)
                             .size() > 0;
        if (!triggered) {
            return false;
        }

        protocol::messages::PublishMessage forwardedPublication = publication;
        forwardedPublication.frame = internal::updateFrameTraversedAndTTL(
            forwardedPublication.frame, edge->edgeInfo.bridgedNodeID);

        if (!edge->shouldForwardMessage(forwardedPublication.frame)) {
            return false;
        }

        if (edge->edgeInfo.bridgedNodeMaxMessageSize != NO_MAX_MESSAGE_SIZE &&
            forwardedPublication.payload.size() >
                edge->edgeInfo.bridgedNodeMaxMessageSize) {
            return false;
        }

        auto result =
            edge->encoder->publish(forwardedPublication, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectingState("Failed to publish message: " +
                                        std::string(result.error));
            return false;
        }

        return true;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        if (edge->isOriginOfFrame(subscription.frame)) {
            return;
        }

        protocol::messages::SubscribeMessage subscriptionToForward =
            subscription;
        subscriptionToForward.frame = internal::updateFrameTraversedAndTTL(
            subscriptionToForward.frame, edge->edgeInfo.bridgedNodeID);

        if (!edge->shouldForwardMessage(subscriptionToForward.frame)) {
            return;
        }

        auto result =
            edge->encoder->subscribe(subscriptionToForward, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectingState("Failed to subscribe: " +
                                        std::string(result.error));
            return;
        }
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        if (edge->isOriginOfFrame(unsubscription.frame)) {
            return;
        }

        protocol::messages::UnsubscribeMessage unsubscriptionToForward =
            unsubscription;
        unsubscriptionToForward.frame = internal::updateFrameTraversedAndTTL(
            unsubscriptionToForward.frame, edge->edgeInfo.bridgedNodeID);

        if (!edge->shouldForwardMessage(unsubscriptionToForward.frame)) {
            return;
        }

        auto result =
            edge->encoder->unsubscribe(unsubscriptionToForward, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectingState("Failed to unsubscribe: " +
                                        std::string(result.error));
            return;
        }
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        protocol::messages::TerminateNetworkMessage terminationToForward =
            termination;
        terminationToForward.frame = internal::updateFrameTraversedAndTTL(
            terminationToForward.frame, edge->edgeInfo.bridgedNodeID);
        terminationToForward.frame.ttl = ONLY_DIRECT_CONNECTION_TTL;

        auto result = edge->encoder->terminateNetwork(terminationToForward,
                                                      *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectingState("Failed to terminate network: " +
                                        std::string(result.error));
            return;
        }

        edge->setDisconnectedState(termination.reason);
    }

    /**
     * DecodingHandler interface implementation
     */

    void onSupportedProtocolVersions(
        const protocol::messages::SupportedProtocolVersionsMessage& message)
        override {
        (void)message;
        edge->setDisconnectingState(
            "Unexpected supported protocol versions message in connected "
            "state");
    }

    void onInitConnection(
        const protocol::messages::InitConnectionMessage& message) override {
        (void)message;
        edge->setDisconnectingState(
            "Unexpected init connection message in connected state");
    }

    void onConnectionAccepted(
        const protocol::messages::ConnectionAcceptedMessage& message) override {
        (void)message;
        edge->setDisconnectingState(
            "Unexpected connection accepted message in connected state");
    }

    void onGracefullyClose(
        const protocol::messages::GracefullyCloseMessage& message) override {
        edge->setDisconnectedState(message.reason);
    }

    void onTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& message) override {
        edge->setDisconnectedState(message.reason);
        edge->globalNetwork->terminated(message);
    }

    void onPublish(const protocol::messages::PublishMessage& message) override {
        edge->globalNetwork->published(message);
    }

    void onSubscribe(
        const protocol::messages::SubscribeMessage& message) override {
        auto oldTopics =
            edge->bridgedNodeSubscriptions.getOnlyTopLevelSubscribedTopics();
        edge->bridgedNodeSubscriptions.addSubscription(
            std::make_shared<std::string>(message.topic), nullptr);
        updateSubscriptions(oldTopics, message.frame);
    }

    void onUnsubscribe(
        const protocol::messages::UnsubscribeMessage& message) override {
        auto subscriptionId = findSubscriptionUsingTopicPattern(message.topic);
        if (subscriptionId == nullptr) {
            return;
        }

        auto oldTopics =
            edge->bridgedNodeSubscriptions.getOnlyTopLevelSubscribedTopics();
        edge->bridgedNodeSubscriptions.removeSubscription(*subscriptionId);
        updateSubscriptions(oldTopics, message.frame);
    }

    void onMalformedMessage(
        const protocol::messages::MalformedMessage& message) override {
        (void)message;
        edge->setDisconnectingState("Malformed message received");
    }
};
}  // namespace directmq::network::edge
