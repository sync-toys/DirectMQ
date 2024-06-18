#pragma once
#include <inttypes.h>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct SupportedProtocolVersionsMessage {
    DataFrame frame;
    uint32_t* supportedVersions;
    size_t supportedVersionsCount;
};
}  // namespace directmq::protocol::messages
