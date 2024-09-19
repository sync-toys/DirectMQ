#pragma once

#include <pb.h>

#include <memory>
#include <queue>

#include "portal.hpp"

using namespace directmq;

class TestWriter : public portal::DataWriter {
   private:
    std::queue<std::shared_ptr<portal::Packet>>* packets;
    const size_t packetSize;
    pb_byte_t* data;
    size_t currentPosition = 0;

   public:
    TestWriter(std::queue<std::shared_ptr<portal::Packet>>* packets,
               const size_t packetSize)
        : packets(packets),
          packetSize(packetSize),
          data(new pb_byte_t[packetSize]) {}

    ~TestWriter() { delete[] data; }

    bool write(const uint8_t* block, const size_t blockSize) {
        for (size_t blockPos = 0; blockPos < blockSize; blockPos++) {
            data[this->currentPosition + blockPos] = block[blockPos];
        }

        this->currentPosition += blockSize;
        return true;
    }

    void end() {
        auto packet =
            portal::Packet::fromData(this->data, this->currentPosition);

        packets->push(packet);
    };
};
