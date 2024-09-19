#pragma once

#include <cstddef>
#include <cstdint>
#include <string>

namespace directmq::network {
struct NetworkNodeConfig {
    int32_t hostTTL;
    size_t hostMaxIncomingMessageSize;
    std::string hostID;
};
}  // namespace directmq::network
