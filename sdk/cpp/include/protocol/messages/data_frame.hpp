#pragma once

#include <inttypes.h>
#include <stddef.h>

namespace directmq::protocol::messages {
struct DataFrame {
    int32_t ttl;
    char** traversed;
    size_t traversedCount;
};
}  // namespace directmq::protocol::messages
