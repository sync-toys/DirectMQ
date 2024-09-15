#pragma once

#include <string>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct SubscribeMessage {
    DataFrame frame;
    std::string topic;
};
}  // namespace directmq::protocol::messages
