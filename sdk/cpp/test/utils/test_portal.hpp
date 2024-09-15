#pragma once
#include <queue>

#include "portal.hpp"
#include "test_writer.hpp"

class TestPortal : public portal::Portal {
   private:
    std::queue<portal::Packet> packets;

   public:
    void close() {}

    portal::DataWriter* beginWrite(const size_t packetSize) {
        auto writer = new TestWriter(&this->packets, packetSize);
        return writer;
    }

    portal::Packet readPacket() {
        if (this->packets.empty()) {
            auto packet = portal::Packet{0, nullptr};
            return packet;
        }

        auto packet = packets.front();
        packets.pop();

        return packet;
    }
};
