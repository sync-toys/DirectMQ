#pragma once

namespace directmq::network::edge {
enum class EdgeStateName : unsigned char {
    CONNECTING,
    CONNECTED,
    DISCONNECTING,
    DISCONNECTED,
};
}
