#pragma once

#include <functional>
#include <string>

#include "../network/participant.hpp"
#include "../portal.hpp"
#include "../protocol/messages.hpp"

namespace directmq::api {
class DiagnosticsAPI : public network::NetworkParticipant {
   public:
    virtual ~DiagnosticsAPI() = default;

    virtual void onConnectionEstablished(const std::string& bridgedNodeID,
                                         directmq::portal::Portal& portal) = 0;

    virtual void onConnectionLost(const std::string& bridgedNodeID,
                                  const std::string& reason,
                                  directmq::portal::Portal& portal) = 0;

    // required for NetworkNode setup
    virtual void setOnConnectionLostHandler(
        std::function<void(const std::string& bridgedNodeID,
                           const std::string& reason,
                           directmq::portal::Portal& portal)>
            handler) = 0;
};

class DiagnosticsAPILambdaHandler : public DiagnosticsAPI {
   private:
    std::function<void(const std::string& bridgedNodeID,
                       directmq::portal::Portal& portal)>
        onConnectionEstablishedHandler;

    std::function<void(const std::string& bridgedNodeID,
                       const std::string& reason,
                       directmq::portal::Portal& portal)>
        onConnectionLostHandler;

    std::function<void(const protocol::messages::PublishMessage& publication)>
        onPublicationHandler;

    std::function<void(
        const protocol::messages::SubscribeMessage& subscription)>
        onSubscriptionHandler;

    std::function<void(
        const protocol::messages::UnsubscribeMessage& unsubscription)>
        onUnsubscribeHandler;

    std::function<void(
        const protocol::messages::TerminateNetworkMessage& termination)>
        onTerminateNetworkHandler;

   public:
    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return {};
    }

    bool willHandleTopic(const std::string& topic) const override {
        return false;
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        return false;
    }

    void setOnConnectionEstablishedHandler(
        std::function<void(const std::string& bridgedNodeID,
                           directmq::portal::Portal& portal)>
            handler) {
        onConnectionEstablishedHandler = handler;
    }

    void onConnectionEstablished(const std::string& bridgedNodeID,
                                 directmq::portal::Portal& portal) override {
        if (onConnectionEstablishedHandler) {
            onConnectionEstablishedHandler(bridgedNodeID, portal);
        }
    }

    void setOnConnectionLostHandler(
        std::function<void(const std::string& bridgedNodeID,
                           const std::string& reason,
                           directmq::portal::Portal& portal)>
            handler) {
        onConnectionLostHandler = handler;
    }

    void onConnectionLost(const std::string& bridgedNodeID,
                          const std::string& reason,
                          directmq::portal::Portal& portal) override {
        if (onConnectionLostHandler) {
            onConnectionLostHandler(bridgedNodeID, reason, portal);
        }
    }

    void setOnPublicationHandler(
        std::function<
            void(const protocol::messages::PublishMessage& publication)>
            handler) {
        onPublicationHandler = handler;
    }

    bool handlePublish(
        const protocol::messages::PublishMessage& publication) override {
        if (onPublicationHandler) {
            onPublicationHandler(publication);
        }

        // this is diagnostics API, so it should
        // never handle any publications on its own
        return false;
    }

    void setOnSubscriptionHandler(
        std::function<
            void(const protocol::messages::SubscribeMessage& subscription)>
            handler) {
        onSubscriptionHandler = handler;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        if (onSubscriptionHandler) {
            onSubscriptionHandler(subscription);
        }
    }

    void setOnUnsubscribeHandler(
        std::function<
            void(const protocol::messages::UnsubscribeMessage& unsubscription)>
            handler) {
        onUnsubscribeHandler = handler;
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        if (onUnsubscribeHandler) {
            onUnsubscribeHandler(unsubscription);
        }
    }

    void setOnTerminateNetworkHandler(
        std::function<void(
            const protocol::messages::TerminateNetworkMessage& termination)>
            handler) {
        onTerminateNetworkHandler = handler;
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        if (onTerminateNetworkHandler) {
            onTerminateNetworkHandler(termination);
        }
    }
};
}  // namespace directmq::api
