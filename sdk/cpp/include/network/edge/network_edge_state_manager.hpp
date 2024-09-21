#pragma once

#include "../../portal.hpp"
#include "../../protocol/decoder.hpp"
#include "../../protocol/encoder.hpp"
#include "../../subscription_list.hpp"
#include "../constants.hpp"
#include "../global.hpp"
#include "../participant.hpp"
#include "edge_info.hpp"
#include "network_edge_state.hpp"
#include "state_name.hpp"
#include "utilities.hpp"

namespace directmq::network::edge {
// Purpose of NetworkEdgeStateManager is to break
// the circular dependency between network edge states:
// !!! disconnected !!! -> connecting -> connected -> disconnecting -> !!!
// disconnected !!!
class NetworkEdgeStateManager : public NetworkParticipant,
                                public protocol::DecodingHandler {
   private:
    std::shared_ptr<GlobalNetwork> globalNetwork;

    std::shared_ptr<portal::Portal> portal;

    std::shared_ptr<protocol::Decoder> decoder;
    std::shared_ptr<protocol::Encoder> encoder;

    EdgeInfo edgeInfo;

    subscriptions::SubscriptionList<void*> bridgedNodeSubscriptions;

    void setState(NetworkEdgeState* newState) {
        auto oldState = state;

        state = newState;
        state->onSet();

        delete oldState;
    }

    friend class NetworkEdgeStateDisconnected;
    friend class NetworkEdgeStateConnecting;
    friend class NetworkEdgeStateConnected;
    friend class NetworkEdgeStateDisconnecting;

   protected:
    NetworkEdgeState* state = nullptr;

   public:
    virtual ~NetworkEdgeStateManager() = default;

    virtual void setConnectingState(bool initializeConnection) = 0;
    virtual void setConnectedState() = 0;
    virtual void setDisconnectingState(const std::string& reason) = 0;
    virtual void setDisconnectedState(const std::string& reason) = 0;

    NetworkEdgeStateManager(std::shared_ptr<GlobalNetwork> globalNetwork,
                            std::shared_ptr<portal::Portal> portal,
                            std::shared_ptr<protocol::Decoder> decoder,
                            std::shared_ptr<protocol::Encoder> encoder)
        : globalNetwork(globalNetwork),
          portal(portal),
          decoder(decoder),
          encoder(encoder),
          edgeInfo{.bridgedNodeID = "",
                   .bridgedNodeMaxMessageSize = NO_MAX_MESSAGE_SIZE,
                   .bridgedNodeSupportedProtocolVersions = {},
                   .negotiatedProtocolVersion = 0} {}

    protocol::DecodingResult processIncomingPacket(
        std::shared_ptr<portal::Packet> packet) {
        return decoder->decodePacket(packet, this);
    }

    std::shared_ptr<portal::Portal> getPortal() const { return portal; }

    EdgeStateName getStateName() const { return state->getStateName(); }

    bool shouldForwardMessage(
        const protocol::messages::DataFrame& frame) const {
        bool loopDetected = internal::checkForNetworkLoops(frame);
        if (loopDetected) {
            globalNetwork->terminated(
                protocol::messages::TerminateNetworkMessage{
                    .frame = frame,
                    .reason = internal::buildLoopDetectedMessage(frame)});
        }

        return !loopDetected && frame.ttl > 0;
    }

    /**
     * NetworkParticipant interface implementation
     */

    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return state->getSubscribedTopics();
    }

    bool willHandleTopic(const std::string& topic) const override {
        return state->willHandleTopic(topic);
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        return state->alreadyHandlesPattern(pattern);
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        if (frame.traversed.size() == 0) {
            // in case of host node origin, we are returning false
            // because network edges are used only for bridged nodes
            return false;
        }

        return *frame.traversed.end() == edgeInfo.bridgedNodeID;
    }

    bool handlePublish(
        const protocol::messages::PublishMessage& publication) override {
        return state->handlePublish(publication);
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        state->handleSubscribe(subscription);
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        state->handleUnsubscribe(unsubscription);
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        state->handleTerminateNetwork(termination);
    }

    /**
     * DecodingHandler interface implementation
     */

    void onSupportedProtocolVersions(
        const protocol::messages::SupportedProtocolVersionsMessage& message)
        override {
        state->onSupportedProtocolVersions(message);
    }

    void onInitConnection(
        const protocol::messages::InitConnectionMessage& message) override {
        state->onInitConnection(message);
    }

    void onConnectionAccepted(
        const protocol::messages::ConnectionAcceptedMessage& message) override {
        state->onConnectionAccepted(message);
    }

    void onGracefullyClose(
        const protocol::messages::GracefullyCloseMessage& message) override {
        state->onGracefullyClose(message);
    }

    void onTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& message) override {
        state->onTerminateNetwork(message);
    }

    void onPublish(const protocol::messages::PublishMessage& message) override {
        state->onPublish(message);
    }

    void onSubscribe(
        const protocol::messages::SubscribeMessage& message) override {
        state->onSubscribe(message);
    }

    void onUnsubscribe(
        const protocol::messages::UnsubscribeMessage& message) override {
        state->onUnsubscribe(message);
    }

    void onMalformedMessage(
        const protocol::messages::MalformedMessage& message) override {
        state->onMalformedMessage(message);
    }
};
}  // namespace directmq::network::edge
