#pragma once

#include <array>
#include <boost/asio.hpp>
#include <boost/asio/placeholders.hpp>
#include <cstdint>
#include <functional>
#include <memory>
#include <queue>
#include <vector>

namespace directmq::portal::streams {
template <typename Protocol>
class AsyncStreamFrameWriter
    : public std::enable_shared_from_this<AsyncStreamFrameWriter<Protocol>> {
public:
  using ErrorHandler = std::function<void (const boost::system::error_code &)>;
  using Socket = boost::asio::basic_stream_socket<Protocol>;
  using Pointer = std::shared_ptr<AsyncStreamFrameWriter>;

  static Pointer create(Socket &socket, ErrorHandler errorHandler) {
    return Pointer(new AsyncStreamFrameWriter(socket, errorHandler));
  }

  void send(std::vector<uint8_t> message) {
    messageQueue.push(message);

    if (!transmitting) {
      transmitting = true;
      startTransmission();
    }
  }

private:
  Socket &socket;
  ErrorHandler errorHandler;
  std::queue<std::vector<uint8_t>> messageQueue;
  bool transmitting = false;

  AsyncStreamFrameWriter(Socket &socket, ErrorHandler &errorHandler)
      : socket(socket), errorHandler(errorHandler) {}

  std::array<uint8_t, 4> encodeFrameHeader(std::size_t size) {
    std::array<uint8_t, 4> header;
    header[0] = (size >> 24) & 0xFF;
    header[1] = (size >> 16) & 0xFF;
    header[2] = (size >> 8) & 0xFF;
    header[3] = size & 0xFF;
    return header;
  }

  std::vector<uint8_t> encodeFrame(std::vector<uint8_t> message) {
    std::vector<uint8_t> encoded;

    auto header = encodeFrameHeader(message.size());
    encoded.resize(header.size() + message.size());

    std::copy(header.begin(), header.end(), encoded.begin());
    std::copy(message.begin(), message.end(), encoded.begin() + header.size());

    return encoded;
  }

  void startTransmission() {
    if (messageQueue.empty()) {
      transmitting = false;
      return;
    }

    auto message = messageQueue.front();
    messageQueue.pop();

    auto encoded = encodeFrame(message);

    boost::asio::async_write(
        socket, boost::asio::buffer(encoded),
        std::bind(&AsyncStreamFrameWriter::handleTransmission,
                  this->shared_from_this(), boost::asio::placeholders::error));
  }

  void handleTransmission(const boost::system::error_code &error) {
    if (error) {
      errorHandler(error);
      return;
    }

    startTransmission();
  }
};
} // namespace directmq::portal::streams
