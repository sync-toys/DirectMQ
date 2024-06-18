#pragma once
#include <inttypes.h>

#include "../../bytes.hpp"

namespace directmq::protocol::messages {
struct MalformedMessage {
    bytes message;
    size_t messageSize;
};
}  // namespace directmq::protocol::messages
