#pragma once
#include <inttypes.h>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct InitConnectionMessage {
    DataFrame frame;
    uint64_t maxMessageSize;
};
}  // namespace directmq::protocol::messages
