#pragma once

#include <boost/asio.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/io_service.hpp>
#include <boost/asio/placeholders.hpp>
#include <boost/system/detail/error_code.hpp>
#include <cstdint>
#include <functional>
#include <memory>
#include <thread>
#include <vector>

#include "../../network/node.hpp"
#include "async_stream_portal.hpp"
#include "portal.hpp"

namespace directmq::portal::streams {
class TcpPortalClient : public std::enable_shared_from_this<TcpPortalClient>,
                        public portal::Portal {
   public:
    using Pointer = std::shared_ptr<TcpPortalClient>;
    using AsyncConnectHandler =
        std::function<void(const boost::system::error_code &, Pointer)>;
    using Portal = AsyncStreamPortal<boost::asio::ip::tcp>;

    static Pointer connect(std::shared_ptr<network::EdgeManager> edgeManager,
                           const std::string &host, uint_least16_t port) {
        Pointer portal = create(edgeManager);
        auto error = portal->performConnection(host, port);

        if (error) {
            throw error;
        }

        return portal;
    }

    static void connectAsync(std::shared_ptr<network::EdgeManager> edgeManager,
                             const std::string &host, uint_least16_t port,
                             AsyncConnectHandler asyncConnectHandler) {
        Pointer portal = create(edgeManager);
        portal->performAsyncConnection(host, port, asyncConnectHandler);
    }

    std::shared_ptr<DataWriter> beginWrite(
        const std::size_t messageSize) override {
        return portal_->beginWrite(messageSize);
    }

    void wait() {
        if (alreadyStopped_ || !alreadyStarted_) {
            return;
        }

        ioThread_.join();
    }

    void close() override {
        if (alreadyStopped_) {
            return;
        }

        alreadyStopped_ = true;
        portal_->close();
        ioService_.stop();
        ioThread_.join();
    }

    Portal::Pointer portal() { return portal_; }

    ~TcpPortalClient() { close(); }

   private:
    boost::asio::io_service ioService_;
    std::thread ioThread_;

    Portal::Pointer portal_;

    std::shared_ptr<network::EdgeManager> edgeManager_;
    std::shared_ptr<network::edge::NetworkEdge> networkEdge_;

    bool alreadyStarted_ = false;
    bool alreadyStopped_ = false;

    static Pointer create(std::shared_ptr<network::EdgeManager> edgeManager) {
        return Pointer(new TcpPortalClient(edgeManager));
    }

    TcpPortalClient(std::shared_ptr<network::EdgeManager> edgeManager)
        : portal_(nullptr), networkEdge_(nullptr), edgeManager_(edgeManager) {}

    boost::asio::ip::tcp::socket *initPortal() {
        auto socket = new boost::asio::ip::tcp::socket(ioService_);
        portal_ = Portal::create(
            socket,
            [this](Portal::Pointer portal, std::vector<uint8_t> message) {
                handleMessageReceived(portal, message);
            },
            [this](Portal::Pointer portal,
                   const boost::system::error_code &error) {
                handleProcessingError(portal, error);
            });

        return socket;
    }

    boost::system::error_code performConnection(const std::string &host,
                                                uint_least16_t port) {
        auto socket = initPortal();

        boost::asio::ip::tcp::resolver resolver(ioService_);
        boost::asio::ip::tcp::resolver::query query(host, std::to_string(port));
        boost::asio::ip::tcp::resolver::iterator endpointIterator =
            resolver.resolve(query);

        boost::system::error_code error;
        boost::asio::connect(*socket, endpointIterator, error);

        if (error) {
            handleProcessingError(portal_, error);
            return error;
        }

        portal_->start();
        ioThread_ = std::thread([this]() { ioService_.run(); std::cout << "ioService_.run() finished" << std::endl; });

        this->networkEdge_ = edgeManager_->addConnectingEdge(shared_from_this());
        return boost::system::error_code();
    }

    void performAsyncConnection(const std::string &host, uint_least16_t port,
                                AsyncConnectHandler asyncConnectHandler) {
        auto socket = initPortal();

        boost::asio::ip::tcp::resolver resolver(ioService_);
        boost::asio::ip::tcp::resolver::query query(host, std::to_string(port));
        boost::asio::ip::tcp::resolver::iterator endpointIterator =
            resolver.resolve(query);

        auto self(shared_from_this());
        boost::asio::async_connect(
            *socket, endpointIterator,
            [self, asyncConnectHandler](
                const boost::system::error_code &error,
                const boost::asio::ip::tcp::resolver::iterator &) {
                if (!error) {
                    self->portal_->start();
                    self->networkEdge_ =
                        self->edgeManager_->addConnectingEdge(self);
                }
                asyncConnectHandler(error, self);
            });

        ioThread_ = std::thread([this]() { ioService_.run(); });
    }

    void handleMessageReceived(Portal::Pointer portal,
                               std::vector<uint8_t> message) {
        auto packet = portal::Packet::fromData(message.data(), message.size());
        networkEdge_->processIncomingPacket(packet);
    }

    void handleProcessingError(Portal::Pointer portal,
                               const boost::system::error_code &error) {
        portal->close();
        edgeManager_->removeEdge(portal, error.message());
    }
};
}  // namespace directmq::portal::streams
