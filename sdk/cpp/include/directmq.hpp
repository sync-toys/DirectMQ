#pragma once

#include "api/diagnostics.hpp"
#include "api/native.hpp"
#include "network/edge/network_edge.hpp"
#include "network/node.hpp"
#include "portal.hpp"
#include "protocol/embedded/decoder.hpp"
#include "protocol/embedded/encoder.hpp"
#include "protocol/messages/publish.hpp"

namespace directmq {
class DirectMQNode : public network::NodeManager,
                     public network::EdgeManager,
                     public api::NativeLambdaAPI,
                     public api::DiagnosticsAPI {
   private:
    std::shared_ptr<network::NetworkNode> node;
    std::shared_ptr<api::NativeLambdaAPIAdapter> nativeAPI;
    std::shared_ptr<api::DiagnosticsAPILambdaAdapter> diagnosticsAPI;
    std::shared_ptr<protocol::embedded::EmbeddedProtocolDecoderImplementation>
        decoder;
    std::shared_ptr<protocol::embedded::EmbeddedProtocolEncoderImplementation>
        encoder;

   public:
    DirectMQNode(const network::NetworkNodeConfig& config)
        : node(std::make_shared<network::NetworkNode>(
              config, nativeAPI, diagnosticsAPI, encoder, decoder)),
          nativeAPI(std::make_shared<api::NativeLambdaAPIAdapter>()),
          diagnosticsAPI(std::make_shared<api::DiagnosticsAPILambdaAdapter>()),
          decoder(std::make_shared<
                  protocol::embedded::EmbeddedProtocolDecoderImplementation>()),
          encoder(
              std::make_shared<protocol::embedded::
                                   EmbeddedProtocolEncoderImplementation>()) {}

    ~DirectMQNode() {
        node->closeNode("node is being destroyed", []() {});
    }

    /**
     * NodeManager interface implementation
     */

    void closeNode(const std::string& reason,
                   std::function<void()> onNodeFullyClosed) override {
        node->closeNode(reason, onNodeFullyClosed);
    }

    /**
     * EdgeManager interface implementation
     */

    std::shared_ptr<network::edge::NetworkEdge> addListeningEdge(
        std::shared_ptr<portal::Portal> portal) override {
        return node->addListeningEdge(portal);
    }

    std::shared_ptr<network::edge::NetworkEdge> addConnectingEdge(
        std::shared_ptr<portal::Portal> portal) override {
        return node->addConnectingEdge(portal);
    }

    bool removeEdge(std::shared_ptr<portal::Portal> portal,
                    const std::string& reason) override {
        return node->removeEdge(portal, reason);
    }

    /**
     * NativeLambdaAPI interface implementation
     */

    void publish(
        const std::string& topic, const std::vector<uint8_t>& message,
        const protocol::messages::DeliveryStrategy& deliveryStrategy) override {
        nativeAPI->publish(topic, message, deliveryStrategy);
    }

    int subscribe(const std::string& topic,
                  std::function<api::LambdaMessageHandler> handler) override {
        return nativeAPI->subscribe(topic, handler);
    }

    void unsubscribe(int subscriptionId) override {
        nativeAPI->unsubscribe(subscriptionId);
    }

    /**
     * DiagnosticsAPI interface implementation
     */

    void setOnConnectionEstablishedHandler(
        std::function<api::ConnectionEstablishedHandler> handler) override {
        diagnosticsAPI->setOnConnectionEstablishedHandler(handler);
    }

    void setOnEdgeDisconnectionHandler(
        std::function<api::EdgeDisconnectedHandler> handler) override {
        diagnosticsAPI->setOnEdgeDisconnectionHandler(handler);
    }

    void setOnPublicationHandler(
        std::function<api::PublicationHandler> handler) override {
        diagnosticsAPI->setOnPublicationHandler(handler);
    }

    void setOnSubscriptionHandler(
        std::function<api::SubscriptionHandler> handler) override {
        diagnosticsAPI->setOnSubscriptionHandler(handler);
    }

    void setOnUnsubscribeHandler(
        std::function<api::UnsubscribeHandler> handler) override {
        diagnosticsAPI->setOnUnsubscribeHandler(handler);
    }

    void setOnTerminateNetworkHandler(
        std::function<api::TerminateNetworkHandler> handler) override {
        diagnosticsAPI->setOnTerminateNetworkHandler(handler);
    }
};
}  // namespace directmq
