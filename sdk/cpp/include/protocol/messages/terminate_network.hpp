#pragma once
#include "data_frame.hpp"

namespace directmq::protocol::messages {
struct TerminateNetworkMessage {
    DataFrame frame;
    char* reason;
};
}  // namespace directmq::protocol::messages
