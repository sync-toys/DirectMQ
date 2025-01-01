#pragma once

#include <algorithm>
#include <boost/asio.hpp>
#include <boost/asio/placeholders.hpp>
#include <cstdint>
#include <functional>
#include <memory>
#include <vector>

#include "async_stream_frame_reader.hpp"
#include "async_stream_frame_writer.hpp"
#include "portal.hpp"

namespace directmq::portal::streams {
template <typename Protocol>
class AsyncStreamPortal
    : public portal::Portal,
      public std::enable_shared_from_this<AsyncStreamPortal<Protocol>> {
   public:
    using Pointer = std::shared_ptr<AsyncStreamPortal>;
    using ReceiveHandler = std::function<void(Pointer, std::vector<uint8_t>)>;
    using ErrorHandler =
        std::function<void(Pointer, const boost::system::error_code &)>;
    using Socket = boost::asio::basic_stream_socket<Protocol>;

    class AsyncDataWriter : public portal::DataWriter {
       private:
        std::vector<uint8_t> data;
        std::size_t currentOffset = 0;
        std::size_t messageSize;
        typename AsyncStreamFrameWriter<Protocol>::Pointer frameWriter_;

       public:
        using Pointer = std::shared_ptr<AsyncDataWriter>;

        AsyncDataWriter(
            std::size_t messageSize,
            typename AsyncStreamFrameWriter<Protocol>::Pointer frameWriter)
            : messageSize(messageSize), frameWriter_(frameWriter) {
            data.resize(messageSize);
        }

        bool write(const uint8_t *block, const std::size_t blockSize) override {
            if (currentOffset + blockSize > messageSize) {
                return false;
            }

            std::copy(block, block + blockSize, data.begin() + currentOffset);
            currentOffset += blockSize;

            return true;
        }

        void end() override { frameWriter_->send(data); }
    };

    static Pointer create(Socket *socket, ReceiveHandler receiveHandler,
                          ErrorHandler errorHandler) {
        return Pointer(
            new AsyncStreamPortal(socket, receiveHandler, errorHandler));
    }

    void start() {
        if (alreadyStarted_) {
            // TODO: exception?
            return;
        }

        alreadyStarted_ = true;
        initComponents();

        frameReader_->start();
    }

    std::shared_ptr<DataWriter> beginWrite(
        const std::size_t messageSize) override {
        if (!alreadyStarted_) {
            throw std::runtime_error(
                "cannot write to portal, portal has not been started");
        }

        if (errorOcurred_) {
            throw std::runtime_error("cannot write to portal, error ocurred");
        }

        if (isClosed_) {
            throw std::runtime_error(
                "cannot write to portal, portal is closed");
        }

        return AsyncDataWriter::Pointer(
            new AsyncDataWriter(messageSize, frameWriter_));
    }

    void close() override {
        if (isClosed_) {
            return;
        }

        isClosed_ = true;
        socket_->close();
        delete socket_;
    }

   private:
    Socket *socket_;
    ReceiveHandler receiveHandler_;
    ErrorHandler errorHandler_;

    bool alreadyStarted_ = false;
    bool errorOcurred_ = false;
    bool isClosed_ = false;

    typename AsyncStreamFrameReader<Protocol>::Pointer frameReader_;
    typename AsyncStreamFrameWriter<Protocol>::Pointer frameWriter_;

    AsyncStreamPortal(Socket *socket, ReceiveHandler receiveHandler,
                      ErrorHandler errorHandler)
        : socket_(socket),
          receiveHandler_(receiveHandler),
          errorHandler_(errorHandler) {}

    void initComponents() {
        auto receiveHandler =
            std::bind(&AsyncStreamPortal::handleMessage,
                      this->shared_from_this(), std::placeholders::_1);
        auto errorHandler =
            std::bind(&AsyncStreamPortal::handleError, this->shared_from_this(),
                      std::placeholders::_1);

        frameReader_ = AsyncStreamFrameReader<Protocol>::create(
            *socket_, receiveHandler, errorHandler);
        frameWriter_ =
            AsyncStreamFrameWriter<Protocol>::create(*socket_, errorHandler);
    }

    void handleMessage(std::vector<uint8_t> message) {
        receiveHandler_(this->shared_from_this(), message);
    }

    void handleError(const boost::system::error_code &error) {
        errorOcurred_ = true;
        errorHandler_(this->shared_from_this(), error);
    }
};
}  // namespace directmq::portal::streams
