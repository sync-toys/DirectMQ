#pragma once
#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct GracefullyCloseMessage {
    DataFrame frame;
    char* reason;
};
}  // namespace directmq::protocol::messages
