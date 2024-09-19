#pragma once

#include <functional>
#include <list>

#include "../api/diagnostics.hpp"
#include "../api/native.hpp"
#include "../portal.hpp"
#include "edge/network_edge.hpp"
#include "global.hpp"

namespace directmq::network {
class EdgeManager {
   public:
    virtual ~EdgeManager() = default;

    virtual std::shared_ptr<edge::NetworkEdge> addListeningEdge(
        std::shared_ptr<portal::Portal> portal) = 0;
    virtual std::shared_ptr<edge::NetworkEdge> addConnectingEdge(
        std::shared_ptr<portal::Portal> portal) = 0;

    virtual bool removeEdge(std::shared_ptr<portal::Portal> portal,
                            const std::string& reason) = 0;
};

using ConnectionLostHandler = void(const std::string& bridgedNodeID,
                                   const std::string& reason,
                                   portal::Portal& portal);

class NetworkNode : public EdgeManager, public api::DiagnosticsAPI {
   private:
    std::shared_ptr<network::GlobalNetwork> globalNetwork;
    std::list<std::shared_ptr<edge::NetworkEdge>> edges;

    std::shared_ptr<api::NativeAPIAdapter> nativeAPI;
    std::shared_ptr<api::DiagnosticsAPI> diagnosticsAPI;

    std::shared_ptr<protocol::Encoder> encoder;
    std::shared_ptr<protocol::Decoder> decoder;

    std::function<ConnectionLostHandler> onConnectionLostCallback;
    std::function<void()> onNodeFullyClosedCallback;

    void handleEdgeConnectionLost(const std::string& bridgedNodeID,
                                  const std::string& reason,
                                  portal::Portal& portal) {
        auto edge = std::find_if(edges.begin(), edges.end(),
                                 [&](std::shared_ptr<edge::NetworkEdge> edge) {
                                     return edge->getPortal().get() == &portal;
                                 });

        if (edge != edges.end()) {
            globalNetwork->removeParticipant(*edge);
            edges.erase(edge);
        }

        if (onConnectionLostCallback) {
            onConnectionLostCallback(bridgedNodeID, reason, portal);
        }

        if (edges.size() == 0 && onNodeFullyClosedCallback) {
            onNodeFullyClosedCallback();
            onNodeFullyClosedCallback = nullptr;
        }
    }

   public:
    NetworkNode(const network::NetworkNodeConfig& config,
                std::shared_ptr<api::NativeAPIAdapter> nativeAPI,
                std::shared_ptr<api::DiagnosticsAPI> diagnosticsAPI,
                std::shared_ptr<protocol::Encoder> encoder,
                std::shared_ptr<protocol::Decoder> decoder)
        : nativeAPI(nativeAPI),
          diagnosticsAPI(diagnosticsAPI),
          encoder(encoder),
          decoder(decoder) {
        globalNetwork = std::make_shared<network::GlobalNetwork>(
            config, nativeAPI, diagnosticsAPI);

        nativeAPI->setup(globalNetwork);
        diagnosticsAPI->setOnConnectionLostHandler(
            [&](const std::string& bridgedNodeID, const std::string& reason,
                portal::Portal& portal) {
                handleEdgeConnectionLost(bridgedNodeID, reason, portal);
            });
    }

    std::shared_ptr<edge::NetworkEdge> addListeningEdge(
        std::shared_ptr<portal::Portal> portal) override {
        auto edge = std::make_shared<edge::NetworkEdge>(globalNetwork, portal,
                                                        decoder, encoder);

        edges.push_back(edge);
        globalNetwork->addParticipant(edge);

        edge->setConnectingState(false);

        return edge;
    }

    std::shared_ptr<edge::NetworkEdge> addConnectingEdge(
        std::shared_ptr<portal::Portal> portal) override {
        auto edge = std::make_shared<edge::NetworkEdge>(globalNetwork, portal,
                                                        decoder, encoder);

        edges.push_back(edge);
        globalNetwork->addParticipant(edge);

        edge->setConnectingState(true);

        return edge;
    }

    bool removeEdge(std::shared_ptr<portal::Portal> portal,
                    const std::string& reason) override {
        auto edge = std::find_if(edges.begin(), edges.end(), [&](auto& edge) {
            return edge->getPortal() == portal;
        });

        if (edge == edges.end()) {
            return false;
        }

        if ((*edge)->getStateName() == edge::EdgeStateName::DISCONNECTING ||
            (*edge)->getStateName() == edge::EdgeStateName::DISCONNECTED) {
            return true;
        }

        (*edge)->setDisconnectingState(reason);

        return true;
    }

    void closeNode(const std::string& reason,
                   std::function<void()> onNodeFullyClosed) {
        onNodeFullyClosedCallback = onNodeFullyClosed;

        for (auto& edge : edges) {
            edge->setDisconnectingState(reason);
        }
    }
};
}  // namespace directmq::network
