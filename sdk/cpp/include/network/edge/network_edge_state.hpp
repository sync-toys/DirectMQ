#pragma once

#include "../../protocol/decoder.hpp"
#include "../participant.hpp"
#include "state_name.hpp"

namespace directmq::network::edge {
class NetworkEdgeStateManager;

class NetworkEdgeState : public NetworkParticipant,
                         public protocol::DecodingHandler {
   protected:
    NetworkEdgeStateManager* edge;

   public:
    NetworkEdgeState(NetworkEdgeStateManager* edge) : edge(edge) {}
    virtual ~NetworkEdgeState() = default;

    virtual EdgeStateName getStateName() const = 0;
    virtual void onSet() = 0;
};
}  // namespace directmq::network::edge
