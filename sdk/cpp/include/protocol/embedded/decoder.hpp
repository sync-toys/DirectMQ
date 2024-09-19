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
        std::list<std::string> traversed;
        for (size_t i = 0; i < frame.traversed_count; i++) {
            traversed.push_back(std::string(frame.traversed[i]));
        }

        messages::DataFrame decodedFrame{.ttl = *frame.ttl,
                                         .traversed = traversed};

        switch (frame.which_message) {
            case directmq_v1_DataFrame_supported_protocol_versions_tag: {
                directmq_v1_SupportedProtocolVersions encoded =
                    *frame.message.supported_protocol_versions;

                std::vector<uint32_t> supportedVersions;
                for (size_t i = 0;
                     i < encoded.supported_protocol_versions_count; i++) {
                    supportedVersions.push_back(
                        encoded.supported_protocol_versions[i]);
                }

                messages::SupportedProtocolVersionsMessage
                    supportedProtolVersionsMessage{
                        .frame = decodedFrame,
                        .supportedVersions = supportedVersions};

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
                    .payload = std::vector<uint8_t>(encoded.payload->size)};

                memcpy(publishMessage.payload.data(), encoded.payload->bytes,
                       encoded.payload->size);

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
    DecodingResult decodePacket(std::shared_ptr<portal::Packet> packet,
                                DecodingHandler* handler) {
        if (packet->size == 0) {
            return DecodingResult{nullptr};
        }

        directmq_v1_DataFrame frame = directmq_v1_DataFrame_init_zero;

        pb_istream_t stream =
            pb_istream_from_buffer(packet->data, packet->size);
        bool successfull =
            pb_decode(&stream, directmq_v1_DataFrame_fields, &frame);

        if (!successfull) {
            auto error = PB_GET_ERROR(&stream);
            messages::MalformedMessage malformedMessage{
                .bytes = std::vector<uint8_t>(packet->size),
                .error = std::string(error)};

            std::copy(packet->data, packet->data + packet->size,
                      malformedMessage.bytes.begin());

            handler->onMalformedMessage(malformedMessage);
            return DecodingResult{error};
        }

        return this->handleMessage(frame, handler);
    }
};
}  // namespace directmq::protocol::embedded
