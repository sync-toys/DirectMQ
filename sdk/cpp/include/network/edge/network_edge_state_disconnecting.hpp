#pragma once

#include <stdexcept>

#include "../../protocol/decoder.hpp"
#include "../../protocol/messages/unsubscribe.hpp"
#include "../constants.hpp"
#include "../participant.hpp"
#include "network_edge_state.hpp"
#include "state_name.hpp"

namespace directmq::network::edge {
class NetworkEdgeStateDisconnecting : public NetworkEdgeState {
   private:
    std::string reason;

   public:
    NetworkEdgeStateDisconnecting(NetworkEdgeStateManager* edge,
                                  const std::string& reason)
        : NetworkEdgeState(edge), reason(reason) {}

    virtual EdgeStateName getStateName() const {
        return EdgeStateName::DISCONNECTING;
    };

    void onSet() override {
        protocol::messages::GracefullyCloseMessage gcMessage{
            .frame =
                protocol::messages::DataFrame{
                    .ttl = ONLY_DIRECT_CONNECTION_TTL,
                    .traversed = {edge->globalNetwork->config.hostID},
                },
            .reason = reason};

        edge->encoder->gracefullyClose(gcMessage, *edge->portal);

        edge->setDisconnectedState(reason);
    };

    /**
     * NetworkParticipant interface implementation
     */

    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return {};
    }

    bool willHandleTopic(const std::string& topic) const override {
        (void)topic;
        return false;
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        (void)pattern;
        return false;
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
        // we are disconnecting, we cannot handle any publications
        (void)publication;
        return false;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        // we are disconnecting, we cannot handle any subscriptions
        (void)subscription;
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        // we are disconnecting, we cannot handle any unsubscriptions
        (void)unsubscription;
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        // we are disconnecting, we cannot handle any network terminations
        (void)termination;
    }

    /**
     * DecodingHandler interface implementation
     */

    void onSupportedProtocolVersions(
        const protocol::messages::SupportedProtocolVersionsMessage& message)
        override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onInitConnection(
        const protocol::messages::InitConnectionMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onConnectionAccepted(
        const protocol::messages::ConnectionAcceptedMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onGracefullyClose(
        const protocol::messages::GracefullyCloseMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onPublish(const protocol::messages::PublishMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onSubscribe(
        const protocol::messages::SubscribeMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onUnsubscribe(
        const protocol::messages::UnsubscribeMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }

    void onMalformedMessage(
        const protocol::messages::MalformedMessage& message) override {
        // we are disconnecting, we cannot handle any messages
        (void)message;
    }
};
}  // namespace directmq::network::edge
