#pragma once

#include <list>
#include <memory>
#include <string>

#include "../protocol/messages.hpp"

namespace directmq::network {
class NetworkParticipant {
   public:
    virtual ~NetworkParticipant() = default;

    virtual std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const = 0;
    virtual bool willHandleTopic(const std::string& topic) const = 0;
    virtual bool alreadyHandlesPattern(const std::string& pattern) const = 0;
    virtual bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const = 0;

    virtual bool handlePublish(
        const protocol::messages::PublishMessage& publication) = 0;
    virtual void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) = 0;
    virtual void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) = 0;
    virtual void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination) = 0;
};
}  // namespace directmq::network
