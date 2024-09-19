#pragma once

#include <cstddef>
#include <cstdint>

namespace directmq::network {
constexpr int32_t DEFAULT_TTL = 32;
constexpr int32_t ONLY_DIRECT_CONNECTION_TTL = 1;
constexpr int32_t ONLY_DIRECT_CONNECTION_WITH_RESPONSE_TTL = 2;
constexpr size_t NO_MAX_MESSAGE_SIZE = 0;
constexpr uint32_t UNKNOWN_PROTOCOL_VERSION = 0;
constexpr uint32_t DIRECTMQ_V1 = 1;
}  // namespace directmq::network
