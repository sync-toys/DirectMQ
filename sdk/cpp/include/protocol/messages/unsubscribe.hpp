#pragma once
#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct UnsubscribeMessage {
    DataFrame frame;
    char* topic;
};
}  // namespace directmq::protocol::messages
