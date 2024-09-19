#pragma once

#include <inttypes.h>

#include <string>
#include <vector>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
enum DeliveryStrategy { AT_LEAST_ONCE = 0, AT_MOST_ONCE = 1 };

struct PublishMessage {
    DataFrame frame;
    std::string topic;
    DeliveryStrategy deliveryStrategy;
    std::vector<uint8_t> payload;
};
}  // namespace directmq::protocol::messages
