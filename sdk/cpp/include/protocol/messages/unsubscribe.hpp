#pragma once

#include <string>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct UnsubscribeMessage {
    DataFrame frame;
    std::string topic;
};
}  // namespace directmq::protocol::messages
