#pragma once

#include <inttypes.h>

#include <vector>

#include "../../bytes.hpp"

namespace directmq::protocol::messages {
struct MalformedMessage {
    std::vector<uint8_t> bytes;
    std::string error;
};
}  // namespace directmq::protocol::messages
