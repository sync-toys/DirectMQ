#pragma once
#include <pb_encode.h>

#include <list>
#include <vector>

#include "../../portal.hpp"
#include "../encoder.hpp"
#include "directmq/v1/data_frame.pb.h"

namespace directmq::protocol::embedded {
class EmbeddedProtocolEncoderImplementation : public Encoder {
   private:
    template <typename T>
    std::vector<T> listToVector(const std::list<T>& lst) {
        std::vector<T> vec(lst.size());
        std::copy(lst.begin(), lst.end(), vec.begin());
        return vec;
    }

    EncodingResult writeFrame(const directmq_v1_DataFrame& frame,
                              portal::PacketWriter& packetWriter) {
        size_t frameSize = 0;
        const auto sizeCalculationSuccessful = pb_get_encoded_size(
            &frameSize, directmq_v1_DataFrame_fields, &frame);
        if (!sizeCalculationSuccessful) {
            return EncodingResult{"failed to calculate size of frame"};
        }

        pb_byte_t buffer[frameSize] = {0};
        auto stream = pb_ostream_from_buffer(buffer, frameSize);

        auto dataWriter = packetWriter.beginWrite(frameSize);

        stream.state = dataWriter.get();
        stream.callback = [](pb_ostream_t* stream, const pb_byte_t* data,
                             const size_t size) {
            portal::DataWriter* dataWriter =
                static_cast<portal::DataWriter*>(stream->state);
            return dataWriter->write(data, size);
        };

        auto success = pb_encode(&stream, directmq_v1_DataFrame_fields, &frame);
        dataWriter->end();

        if (!success) {
            return EncodingResult{PB_GET_ERROR(&stream)};
        }

        return EncodingResult{nullptr};
    }

    directmq_v1_DataFrame createFrameOf(const size_t& whichMessage,
                                        messages::DataFrame& config) {
        directmq_v1_DataFrame frame = directmq_v1_DataFrame_init_zero;

        const char** traversed = new const char*[config.traversed.size()];
        std::list<std::string>::iterator it = config.traversed.begin();
        for (size_t i = 0; i < config.traversed.size(); i++, it++) {
            traversed[i] = (*it).c_str();
        }

        frame.which_message = whichMessage;
        frame.ttl = &config.ttl;
        frame.traversed = const_cast<char**>(traversed);
        frame.traversed_count = config.traversed.size();

        return frame;
    }

    void deleteFrame(const directmq_v1_DataFrame& frame) {
        delete[] frame.traversed;
    }

   public:
    EncodingResult supportedProtocolVersions(
        messages::SupportedProtocolVersionsMessage& message,
        portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_supported_protocol_versions_tag,
            message.frame);

        directmq_v1_SupportedProtocolVersions encoded =
            directmq_v1_SupportedProtocolVersions_init_zero;

        auto supportedVersions = message.supportedVersions;

        encoded.supported_protocol_versions = supportedVersions.data();
        encoded.supported_protocol_versions_count = supportedVersions.size();

        frame.message.supported_protocol_versions = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    }

    EncodingResult initConnection(messages::InitConnectionMessage& message,
                                  portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_init_connection_tag, message.frame);

        directmq_v1_InitConnection encoded =
            directmq_v1_InitConnection_init_zero;

        encoded.max_message_size = &message.maxMessageSize;

        frame.message.init_connection = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    }

    EncodingResult connectionAccepted(
        messages::ConnectionAcceptedMessage& message,
        portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_connection_accepted_tag, message.frame);

        directmq_v1_ConnectionAccepted encoded =
            directmq_v1_ConnectionAccepted_init_zero;

        encoded.max_message_size = &message.maxMessageSize;

        frame.message.connection_accepted = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };

    EncodingResult gracefullyClose(
        messages::GracefullyCloseMessage& message,
        portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_gracefully_close_tag, message.frame);

        directmq_v1_GracefullyClose encoded =
            directmq_v1_GracefullyClose_init_zero;

        encoded.reason = const_cast<char*>(message.reason.c_str());

        frame.message.gracefully_close = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };

    EncodingResult terminateNetwork(
        messages::TerminateNetworkMessage& message,
        portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(
            directmq_v1_DataFrame_terminate_network_tag, message.frame);

        directmq_v1_TerminateNetwork encoded =
            directmq_v1_TerminateNetwork_init_zero;

        encoded.reason = const_cast<char*>(message.reason.c_str());

        frame.message.terminate_network = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };

    EncodingResult publish(messages::PublishMessage& message,
                           portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_publish_tag,
                                         message.frame);

        directmq_v1_DeliveryStrategy deliveryStrategy =
            message.deliveryStrategy ==
                    messages::DeliveryStrategy::AT_LEAST_ONCE
                ? directmq_v1_DeliveryStrategy_DELIVERY_STRATEGY_AT_LEAST_ONCE_UNSPECIFIED
                : directmq_v1_DeliveryStrategy_DELIVERY_STRATEGY_AT_MOST_ONCE;

        directmq_v1_Publish encoded = directmq_v1_Publish_init_zero;

        encoded.topic = const_cast<char*>(message.topic.c_str());
        encoded.delivery_strategy = &deliveryStrategy;

        auto payloadSize = message.payload.size();
        encoded.size = &payloadSize;

        size_t totalSize = PB_BYTES_ARRAY_T_ALLOCSIZE(message.payload.size());
        pb_bytes_array_t* payloadPtr = (pb_bytes_array_t*)malloc(totalSize);
        if (payloadPtr == nullptr) {
            return EncodingResult{"failed to allocate memory for payload"};
        }

        payloadPtr->size = message.payload.size();
        memcpy(payloadPtr->bytes, message.payload.data(),
               message.payload.size());

        encoded.payload = payloadPtr;

        frame.message.publish = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };

    EncodingResult subscribe(messages::SubscribeMessage& message,
                             portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_subscribe_tag,
                                         message.frame);

        directmq_v1_Subscribe encoded = directmq_v1_Subscribe_init_zero;

        encoded.topic = const_cast<char*>(message.topic.c_str());

        frame.message.subscribe = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };

    EncodingResult unsubscribe(messages::UnsubscribeMessage& message,
                               portal::PacketWriter& packetWriter) override {
        auto frame = this->createFrameOf(directmq_v1_DataFrame_unsubscribe_tag,
                                         message.frame);

        directmq_v1_Unsubscribe encoded = directmq_v1_Unsubscribe_init_zero;

        encoded.topic = const_cast<char*>(message.topic.c_str());

        frame.message.unsubscribe = &encoded;

        auto result = this->writeFrame(frame, packetWriter);
        this->deleteFrame(frame);

        return result;
    };
};
}  // namespace directmq::protocol::embedded
