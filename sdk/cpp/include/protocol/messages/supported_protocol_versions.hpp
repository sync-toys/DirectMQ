#pragma once

#include <inttypes.h>

#include <list>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct SupportedProtocolVersionsMessage {
    DataFrame frame;
    std::list<uint32_t> supportedVersions;
};
}  // namespace directmq::protocol::messages
