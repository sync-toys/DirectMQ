#pragma once
#include <functional>

#include "protocol/decoder.hpp"

using namespace directmq;
using namespace directmq::protocol;

class TestHandler : public DecodingHandler {
   private:
    std::function<void(const messages::SupportedProtocolVersionsMessage&)>
        onSupportedProtocolVersionsCallback;

    std::function<void(const messages::InitConnectionMessage&)>
        onInitConnectionCallback;

    std::function<void(const messages::ConnectionAcceptedMessage&)>
        onConnectionAcceptedCallback;

    std::function<void(const messages::GracefullyCloseMessage&)>
        onGracefullyCloseCallback;

    std::function<void(const messages::TerminateNetworkMessage&)>
        onTerminateNetworkCallback;

    std::function<void(const messages::PublishMessage&)> onPublishCallback;

    std::function<void(const messages::SubscribeMessage&)> onSubscribeCallback;

    std::function<void(const messages::UnsubscribeMessage&)>
        onUnsubscribeCallback;

    std::function<void(const messages::MalformedMessage&)>
        onMalformedMessageCallback;

   public:
    void setOnSupportedProtocolVersions(
        std::function<void(const messages::SupportedProtocolVersionsMessage&)>
            onSupportedProtocolVersions) {
        this->onSupportedProtocolVersionsCallback = onSupportedProtocolVersions;
    }
    void onSupportedProtocolVersions(
        const messages::SupportedProtocolVersionsMessage& message) {
        if (this->onSupportedProtocolVersionsCallback == nullptr) {
            return;
        }

        this->onSupportedProtocolVersionsCallback(message);
    }

    void setOnInitConnection(
        std::function<void(const messages::InitConnectionMessage&)>
            onInitConnection) {
        this->onInitConnectionCallback = onInitConnection;
    }
    void onInitConnection(const messages::InitConnectionMessage& message) {
        if (this->onInitConnectionCallback == nullptr) {
            return;
        }

        this->onInitConnectionCallback(message);
    }

    void setOnConnectionAccepted(
        std::function<void(const messages::ConnectionAcceptedMessage&)>
            onConnectionAccepted) {
        this->onConnectionAcceptedCallback = onConnectionAccepted;
    }
    void onConnectionAccepted(
        const messages::ConnectionAcceptedMessage& message) {
        if (this->onConnectionAcceptedCallback == nullptr) {
            return;
        }

        this->onConnectionAcceptedCallback(message);
    }

    void setOnGracefullyClose(
        std::function<void(const messages::GracefullyCloseMessage&)>
            onGracefullyClose) {
        this->onGracefullyCloseCallback = onGracefullyClose;
    }
    void onGracefullyClose(const messages::GracefullyCloseMessage& message) {
        if (this->onGracefullyCloseCallback == nullptr) {
            return;
        }

        this->onGracefullyCloseCallback(message);
    }

    void setOnTerminateNetwork(
        std::function<void(const messages::TerminateNetworkMessage&)>
            onTerminateNetwork) {
        this->onTerminateNetworkCallback = onTerminateNetwork;
    }
    void onTerminateNetwork(const messages::TerminateNetworkMessage& message) {
        if (this->onTerminateNetworkCallback == nullptr) {
            return;
        }

        this->onTerminateNetworkCallback(message);
    }

    void setOnPublish(
        std::function<void(const messages::PublishMessage&)> onPublish) {
        this->onPublishCallback = onPublish;
    }
    void onPublish(const messages::PublishMessage& message) {
        if (this->onPublishCallback == nullptr) {
            return;
        }

        this->onPublishCallback(message);
    }

    void setOnSubscribe(
        std::function<void(const messages::SubscribeMessage&)> onSubscribe) {
        this->onSubscribeCallback = onSubscribe;
    }
    void onSubscribe(const messages::SubscribeMessage& message) {
        if (this->onSubscribeCallback == nullptr) {
            return;
        }

        this->onSubscribeCallback(message);
    }

    void setOnUnsubscribe(
        std::function<void(const messages::UnsubscribeMessage&)>
            onUnsubscribe) {
        this->onUnsubscribeCallback = onUnsubscribe;
    }
    void onUnsubscribe(const messages::UnsubscribeMessage& message) {
        if (this->onUnsubscribeCallback == nullptr) {
            return;
        }

        this->onUnsubscribeCallback(message);
    }

    void setOnMalformedMessage(
        std::function<void(const messages::MalformedMessage&)>
            onMalformedMessage) {
        this->onMalformedMessageCallback = onMalformedMessage;
    }
    void onMalformedMessage(const messages::MalformedMessage& message) {
        if (this->onMalformedMessageCallback == nullptr) {
            return;
        }

        this->onMalformedMessageCallback(message);
    }
};
