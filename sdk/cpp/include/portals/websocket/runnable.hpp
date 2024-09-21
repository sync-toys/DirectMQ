#pragma once

namespace directmq::portal::websocket {
class Runnable {
   public:
    ~Runnable() = default;
    virtual void run(int timeoutInMs) = 0;
};
}  // namespace directmq::portal::websocket
