#pragma once

#include <inttypes.h>

#include <vector>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct SupportedProtocolVersionsMessage {
    DataFrame frame;
    std::vector<uint32_t> supportedVersions;
};
}  // namespace directmq::protocol::messages
