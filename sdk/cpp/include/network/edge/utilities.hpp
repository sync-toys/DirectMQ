#pragma once

#include <list>
#include <memory>
#include <string>

#include "../../protocol/messages/data_frame.hpp"

namespace directmq::network::edge::internal {
bool checkForNetworkLoops(const protocol::messages::DataFrame& frame) {
    std::list<const std::string*> existingHosts;

    for (auto& host : frame.traversed) {
        for (auto existingHost : existingHosts) {
            if (&host == existingHost) {
                return true;
            }
        }

        existingHosts.push_back(&host);
    }

    return false;
}

protocol::messages::DataFrame updateFrameTraversedAndTTL(
    const protocol::messages::DataFrame& frame, const std::string& hostID) {
    if (frame.traversed.size() > 0 && frame.traversed.back() == hostID) {
        return frame;
    }

    protocol::messages::DataFrame updatedFrame = frame;
    updatedFrame.traversed.push_back(hostID);
    updatedFrame.ttl--;

    return updatedFrame;
}

std::string buildLoopDetectedMessage(
    const protocol::messages::DataFrame& frame) {
    std::string message = "Network loop detected: ";

    for (auto& host : frame.traversed) {
        if (host == frame.traversed.back()) {
            break;
        }

        message += host + " -> ";
    }

    message += frame.traversed.back();

    return message;
}
}  // namespace directmq::network::edge::internal
