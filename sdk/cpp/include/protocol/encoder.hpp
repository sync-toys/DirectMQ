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
        PacketWriter& packetWriter) = 0;

    virtual EncodingResult initConnection(
        messages::InitConnectionMessage& message,
        PacketWriter& packetWriter) = 0;

    virtual EncodingResult connectionAccepted(
        messages::ConnectionAcceptedMessage& message,
        PacketWriter& packetWritere) = 0;

    virtual EncodingResult gracefullyClose(
        messages::GracefullyCloseMessage& message,
        PacketWriter& packetWriter) = 0;

    virtual EncodingResult terminateNetwork(
        messages::TerminateNetworkMessage& message,
        PacketWriter& packetWriter) = 0;

    virtual EncodingResult publish(messages::PublishMessage& message,
                                   PacketWriter& packetWriter) = 0;

    virtual EncodingResult subscribe(messages::SubscribeMessage& message,
                                     PacketWriter& packetWriter) = 0;

    virtual EncodingResult unsubscribe(messages::UnsubscribeMessage& message,
                                       PacketWriter& packetWriter) = 0;
};
}  // namespace directmq::protocol
