#pragma once

#include <string>

#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct GracefullyCloseMessage {
    DataFrame frame;
    std::string reason;
};
}  // namespace directmq::protocol::messages
