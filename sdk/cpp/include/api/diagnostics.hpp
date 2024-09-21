#pragma once

#include <functional>
#include <string>

#include "../network/participant.hpp"
#include "../portal.hpp"
#include "../protocol/messages.hpp"


// break the circular dependency
namespace directmq::network::edge {
class NetworkEdge;
}  // namespace directmq::network::edge

namespace directmq::api {


using ConnectionEstablishedHandler = void(const std::string& bridgedNodeID,
                                          portal::Portal& portal);

using ConnectionLostHandler = void(const std::string& bridgedNodeID,
                                   const std::string& reason,
                                   portal::Portal& portal);

using EdgeDisconnectedHandler = void(const std::string& bridgedNodeID,
                                     const std::string& reason,
                                     network::edge::NetworkEdge& edge,
                                     portal::Portal& portal);

using PublicationHandler = void(const protocol::messages::PublishMessage& publication);

using SubscriptionHandler = void(const protocol::messages::SubscribeMessage& subscription);

using UnsubscribeHandler = void(const protocol::messages::UnsubscribeMessage& unsubscription);

using TerminateNetworkHandler = void(const protocol::messages::TerminateNetworkMessage& termination);


class DiagnosticsAPI {
   public:
    virtual ~DiagnosticsAPI() = default;

    virtual void setOnConnectionEstablishedHandler(
        std::function<ConnectionEstablishedHandler>
            handler) = 0;

    virtual void setOnEdgeDisconnectionHandler(
        std::function<EdgeDisconnectedHandler>
            handler) = 0;

    virtual void setOnPublicationHandler(
        std::function<PublicationHandler>
            handler) = 0;

    virtual void setOnSubscriptionHandler(
        std::function<SubscriptionHandler>
            handler) = 0;

    virtual void setOnUnsubscribeHandler(
        std::function<UnsubscribeHandler>
            handler) = 0;

    virtual void setOnTerminateNetworkHandler(
        std::function<TerminateNetworkHandler>
            handler) = 0;
};

class DiagnosticsAPIAdapter : public DiagnosticsAPI, public network::NetworkParticipant {
   public:
    virtual ~DiagnosticsAPIAdapter() = default;

    // required for NetworkNode setup
    virtual void setOnConnectionLostHandler(
        std::function<ConnectionLostHandler>
            handler) = 0;

    virtual void onConnectionEstablished(const std::string& bridgedNodeID,
                                         directmq::portal::Portal& portal) = 0;

    virtual void onConnectionLost(const std::string& bridgedNodeID,
                                  const std::string& reason,
                                  directmq::portal::Portal& portal) = 0;

    virtual void onEdgeDisconnection(const std::string& bridgedNodeID,
                                     const std::string& reason,
                                     directmq::network::edge::NetworkEdge& edge,
                                     directmq::portal::Portal& portal) = 0;
};

class DiagnosticsAPILambdaAdapter : public DiagnosticsAPIAdapter {
   private:
    std::function<ConnectionEstablishedHandler>
        onConnectionEstablishedHandler;

    std::function<ConnectionLostHandler>
        onConnectionLostHandler;

    std::function<EdgeDisconnectedHandler>
        onEdgeDisconnectionHandler;

    std::function<PublicationHandler>
        onPublicationHandler;

    std::function<SubscriptionHandler>
        onSubscriptionHandler;

    std::function<UnsubscribeHandler>
        onUnsubscribeHandler;

    std::function<TerminateNetworkHandler>
        onTerminateNetworkHandler;

   public:
    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return {};
    }

    bool willHandleTopic(const std::string& topic) const override {
        return false;
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        return false;
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        return false;
    }

    void setOnConnectionEstablishedHandler(
        std::function<ConnectionEstablishedHandler>
            handler) override {
        onConnectionEstablishedHandler = handler;
    }

    void onConnectionEstablished(const std::string& bridgedNodeID,
                                 directmq::portal::Portal& portal) override {
        if (onConnectionEstablishedHandler) {
            onConnectionEstablishedHandler(bridgedNodeID, portal);
        }
    }

    void setOnConnectionLostHandler(
        std::function<ConnectionLostHandler>
            handler) override {
        onConnectionLostHandler = handler;
    }

    void onConnectionLost(const std::string& bridgedNodeID,
                          const std::string& reason,
                          directmq::portal::Portal& portal) override {
        if (onConnectionLostHandler) {
            onConnectionLostHandler(bridgedNodeID, reason, portal);
        }
    }

    void setOnEdgeDisconnectionHandler(
        std::function<EdgeDisconnectedHandler>
            handler) override {
        onEdgeDisconnectionHandler = handler;
    }

    void onEdgeDisconnection(const std::string& bridgedNodeID,
                             const std::string& reason,
                             directmq::network::edge::NetworkEdge& edge,
                             directmq::portal::Portal& portal) override {
        if (onEdgeDisconnectionHandler) {
            onEdgeDisconnectionHandler(bridgedNodeID, reason, edge, portal);
        }
    }

    void setOnPublicationHandler(
        std::function<PublicationHandler>
            handler) override {
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
        std::function<SubscriptionHandler>
            handler) override {
        onSubscriptionHandler = handler;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        if (onSubscriptionHandler) {
            onSubscriptionHandler(subscription);
        }
    }

    void setOnUnsubscribeHandler(
        std::function<UnsubscribeHandler>
            handler) override {
        onUnsubscribeHandler = handler;
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        if (onUnsubscribeHandler) {
            onUnsubscribeHandler(unsubscription);
        }
    }

    void setOnTerminateNetworkHandler(
        std::function<TerminateNetworkHandler>
            handler) override {
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
