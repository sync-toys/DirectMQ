#pragma once
#include <pb_decode.h>

#include <iostream>

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
            case directmq_v1_DataFrame_supported_protocol_versions_tag:
                directmq_v1_SupportedProtocolVersions encoded =
                    *frame.message.supported_protocol_versions;
                messages::SupportedProtocolVersionsMessage message{
                    .frame = decodedFrame,
                    .supportedVersions = encoded.supported_protocol_versions,
                    .supportedVersionsCount =
                        encoded.supported_protocol_versions_count};
                handler->onSupportedProtocolVersions(message);
                return DecodingResult{nullptr};
        }

        return DecodingResult{"unsupported tag provided"};
    }

   public:
    DecodingResult readMessage(PacketReader* packetReader,
                               DecodingHandler* handler) {
        auto packet = packetReader->readPacket();
        if (packet.size == 0) {
            std::cout << "No data to decode" << std::endl;
            return DecodingResult{nullptr};
        }

        directmq_v1_DataFrame frame = directmq_v1_DataFrame_init_zero;

        pb_istream_t stream = pb_istream_from_buffer(packet.data, packet.size);
        bool successfull =
            pb_decode(&stream, directmq_v1_DataFrame_fields, &frame);

        if (!successfull) {
            std::cout << "Failed to decode frame" << std::endl;
            return DecodingResult{PB_GET_ERROR(&stream)};
        }

        std::cout << "Decoded frame" << std::endl;
        return this->handleMessage(frame, handler);
    }
};
}  // namespace directmq::protocol::embedded
