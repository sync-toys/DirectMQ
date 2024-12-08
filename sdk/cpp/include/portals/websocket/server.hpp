#pragma once

/**
 * DISCLAIMER: this C++ SDK Websocket portal is not working at the moment.
 *
 * This is work in progress implementation, TBD.
 */

#include <iostream>
#include <memory>
#include <string>
#include <vector>
#include <websocketpp/config/asio_no_tls.hpp>
#include <websocketpp/server.hpp>

#include "../../network/node.hpp"

namespace directmq::portal::websocket::server {
typedef websocketpp::server<websocketpp::config::asio> server;

using websocketpp::lib::bind;
using websocketpp::lib::placeholders::_1;
using websocketpp::lib::placeholders::_2;

typedef websocketpp::config::asio::message_type::ptr message_ptr;

namespace internals {
class WebsocketppWriter : public portal::DataWriter {
   private:
    websocketpp::connection_hdl hdl;
    server* srv;

    uint8_t* message;
    size_t messageSize;
    size_t currentSize = 0;

   public:
    WebsocketppWriter(websocketpp::connection_hdl hdl, server* srv,
                      size_t messageSize)
        : hdl(hdl),
          srv(srv),
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
        srv->send(hdl, message, messageSize,
                  websocketpp::frame::opcode::BINARY);
        delete[] message;
    }
};

class WebsocketConnection : public portal::Portal {
   private:
    websocketpp::connection_hdl hdl;
    server* srv;
    std::shared_ptr<network::edge::NetworkEdge> edge;

   public:
    WebsocketConnection(websocketpp::connection_hdl hdl, server* srv)
        : hdl(hdl), srv(srv) {}

    void setEdge(std::shared_ptr<network::edge::NetworkEdge> edge) {
        this->edge = edge;
    }

    websocketpp::connection_hdl getHdl() { return hdl; }

    void onMessage(websocketpp::connection_hdl hdl, message_ptr msg) {
        if (msg->get_opcode() == websocketpp::frame::opcode::TEXT) {
            return;
        }

        this->edge->processIncomingPacket(Packet::fromData(
            (uint8_t*)msg->get_payload().data(), msg->get_payload().size()));
    }

    std::shared_ptr<DataWriter> beginWrite(const size_t packetSize) override {
        return std::make_shared<WebsocketppWriter>(hdl, srv, packetSize);
    }

    void close() override {
        srv->close(hdl, websocketpp::close::status::going_away, "");
    }
};
}  // namespace internals

class WebsocketServer {
   private:
    server srv;
    network::EdgeManager* edgeManager;

    std::vector<std::shared_ptr<internals::WebsocketConnection>> connections;

    void onOpen(websocketpp::connection_hdl hdl) {
        std::shared_ptr<internals::WebsocketConnection> connection =
            std::make_shared<internals::WebsocketConnection>(hdl, &srv);

        // TODO: fix another memory leak there, connection is included in edge
        // and edge is included in connection (use weak_ptr)
        auto edge = edgeManager->addListeningEdge(connection);
        connection->setEdge(edge);

        connections.push_back(connection);
    }

    void onMessage(websocketpp::connection_hdl hdl, message_ptr msg) {
        auto connection = std::find_if(
            connections.begin(), connections.end(),
            [&](std::shared_ptr<internals::WebsocketConnection> connection) {
                return connection->getHdl().lock() == hdl.lock();
            });

        if (connection != connections.end()) {
            (*connection)->onMessage(hdl, msg);
        }
    }

    void onClose(websocketpp::connection_hdl hdl) {
        auto connection = std::find_if(
            connections.begin(), connections.end(),
            [&](std::shared_ptr<internals::WebsocketConnection> connection) {
                return connection->getHdl().lock() == hdl.lock();
            });

        if (connection != connections.end()) {
            edgeManager->removeEdge(*connection, "websocket connection closed");
            connections.erase(connection);
        }
    }

   public:
    WebsocketServer(network::EdgeManager* edgeManager)
        : edgeManager(edgeManager) {}

    void run(uint16_t port) {
        srv.set_open_handler(bind(&WebsocketServer::onOpen, this, _1));
        srv.set_close_handler(bind(&WebsocketServer::onClose, this, _1));
        srv.set_message_handler(
            bind(&WebsocketServer::onMessage, this, _1, _2));
        srv.init_asio();
        srv.listen(port);
        srv.start_perpetual();
        srv.run();
    }
};
}  // namespace directmq::portal::websocket::server
