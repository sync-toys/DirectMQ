#pragma once
#include <inttypes.h>

#include "../../bytes.hpp"
#include "data_frame.hpp"

namespace directmq::protocol::messages {
enum DeliveryStrategy { AT_LEAST_ONCE = 0, AT_MOST_ONCE = 1 };

struct PublishMessage {
    DataFrame frame;
    char* topic;
    DeliveryStrategy deliveryStrategy;
    bytes payload;
    size_t payloadSize;
};
}  // namespace directmq::protocol::messages
