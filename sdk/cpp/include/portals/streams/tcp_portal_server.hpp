#pragma once

#include <algorithm>
#include <boost/asio.hpp>
#include <boost/asio/io_service.hpp>
#include <boost/asio/placeholders.hpp>
#include <cstdint>
#include <memory>
#include <thread>
#include <vector>

#include "../../network/node.hpp"
#include "async_stream_portal.hpp"

namespace directmq::portal::streams {
class TcpPortalServer : public std::enable_shared_from_this<TcpPortalServer> {
   public:
    using Pointer = std::shared_ptr<TcpPortalServer>;

    static Pointer create(std::shared_ptr<network::EdgeManager> edgeManager,
                          uint_least16_t port, std::size_t maxConnections) {
        return Pointer(new TcpPortalServer(edgeManager, port, maxConnections));
    }

    void start() {
        if (alreadyStarted_) {
            return;
        }

        alreadyStarted_ = true;
        ioThread_ = std::thread([this]() { ioService_.run(); });
        acceptConnection();
    }

    void stop() {
        if (alreadyStopped_) {
            return;
        }

        alreadyStopped_ = true;

        closeAllPortals();
        ioService_.stop();
        ioThread_.join();
    }

    ~TcpPortalServer() { stop(); }

   private:
    using Portal = AsyncStreamPortal<boost::asio::ip::tcp>;

    struct TcpPortalServerConnection {
        Portal::Pointer portal;
        std::shared_ptr<network::edge::NetworkEdge> edge;
    };

    boost::asio::io_service ioService_;
    std::thread ioThread_;

    boost::asio::ip::tcp::acceptor acceptor_;
    std::vector<TcpPortalServerConnection> portals_;
    std::size_t maxConnections_;

    std::shared_ptr<network::EdgeManager> edgeManager_;

    bool alreadyStarted_ = false;
    bool alreadyStopped_ = false;

    TcpPortalServer(std::shared_ptr<network::EdgeManager> edgeManager,
                    uint_least16_t port, std::size_t maxConnections)
        : edgeManager_(edgeManager),
          acceptor_(ioService_, boost::asio::ip::tcp::endpoint(
                                    boost::asio::ip::tcp::v4(), port)),
          maxConnections_(maxConnections) {}

    void acceptConnection() {
        Portal::Socket *socket = new boost::asio::ip::tcp::socket(ioService_);

        auto newPortal = Portal::create(
            socket,
            [this](Portal::Pointer portal, std::vector<uint8_t> message) {
                handleMessageReceived(portal, message);
            },
            [this](Portal::Pointer portal,
                   const boost::system::error_code &error) {
                handleConnectionError(portal, error);
            });

        acceptor_.async_accept(
            *socket, [this, newPortal](const boost::system::error_code &error) {
                if (!error) {
                    addPortal(newPortal);
                }

                acceptConnection();
            });
    }

    void addPortal(Portal::Pointer newPortal) {
        if (maxConnections_ != 0 && portals_.size() >= maxConnections_) {
            newPortal->close();
            return;
        }

        auto edge = edgeManager_->addListeningEdge(newPortal);
        portals_.push_back({newPortal, edge});

        newPortal->start();
    }

    void handleMessageReceived(Portal::Pointer portal,
                               std::vector<uint8_t> message) {
        auto connection =
            std::find_if(portals_.begin(), portals_.end(),
                         [portal](const TcpPortalServerConnection &connection) {
                             return connection.portal == portal;
                         });

        if (connection == portals_.end()) {
            return;
        }

        auto packet = portal::Packet::fromData(message.data(), message.size());
        connection->edge->processIncomingPacket(packet);
    }

    void handleConnectionError(Portal::Pointer portal,
                               const boost::system::error_code &error) {
        edgeManager_->removeEdge(portal, error.message());
        portals_.erase(std::remove(portals_.begin(), portals_.end(), portal),
                       portals_.end());
    }

    void closeAllPortals() {
        for (auto conn : portals_) {
            edgeManager_->removeEdge(conn.portal, "server closed");
            conn.portal->close();
        }

        portals_.clear();
    }
};
}  // namespace directmq::portal::streams
