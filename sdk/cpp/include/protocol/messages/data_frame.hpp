#pragma once

#include <inttypes.h>
#include <stddef.h>

#include <list>
#include <string>

namespace directmq::protocol::messages {
struct DataFrame {
    int32_t ttl;
    std::list<std::string> traversed;
};
}  // namespace directmq::protocol::messages
