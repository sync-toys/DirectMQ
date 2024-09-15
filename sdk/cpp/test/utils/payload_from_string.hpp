#pragma once

#include <string>
#include <vector>

std::vector<uint8_t> payloadFromString(const std::string& str) {
    std::vector<uint8_t> payload(str.size());
    std::copy(str.begin(), str.end(), payload.begin());
    return payload;
}