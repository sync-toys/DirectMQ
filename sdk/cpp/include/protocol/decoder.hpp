#pragma once
#include "../portal.hpp"
#include "messages.hpp"

namespace directmq::protocol {
struct DecodingResult {
    const char* error;
};

class DecodingHandler {
   public:
    virtual void onSupportedProtocolVersions(
        const messages::SupportedProtocolVersionsMessage& message) = 0;
    virtual void onInitConnection(
        const messages::InitConnectionMessage& message) = 0;
    virtual void onConnectionAccepted(
        const messages::ConnectionAcceptedMessage& message) = 0;
    virtual void onGracefullyClose(
        const messages::GracefullyCloseMessage& message) = 0;
    virtual void onTerminateNetwork(
        const messages::TerminateNetworkMessage& message) = 0;
    virtual void onPublish(const messages::PublishMessage& message) = 0;
    virtual void onSubscribe(const messages::SubscribeMessage& message) = 0;
    virtual void onUnsubscribe(const messages::UnsubscribeMessage& message) = 0;
    virtual void onMalformedMessage(
        const messages::MalformedMessage& message) = 0;
};

class Decoder {
   public:
    virtual DecodingResult readMessage(portal::PacketReader* packetReader,
                                       DecodingHandler* handler) = 0;
};
}  // namespace directmq::protocol
