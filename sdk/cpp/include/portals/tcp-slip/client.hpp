#pragma once

#include <iostream>
#include <boost/asio.hpp>
#include <tinyslip/slip.h>

#include "../../network/node.hpp"

namespace directmq::portal::tcp_slip::client {
    using boost::asio::ip::tcp;

    namespace internals {
        class TcpSlipWriter : public portal::DataWriter {
   private:
    tcp::socket *socket;

    uint8_t* message;
    size_t messageSize;
    size_t currentSize = 0;

   public:
    TcpSlipWriter(tcp::socket *socket,
                      size_t messageSize)
        : socket(socket),
          message(new uint8_t[messageSize]),
          messageSize(messageSize) {}

        bool write(const uint8_t* block, const size_t blockSize) override {
            if (!message) {
                return false;
            }

            if (currentSize + blockSize > messageSize) {
                return false;
            }

            std::copy(block, block + blockSize, message + currentSize);
            currentSize += blockSize;

            return true;
        }

        void end() override {
            endpoint->send(hdl, message, messageSize,
                        websocketpp::frame::opcode::BINARY);
            delete[] message;
        }
    };
    }

    class TcpSlipClient : public portal::Portal {
        private:
            network::EdgeManager* edgeManager;
            tcp::resolver resolver;
            tcp::socket socket;
            std::array<uint8_t, 1024> readBuffer;

            std::shared_ptr<network::edge::NetworkEdge> edge;

        public:
             TcpSlipClient(boost::asio::io_context& io_context, const std::string& host, const std::string& port)
        : socket(io_context), resolver(io_context) {
        auto endpoints = resolver.resolve(host, port);

        boost::asio::async_connect(socket, endpoints,
            [this](const boost::system::error_code& ec, const tcp::endpoint&) {
                if (!ec) {
                    std::cout << "Połączono z serwerem!" << std::endl;
                    doRead();
                } else {
                    std::cerr << "Błąd połączenia: " << ec.message() << std::endl;
                }
            });
        }
    };
}
