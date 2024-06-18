#pragma once
#include <pb_decode.h>

#include "../../portal.hpp"
#include "../decoder.hpp"
#include "directmq/v1/data_frame.pb.h"

namespace directmq::protocol::embedded {
class EmbeddedProtocolDecoderImplementation : public Decoder {
   private:
    DecodingResult handleMessage(directmq_v1_DataFrame& frame,
                                 DecodingHandler* handler) {
        messages::DataFrame decodedFrame{
            .ttl = *frame.ttl,
            .traversed = frame.traversed,
            .traversedCount = frame.traversed_count};

        switch (frame.which_message) {
            case directmq_v1_DataFrame_supported_protocol_versions_tag: {
                directmq_v1_SupportedProtocolVersions encoded =
                    *frame.message.supported_protocol_versions;

                messages::SupportedProtocolVersionsMessage
                    supportedProtolVersionsMessage{
                        .frame = decodedFrame,
                        .supportedVersions =
                            encoded.supported_protocol_versions,
                        .supportedVersionsCount =
                            encoded.supported_protocol_versions_count};

                handler->onSupportedProtocolVersions(
                    supportedProtolVersionsMessage);
                return DecodingResult{nullptr};
            }
            case directmq_v1_DataFrame_init_connection_tag: {
                directmq_v1_InitConnection encoded =
                    *frame.message.init_connection;

                messages::InitConnectionMessage initConnectionMessage{
                    .frame = decodedFrame,
                    .maxMessageSize = *encoded.max_message_size};

                handler->onInitConnection(initConnectionMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_connection_accepted_tag: {
                directmq_v1_ConnectionAccepted encoded =
                    *frame.message.connection_accepted;

                messages::ConnectionAcceptedMessage connectionAcceptedMessage{
                    .frame = decodedFrame,
                    .maxMessageSize = *encoded.max_message_size};

                handler->onConnectionAccepted(connectionAcceptedMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_gracefully_close_tag: {
                directmq_v1_GracefullyClose encoded =
                    *frame.message.gracefully_close;

                messages::GracefullyCloseMessage gracefullyCloseMessage{
                    .frame = decodedFrame, .reason = encoded.reason};

                handler->onGracefullyClose(gracefullyCloseMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_terminate_network_tag: {
                directmq_v1_TerminateNetwork encoded =
                    *frame.message.terminate_network;

                messages::TerminateNetworkMessage terminateNetworkMessage{
                    .frame = decodedFrame, .reason = encoded.reason};

                handler->onTerminateNetwork(terminateNetworkMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_publish_tag: {
                directmq_v1_Publish encoded = *frame.message.publish;

                messages::PublishMessage publishMessage{
                    .frame = decodedFrame,
                    .topic = encoded.topic,
                    .deliveryStrategy = static_cast<messages::DeliveryStrategy>(
                        *encoded.delivery_strategy),
                    .payload = encoded.payload->bytes,
                    .payloadSize = *encoded.size};

                handler->onPublish(publishMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_subscribe_tag: {
                directmq_v1_Subscribe encoded = *frame.message.subscribe;

                messages::SubscribeMessage subscribeMessage{
                    .frame = decodedFrame, .topic = encoded.topic};

                handler->onSubscribe(subscribeMessage);
                return DecodingResult{nullptr};
            }

            case directmq_v1_DataFrame_unsubscribe_tag: {
                directmq_v1_Unsubscribe encoded = *frame.message.unsubscribe;

                messages::UnsubscribeMessage unsubscribeMessage{
                    .frame = decodedFrame, .topic = encoded.topic};

                handler->onUnsubscribe(unsubscribeMessage);
                return DecodingResult{nullptr};
            }
        }

        return DecodingResult{"unsupported tag provided"};
    }

   public:
    DecodingResult readMessage(PacketReader* packetReader,
                               DecodingHandler* handler) {
        auto packet = packetReader->readPacket();
        if (packet.size == 0) {
            return DecodingResult{nullptr};
        }

        directmq_v1_DataFrame frame = directmq_v1_DataFrame_init_zero;

        pb_istream_t stream = pb_istream_from_buffer(packet.data, packet.size);
        bool successfull =
            pb_decode(&stream, directmq_v1_DataFrame_fields, &frame);

        if (!successfull) {
            return DecodingResult{PB_GET_ERROR(&stream)};
        }

        return this->handleMessage(frame, handler);
    }
};
}  // namespace directmq::protocol::embedded
