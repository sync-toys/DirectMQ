#pragma once

#include <algorithm>
#include <boost/asio.hpp>
#include <boost/asio/io_service.hpp>
#include <boost/asio/placeholders.hpp>
#include <cstdint>
#include <memory>
#include <thread>
#include <vector>

#include "async_stream_portal.hpp"

#include "../../network/node.hpp"

namespace directmq::portal::streams {
class TcpPortalServerConnection {
    public:

};

class TcpPortalServer : public std::enable_shared_from_this<TcpPortalServer> {
public:
  using Pointer = std::shared_ptr<TcpPortalServer>;

  static Pointer create(std::shared_ptr<network::EdgeManager> edgeManager, uint_least16_t port, std::size_t maxConnections) {
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

  boost::asio::io_service ioService_;
  std::thread ioThread_;

  boost::asio::ip::tcp::acceptor acceptor_;
  std::vector<Portal::Pointer> portals_;
  std::size_t maxConnections_;

  std::shared_ptr<network::EdgeManager> edgeManager_;

  bool alreadyStarted_ = false;
  bool alreadyStopped_ = false;

  TcpPortalServer(std::shared_ptr<network::EdgeManager> edgeManager, uint_least16_t port, std::size_t maxConnections)
      : edgeManager_(edgeManager), acceptor_(ioService_, boost::asio::ip::tcp::endpoint(
                                  boost::asio::ip::tcp::v4(), port)),
        maxConnections_(maxConnections) {}

  void acceptConnection() {
    Portal::Socket *socket = new boost::asio::ip::tcp::socket(ioService_);

    auto newPortal = Portal::create(
        socket,
        [this](Portal::Pointer portal, std::vector<uint8_t> message) {
          handleMessageReceived(portal, message);
        },
        [this](Portal::Pointer portal, const boost::system::error_code &error) {
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

    portals_.push_back(newPortal);
    newPortal->start();

    edgeManager_->addListeningEdge(newPortal);
  }

  void handleMessageReceived(Portal::Pointer portal,
                             std::vector<uint8_t> message) {

                                // TODO: write TcpPortalServerConnection,
                                // that will contain Portal::Pointer and networkEdge from
                                // addPortal edgeManager_->addListeningEdge,
                                // then rewrite handleConnectionError and closeAllPortals
                                // to remove items correctly
                                //
                                // or maybe AsyncStreamPortal should contain
                                // reference to networkEdge? but only if networkEdge does
                                // not have reference to AsyncStreamPortal (weak_ptr as solution?)

        auto packet = portal::Packet::fromData(message.data(), message.size());
        networkEdge_->processIncomingPacket(packet);
  }

  void handleConnectionError(Portal::Pointer portal,
                              const boost::system::error_code &error) {
        edgeManager_->removeEdge(portal, error.message());
        portals_.erase(std::remove(portals_.begin(), portals_.end(), portal), portals_.end());
  }

  void closeAllPortals() {
    for (auto portal : portals_) {
      edgeManager_->removeEdge(portal, "server closed");
      portal->close();
    }

    portals_.clear();
  }
};
} // namespace directmq::portal::streams
