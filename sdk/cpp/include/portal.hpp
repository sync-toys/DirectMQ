#pragma once

#include <cstddef>
#include <memory>

namespace directmq::portal {
struct Packet {
   public:
    // DO NOT CALL THIS CONSTRUCTOR FROM YOUR CODE, USE empty OR fromData
    // instead!
    Packet(size_t size, uint8_t* data) {
        this->size = size;
        this->data = data;
    }

    ~Packet() {
        if (data != nullptr) {
            delete[] data;

            data = nullptr;
            size = 0;
        }
    }

    size_t size;
    uint8_t* data;

    static std::shared_ptr<Packet> empty() {
        return std::make_shared<Packet>(Packet(0, nullptr));
    }

    static std::shared_ptr<Packet> fromData(const uint8_t* data,
                                            const size_t size) {
        auto packet = std::make_shared<Packet>(size, new uint8_t[size]);
        std::copy(data, data + size, packet->data);
        return packet;
    }
};

class PacketReader {
   public:
    virtual std::shared_ptr<Packet> readPacket() = 0;
};

class DataWriter {
   public:
    virtual ~DataWriter() = default;
    virtual bool write(const uint8_t* block, const size_t blockSize) = 0;
    virtual void end() = 0;
};

class PacketWriter {
   public:
    virtual std::shared_ptr<DataWriter> beginWrite(const size_t packetSize) = 0;
};

class Closer {
   public:
    virtual void close() = 0;
};

class Portal : public PacketWriter, public Closer {};
}  // namespace directmq::portal
