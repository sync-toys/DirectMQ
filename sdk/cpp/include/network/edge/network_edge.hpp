#pragma once

#include "network_edge_state_connected.hpp"
#include "network_edge_state_connecting.hpp"
#include "network_edge_state_disconnected.hpp"
#include "network_edge_state_disconnecting.hpp"
#include "network_edge_state_manager.hpp"

namespace directmq::network::edge {
class NetworkEdge : public NetworkEdgeStateManager {
   private:
    void setState(NetworkEdgeState* newState) {
        NetworkEdgeState* oldState = state;

        state = newState;
        state->onSet();

        if (oldState != nullptr) {
            delete state;
        }
    }

   public:
    NetworkEdge(std::shared_ptr<GlobalNetwork> globalNetwork,
                std::shared_ptr<portal::Portal> portal,
                std::shared_ptr<protocol::Decoder> decoder,
                std::shared_ptr<protocol::Encoder> encoder)
        : NetworkEdgeStateManager(globalNetwork, portal, decoder, encoder) {}

    void setDisconnectedState(const std::string& reason) override {
        setState(new NetworkEdgeStateDisconnected(this, reason));
    }

    void setConnectingState(bool initializeConnection) override {
        setState(new NetworkEdgeStateConnecting(this, initializeConnection));
    }

    void setConnectedState() override {
        setState(new NetworkEdgeStateConnected(this));
    }

    void setDisconnectingState(const std::string& reason) override {
        setState(new NetworkEdgeStateDisconnecting(this, reason));
    }
};
}  // namespace directmq::network::edge
