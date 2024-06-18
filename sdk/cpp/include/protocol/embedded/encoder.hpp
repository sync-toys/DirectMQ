#pragma once
#include <pb_encode.h>

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

    EncodingResult initConnection(messages::InitConnectionMessage& message
                                  __attribute__((unused)),
                                  PacketWriter& packetWriter
                                  __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    }

    EncodingResult connectionAccepted(
        messages::ConnectionAcceptedMessage& message __attribute__((unused)),
        PacketWriter& packetWritere __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };

    EncodingResult gracefullyClose(messages::GracefullyCloseMessage& message
                                   __attribute__((unused)),
                                   PacketWriter& packetWriter
                                   __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };

    EncodingResult terminateNetwork(messages::TerminateNetworkMessage& message
                                    __attribute__((unused)),
                                    PacketWriter& packetWriter
                                    __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };

    EncodingResult publish(messages::PublishMessage& message
                           __attribute__((unused)),
                           PacketWriter& packetWriter
                           __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };

    EncodingResult subscribe(messages::SubscribeMessage& message
                             __attribute__((unused)),
                             PacketWriter& packetWriter
                             __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };

    EncodingResult unsubscribe(messages::UnsubscribeMessage& message
                               __attribute__((unused)),
                               PacketWriter& packetWriter
                               __attribute__((unused))) override {
        return EncodingResult{"not implemented"};
    };
};
}  // namespace directmq::protocol::embedded
