#pragma once
#include "../portal.hpp"
#include "messages.hpp"

namespace directmq::protocol {
struct EncodingResult {
    const char* error;
};

class Encoder {
   public:
    virtual EncodingResult supportedProtocolVersions(
        messages::SupportedProtocolVersionsMessage& message,
        portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult initConnection(
        messages::InitConnectionMessage& message,
        portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult connectionAccepted(
        messages::ConnectionAcceptedMessage& message,
        portal::PacketWriter& packetWritere) = 0;

    virtual EncodingResult gracefullyClose(
        messages::GracefullyCloseMessage& message,
        portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult terminateNetwork(
        messages::TerminateNetworkMessage& message,
        portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult publish(messages::PublishMessage& message,
                                   portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult subscribe(messages::SubscribeMessage& message,
                                     portal::PacketWriter& packetWriter) = 0;

    virtual EncodingResult unsubscribe(messages::UnsubscribeMessage& message,
                                       portal::PacketWriter& packetWriter) = 0;
};
}  // namespace directmq::protocol
