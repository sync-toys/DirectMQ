#pragma once

#include <cstddef>
#include <cstdint>
#include <string>
#include <vector>

namespace directmq::network::edge {
struct EdgeInfo {
    std::string bridgedNodeID;
    size_t bridgedNodeMaxMessageSize;
    std::vector<uint32_t> bridgedNodeSupportedProtocolVersions;
    int32_t negotiatedProtocolVersion;
};
}  // namespace directmq::network::edge
