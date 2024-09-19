#pragma once
#include <queue>

#include "portal.hpp"
#include "test_writer.hpp"

class TestPortal : public portal::Portal, public portal::PacketReader {
   private:
    std::queue<std::shared_ptr<portal::Packet>> packets;

   public:
    void close() {}

    std::shared_ptr<portal::DataWriter> beginWrite(const size_t packetSize) {
        auto writer = std::make_shared<TestWriter>(&this->packets, packetSize);
        return writer;
    }

    std::shared_ptr<portal::Packet> readPacket() {
        if (this->packets.empty()) {
            auto packet = portal::Packet::empty();
            return packet;
        }

        auto packet = packets.front();
        packets.pop();

        return packet;
    }
};
