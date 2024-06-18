#pragma once
#include <pb_encode.h>

#include <cstdlib>

#include "../../portal.hpp"
#include "../encoder.hpp"
#include "directmq/v1/data_frame.pb.h"

namespace directmq::protocol::embedded {
class EmbeddedProtocolEncoderImplementation : public Encoder {
   private:
    EncodingResult writeFrame(const directmq_v1_DataFrame& frame,
                              PacketWriter& packetWriter) {
        size_t frameSize = 0;
        const auto sizeCalculationSuccessful = pb_get_encoded_size(
            &frameSize, directmq_v1_DataFrame_fields, &frame);
        if (!sizeCalculationSuccessful) {
            return EncodingResult{"failed to calculate size of frame"};
        }

        pb_byte_t buffer[frameSize] = {0};
        auto stream = pb_ostream_from_buffer(buffer, frameSize);

        DataWriter* dataWriter = packetWriter.beginWrite(frameSize);

        stream.state = dataWriter;
        stream.callback = [](pb_ostream_t* stream, const pb_byte_t* data,
                             const size_t size) {
            DataWriter* dataWriter = static_cast<DataWriter*>(stream->state);
            return dataWriter->write(data, size);
        };

        auto success = pb_encode(&stream, directmq_v1_DataFrame_fields, &frame);
        dataWriter->end();
        delete dataWriter;

        if (!success) {
            return EncodingResult{PB_GET_ERROR(&stream)};
        }

        return EncodingResult{nullptr};
    }

    directmq_v1_DataFrame createFrameOf(const size_t& whichMessage,
                                        messages::DataFrame& config) {
        directmq_v1_DataFrame frame = directmq_v1_DataFrame_init_zero;

        frame.which_message = whichMessage;
        frame.ttl = &config.ttl;
        frame.traversed = config.traversed;
        frame.traversed_count = config.traversedCount;

        return frame;
    }

   public:
    EncodingResult supportedProtocolVersions(
        messages::SupportedProtocolVersionsMessage& message,
        PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_supported_protocol_versions_tag,
            message.frame);

        directmq_v1_SupportedProtocolVersions encoded =
            directmq_v1_SupportedProtocolVersions_init_zero;

        encoded.supported_protocol_versions = message.supportedVersions;
        encoded.supported_protocol_versions_count =
            message.supportedVersionsCount;

        frame.message.supported_protocol_versions = &encoded;

        return this->writeFrame(frame, packetWriter);
    }

    EncodingResult initConnection(messages::InitConnectionMessage& message,
                                  PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_init_connection_tag, message.frame);

        directmq_v1_InitConnection encoded =
            directmq_v1_InitConnection_init_zero;

        encoded.max_message_size = &message.maxMessageSize;

        frame.message.init_connection = &encoded;

        return this->writeFrame(frame, packetWriter);
    }

    EncodingResult connectionAccepted(
        messages::ConnectionAcceptedMessage& message,
        PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_connection_accepted_tag, message.frame);

        directmq_v1_ConnectionAccepted encoded =
            directmq_v1_ConnectionAccepted_init_zero;

        encoded.max_message_size = &message.maxMessageSize;

        frame.message.connection_accepted = &encoded;

        return this->writeFrame(frame, packetWriter);
    };

    EncodingResult gracefullyClose(messages::GracefullyCloseMessage& message,
                                   PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_gracefully_close_tag, message.frame);

        directmq_v1_GracefullyClose encoded =
            directmq_v1_GracefullyClose_init_zero;

        encoded.reason = message.reason;

        frame.message.gracefully_close = &encoded;

        return this->writeFrame(frame, packetWriter);
    };

    EncodingResult terminateNetwork(messages::TerminateNetworkMessage& message,
                                    PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_terminate_network_tag, message.frame);

        directmq_v1_TerminateNetwork encoded =
            directmq_v1_TerminateNetwork_init_zero;

        encoded.reason = message.reason;

        frame.message.terminate_network = &encoded;

        return this->writeFrame(frame, packetWriter);
    };

    EncodingResult publish(messages::PublishMessage& message,
                           PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_publish_tag,
                                         message.frame);

        directmq_v1_DeliveryStrategy deliveryStrategy =
            message.deliveryStrategy ==
                    messages::DeliveryStrategy::AT_LEAST_ONCE
                ? directmq_v1_DeliveryStrategy_DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED
                : directmq_v1_DeliveryStrategy_DELIVERY_STRATEGY_AT_MOST_ONCE;

        directmq_v1_Publish encoded = directmq_v1_Publish_init_zero;

        encoded.topic = message.topic;
        encoded.delivery_strategy = &deliveryStrategy;
        encoded.size = &message.payloadSize;

        size_t totalSize = PB_BYTES_ARRAY_T_ALLOCSIZE(message.payloadSize);
        pb_bytes_array_t* payloadPtr = (pb_bytes_array_t*)malloc(totalSize);
        if (payloadPtr == nullptr) {
            return EncodingResult{"failed to allocate memory for payload"};
        }

        payloadPtr->size = message.payloadSize;
        memcpy(payloadPtr->bytes, message.payload, message.payloadSize);

        encoded.payload = payloadPtr;

        frame.message.publish = &encoded;

        return this->writeFrame(frame, packetWriter);
    };

    EncodingResult subscribe(messages::SubscribeMessage& message,
                             PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_subscribe_tag,
                                         message.frame);

        directmq_v1_Subscribe encoded = directmq_v1_Subscribe_init_zero;

        encoded.topic = message.topic;

        frame.message.subscribe = &encoded;

        return this->writeFrame(frame, packetWriter);
    };

    EncodingResult unsubscribe(messages::UnsubscribeMessage& message,
                               PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_unsubscribe_tag,
                                         message.frame);

        directmq_v1_Unsubscribe encoded = directmq_v1_Unsubscribe_init_zero;

        encoded.topic = message.topic;

        frame.message.unsubscribe = &encoded;

        return this->writeFrame(frame, packetWriter);
    };
};
}  // namespace directmq::protocol::embedded
