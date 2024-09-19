#pragma once

#include <stdexcept>

#include "../../protocol/decoder.hpp"
#include "../../protocol/messages/unsubscribe.hpp"
#include "../constants.hpp"
#include "../participant.hpp"
#include "network_edge_state.hpp"
#include "state_name.hpp"

namespace directmq::network::edge {
class NetworkEdgeStateConnecting : public NetworkEdgeState {
   private:
    bool initializeConnection;

    bool checkIncludesSupportedVersion(
        const std::vector<uint32_t>& supportedVersions) {
        for (auto& version : supportedVersions) {
            if (version == DIRECTMQ_V1) {
                return true;
            }
        }

        return false;
    }

    void respondWithSupportedProtocolVersions(
        const protocol::messages::SupportedProtocolVersionsMessage& message) {
        protocol::messages::SupportedProtocolVersionsMessage response{
            .frame =
                protocol::messages::DataFrame{
                    .ttl = ONLY_DIRECT_CONNECTION_TTL,
                    .traversed = message.frame.traversed,
                },
            .supportedVersions = {DIRECTMQ_V1}};

        response.frame.traversed.push_back(edge->globalNetwork->config.hostID);

        edge->encoder->supportedProtocolVersions(response, *edge->portal);
    }

    void initializeEdgeConnection() {
        protocol::messages::InitConnectionMessage icMessage{
            .frame =
                protocol::messages::DataFrame{
                    .ttl = ONLY_DIRECT_CONNECTION_TTL,
                    .traversed = {edge->globalNetwork->config.hostID},
                },
            .maxMessageSize =
                edge->globalNetwork->config.hostMaxIncomingMessageSize};

        auto result = edge->encoder->initConnection(icMessage, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectedState("Connection initialization failed: " +
                                       std::string(result.error));
        }
    }

    void acceptEdgeConnection() {
        protocol::messages::ConnectionAcceptedMessage caMessage{
            .frame = protocol::messages::DataFrame{
                .ttl = ONLY_DIRECT_CONNECTION_TTL,
                .traversed = {edge->globalNetwork->config.hostID},
            }};

        auto result =
            edge->encoder->connectionAccepted(caMessage, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectedState("Connection acceptance failed: " +
                                       std::string(result.error));
            return;
        }

        edge->setConnectedState();
    }

   public:
    NetworkEdgeStateConnecting(NetworkEdgeStateManager* edge,
                               bool initializeConnection)
        : NetworkEdgeState(edge), initializeConnection(initializeConnection) {}

    virtual EdgeStateName getStateName() const {
        return EdgeStateName::CONNECTING;
    };

    void onSet() override {
        if (!initializeConnection) {
            return;
        }

        protocol::messages::SupportedProtocolVersionsMessage spvMessage{
            .frame =
                protocol::messages::DataFrame{
                    .ttl = ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL,
                    .traversed = {edge->globalNetwork->config.hostID},
                },
            .supportedVersions = {DIRECTMQ_V1}};

        auto result =
            edge->encoder->supportedProtocolVersions(spvMessage, *edge->portal);

        if (result.error != nullptr) {
            edge->setDisconnectedState(
                "Supported protocol version negotiation failed: " +
                std::string(result.error));
        }
    };

    /**
     * NetworkParticipant interface implementation
     */

    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return {};
    }

    bool willHandleTopic(const std::string& topic) const override {
        return false;
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        return false;
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        throw std::runtime_error(
            "this method should not be used, use NetworkEdge::isOriginOfFrame "
            "instead");
    }

    bool handlePublish(
        const protocol::messages::PublishMessage& publication) override {
        // we are connecting, we cannot handle any publications
        return false;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        // we are connecting, we cannot handle any subscriptions
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        // we are connecting, we cannot handle any unsubscriptions
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        // we are connecting, we cannot handle any network terminations
    }

    /**
     * DecodingHandler interface implementation
     */

    void onSupportedProtocolVersions(
        const protocol::messages::SupportedProtocolVersionsMessage& message)
        override {
        if (message.supportedVersions.size() == 0) {
            edge->setDisconnectedState(
                "No supported protocol versions received");
            return;
        }

        if (!checkIncludesSupportedVersion(message.supportedVersions)) {
            edge->setDisconnectedState("No supported protocol versions match");
            return;
        }

        if (message.frame.ttl == ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL) {
            respondWithSupportedProtocolVersions(message);
            return;
        }

        initializeEdgeConnection();
    }

    void onInitConnection(
        const protocol::messages::InitConnectionMessage& message) override {
        if (edge->edgeInfo.negotiatedProtocolVersion ==
            UNKNOWN_PROTOCOL_VERSION) {
            edge->setDisconnectingState(
                "Unknown protocol version, missing protocol negotiation");
            return;
        }

        if (message.frame.traversed.size() != 1) {
            edge->setDisconnectingState(
                "Unexpected number of traversed nodes in init connection "
                "message");
            return;
        }

        edge->edgeInfo.bridgedNodeID = *message.frame.traversed.begin();
        edge->edgeInfo.bridgedNodeMaxMessageSize = message.maxMessageSize;

        acceptEdgeConnection();
    }

    void onConnectionAccepted(
        const protocol::messages::ConnectionAcceptedMessage& message) override {
        if (edge->edgeInfo.negotiatedProtocolVersion ==
            UNKNOWN_PROTOCOL_VERSION) {
            edge->setDisconnectingState(
                "Unknown protocol version, missing protocol negotiation");
            return;
        }

        if (message.frame.traversed.size() != 1) {
            edge->setDisconnectingState(
                "Unexpected number of traversed nodes in init connection "
                "message");
            return;
        }

        edge->edgeInfo.bridgedNodeID = *message.frame.traversed.begin();
        edge->edgeInfo.bridgedNodeMaxMessageSize = message.maxMessageSize;

        edge->setConnectedState();
    }

    void onGracefullyClose(
        const protocol::messages::GracefullyCloseMessage& message) override {
        edge->setDisconnectedState(message.reason);
    }

    void onTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& message) override {
        edge->globalNetwork->terminated(message);
    }

    void onPublish(const protocol::messages::PublishMessage& message) override {
        edge->setDisconnectingState(
            "Unexpected publish message received during connection process");
    }

    void onSubscribe(
        const protocol::messages::SubscribeMessage& message) override {
        edge->setDisconnectingState(
            "Unexpected subscribe message received during connection process");
    }

    void onUnsubscribe(
        const protocol::messages::UnsubscribeMessage& message) override {
        edge->setDisconnectingState(
            "Unexpected unsubscribe message received during connection "
            "process");
    }

    void onMalformedMessage(
        const protocol::messages::MalformedMessage& message) override {
        edge->setDisconnectingState(
            "Malformed message received during connection process");
    }
};
}  // namespace directmq::network::edge
