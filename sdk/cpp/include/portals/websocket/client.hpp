#pragma once

#include <iostream>
#include <memory>
#include <string>
#include <websocketpp/client.hpp>
#include <websocketpp/config/asio_no_tls_client.hpp>

#include "../../../network/node.hpp"

namespace directmq::portal::websocket::client {
typedef websocketpp::client<websocketpp::config::asio_client> client;

using websocketpp::lib::bind;
using websocketpp::lib::placeholders::_1;
using websocketpp::lib::placeholders::_2;

typedef websocketpp::config::asio_client::message_type::ptr message_ptr;

class WebsocketpptWriter : public portal::DataWriter {
   private:
    websocketpp::connection_hdl hdl;
    client* endpoint;

    uint8_t* message;
    size_t messageSize;
    size_t currentSize = 0;

   public:
    WebsocketpptWriter(websocketpp::connection_hdl hdl, client* endpoint,
                       size_t messageSize)
        : hdl(hdl),
          endpoint(endpoint),
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

template <typename T>
struct WebsocketCreationResult {
    std::shared_ptr<T> client;
    std::string error;
};

class WebsocketClient : public portal::Portal {
   private:
    client endpoint;
    websocketpp::connection_hdl hdl;
    network::EdgeManager* edgeManager;
    std::shared_ptr<network::edge::NetworkEdge> edge;

    // TODO: fix memory leak where onOpen and onClose have to have a shared
    // pointer to WebsocketClient so WebsocketClient will never be deleted as
    // its methods have shared pointer of itself binded by parameters (loop)

    void onOpen(std::shared_ptr<portal::Portal> thisPtr,
                websocketpp::connection_hdl hdl) {
        this->hdl = hdl;

        this->edge = this->edgeManager->addConnectingEdge(thisPtr);
    }

    void onMessage(websocketpp::connection_hdl hdl, message_ptr msg) {
        if (msg->get_opcode() == websocketpp::frame::opcode::TEXT) {
            return;
        }

        this->edge->processIncomingPacket(Packet::fromData(
            (uint8_t*)msg->get_payload().data(), msg->get_payload().size()));
    }

    void onClose(std::shared_ptr<portal::Portal> thisPtr,
                 websocketpp::connection_hdl hdl) {
        this->edgeManager->removeEdge(thisPtr, "websocket connection closed");
    }

   public:
    WebsocketClient(network::EdgeManager* edgeManager)
        : edgeManager(edgeManager) {}

    ~WebsocketClient() { close(); }

    std::shared_ptr<DataWriter> beginWrite(const size_t packetSize) override {
        return std::make_shared<WebsocketpptWriter>(hdl, &endpoint, packetSize);
    }

    void close() override {
        endpoint.close(hdl, websocketpp::close::status::going_away, "");
    }

    static WebsocketCreationResult<WebsocketClient> create(
        network::EdgeManager* edgeManager, const std::string& uri) {
        std::shared_ptr<WebsocketClient> client =
            std::make_shared<WebsocketClient>(edgeManager);

        client->endpoint.init_asio();

        client->endpoint.set_open_handler(
            bind(&WebsocketClient::onOpen, client.get(), client, _1));

        client->endpoint.set_message_handler(
            bind(&WebsocketClient::onMessage, client.get(), _1, _2));

        client->endpoint.set_close_handler(
            bind(&WebsocketClient::onClose, client.get(), client, _1));

        websocketpp::lib::error_code errorCode;
        client::connection_ptr con =
            client->endpoint.get_connection(uri, errorCode);
        if (errorCode) {
            return WebsocketCreationResult<WebsocketClient>{
                nullptr, errorCode.message()};
        }

        client->endpoint.connect(con);

        return WebsocketCreationResult<WebsocketClient>{client, ""};
    }

    void run() { endpoint.run(); }
};
}  // namespace directmq::portal::websocket::client
