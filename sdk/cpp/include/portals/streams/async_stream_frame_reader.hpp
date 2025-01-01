#pragma once

#include <array>
#include <boost/asio.hpp>
#include <boost/asio/placeholders.hpp>
#include <cstdint>
#include <functional>
#include <memory>
#include <vector>

namespace directmq::portal::streams {
template <typename Protocol>
class AsyncStreamFrameReader
    : public std::enable_shared_from_this<AsyncStreamFrameReader<Protocol>> {
   public:
    using ReceiveHandler = std::function<void(std::vector<uint8_t>)>;
    using ErrorHandler = std::function<void(const boost::system::error_code &)>;
    using Socket = boost::asio::basic_stream_socket<Protocol>;
    using Pointer = std::shared_ptr<AsyncStreamFrameReader>;

    static Pointer create(Socket &socket, ReceiveHandler receiveHandler,
                          ErrorHandler errorHandler) {
        return Pointer(
            new AsyncStreamFrameReader(socket, receiveHandler, errorHandler));
    }

    void start() { readFrameHeader(); }

   private:
    Socket &socket;
    ReceiveHandler receiveHandler;
    ErrorHandler errorHandler;

    std::array<uint8_t, 4> header;
    std::vector<uint8_t> body;

    AsyncStreamFrameReader(Socket &socket, ReceiveHandler &receiveHandler,
                           ErrorHandler &errorHandler)
        : socket(socket),
          receiveHandler(receiveHandler),
          errorHandler(errorHandler) {}

    void readFrameHeader() {
        boost::asio::async_read(
            socket, boost::asio::buffer(header),
            std::bind(&AsyncStreamFrameReader::readFrameBody,
                      this->shared_from_this(),
                      boost::asio::placeholders::error));
    }

    void readFrameBody(const boost::system::error_code &error) {
        if (error) {
            errorHandler(error);
            return;
        }

        std::size_t size = decodeHeader(header);
        body.resize(size);

        boost::asio::async_read(
            socket, boost::asio::buffer(body),
            std::bind(&AsyncStreamFrameReader::handleFrameBodyReceived,
                      this->shared_from_this(),
                      boost::asio::placeholders::error));
    }

    void handleFrameBodyReceived(const boost::system::error_code &error) {
        if (error) {
            errorHandler(error);
            return;
        }

        receiveHandler(body);
        readFrameHeader();
    }

    std::size_t decodeHeader(std::array<uint8_t, 4> header) {
        constexpr uint32_t SHIFT_24 = 24;
        constexpr uint32_t SHIFT_16 = 16;
        constexpr uint32_t SHIFT_8 = 8;

        return static_cast<uint32_t>(header[0]) << SHIFT_24 |
               static_cast<uint32_t>(header[1]) << SHIFT_16 |
               static_cast<uint32_t>(header[2]) << SHIFT_8 |
               static_cast<uint32_t>(header[3]);
    }
};
}  // namespace directmq::portal::streams
