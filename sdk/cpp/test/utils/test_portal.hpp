#pragma once
#include <queue>

#include "portal.hpp"
#include "test_writer.hpp"

class TestPortal : public Portal {
   private:
    std::queue<Packet> packets;

   public:
    void close() {}

    DataWriter* beginWrite(const size_t packetSize) {
        auto writer = new TestWriter(&this->packets, packetSize);
        return writer;
    }

    Packet readPacket() {
        if (this->packets.empty()) {
            auto packet = Packet{0, nullptr};
            return packet;
        }

        auto packet = packets.front();
        packets.pop();

        return packet;
    }
};
