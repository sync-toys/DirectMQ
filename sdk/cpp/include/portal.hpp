#pragma once

#include <bytes.hpp>
#include <cstddef>

namespace directmq::portal {
struct Packet {
    size_t size;
    bytes data;
};

class PacketReader {
   public:
    virtual Packet readPacket() = 0;
};

class DataWriter {
   public:
    virtual ~DataWriter() = default;
    virtual bool write(bytes block, const size_t blockSize) = 0;
    virtual void end() = 0;
};

class PacketWriter {
   public:
    virtual DataWriter* beginWrite(const size_t packetSize) = 0;
};

class Closer {
   public:
    virtual void close() = 0;
};

class Portal : public PacketReader, public PacketWriter, public Closer {};
}  // namespace directmq::portal
