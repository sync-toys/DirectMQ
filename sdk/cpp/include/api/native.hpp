#pragma once

#include <cstdint>
#include <stdexcept>
#include <string>
#include <vector>

#include "../network/global.hpp"
#include "../network/participant.hpp"
#include "../protocol/messages/publish.hpp"
#include "../subscription_list.hpp"

namespace directmq::network {
// just pre-declare, required to
// make NativeAPIAdapter a friend
class NetworkNode;
}  // namespace directmq::network

namespace directmq::api {
class NativePublicationAPI {
   public:
    virtual ~NativePublicationAPI() = default;

    virtual void publish(
        const std::string& topic, const std::vector<uint8_t>& message,
        const protocol::messages::DeliveryStrategy& deliveryStrategy) = 0;
};

class NativeAPIAdapter : public NativePublicationAPI,
                         public network::NetworkParticipant {
   protected:
    std::shared_ptr<network::GlobalNetwork> globalNetwork;

   private:
    // this method is called by the NetworkNode
    // to setup the globalNetwork pointer
    void setup(std::shared_ptr<network::GlobalNetwork> globalNetwork) {
        if (globalNetwork == nullptr) {
            throw std::runtime_error("globalNetwork cannot be null");
        }

        this->globalNetwork = globalNetwork;
    }

    friend class network::NetworkNode;
};

using LambdaMessageHandler = void(const std::string& topic,
                                  const std::vector<uint8_t>& message);

class NativeLambdaAPI : public NativePublicationAPI {
   public:
    virtual int subscribe(const std::string& topic,
                          std::function<LambdaMessageHandler> handler) = 0;

    virtual void unsubscribe(int subscriptionId) = 0;
};

class NativeLambdaAPIAdapter : public NativeLambdaAPI, public NativeAPIAdapter {
   protected:
    subscriptions::SubscriptionList<std::function<LambdaMessageHandler>>
        subscriptions;

    protocol::messages::DataFrame getInitialDataFrame() const {
        return protocol::messages::DataFrame{
            .ttl = globalNetwork->config.hostTTL,
            .traversed = {globalNetwork->config.hostID},
        };
    }

    void updateSubscriptions(
        const std::list<std::shared_ptr<const std::string>>& oldTopics) {
        auto updatedTopLevelSubscriptions =
            subscriptions.getOnlyTopLevelSubscribedTopics();
        auto diff = topics::getDeduplicatedOverlappingTopicsDiff(
            oldTopics, updatedTopLevelSubscriptions);

        for (auto& topic : diff.added) {
            protocol::messages::SubscribeMessage subscription{
                .frame = getInitialDataFrame(),
                .topic = *topic,
            };

            globalNetwork->subscribed(subscription);
        }

        for (auto& topic : diff.removed) {
            protocol::messages::UnsubscribeMessage unsubscription{
                .frame = getInitialDataFrame(),
                .topic = *topic,
            };

            globalNetwork->unsubscribed(unsubscription);
        }
    }

   public:
    void publish(
        const std::string& topic, const std::vector<uint8_t>& payload,
        const protocol::messages::DeliveryStrategy& deliveryStrategy) override {
        if (topics::isCorrectTopicPattern(topic)) {
            throw std::runtime_error("incorrect topic");
        }

        if (payload.size() == 0) {
            throw std::runtime_error("empty payload is forbidden");
        }

        protocol::messages::PublishMessage publication{
            .frame = getInitialDataFrame(),
            .topic = topic,
            .deliveryStrategy = deliveryStrategy,
            .payload = payload};

        globalNetwork->published(publication);
    }

    int subscribe(const std::string& topic,
                  std::function<LambdaMessageHandler> handler) {
        if (topics::isCorrectTopicPattern(topic)) {
            throw std::runtime_error("incorrect topic");
        }

        if (handler == nullptr) {
            throw std::runtime_error("handler cannot be null");
        }

        auto oldTopics = subscriptions.getOnlyTopLevelSubscribedTopics();
        int subscriptionId = subscriptions.addSubscription(
            std::make_shared<const std::string>(topic), handler);
        updateSubscriptions(oldTopics);
        return subscriptionId;
    }

    void unsubscribe(int subscriptionId) {
        auto oldTopics = subscriptions.getOnlyTopLevelSubscribedTopics();
        subscriptions.removeSubscription(subscriptionId);
        updateSubscriptions(oldTopics);
    }

    std::list<std::shared_ptr<const std::string>> getSubscribedTopics()
        const override {
        return subscriptions.getOnlyTopLevelSubscribedTopics();
    }

    bool willHandleTopic(const std::string& topic) const override {
        return subscriptions.willHandleTopic(topic);
    }

    bool alreadyHandlesPattern(const std::string& pattern) const override {
        return subscriptions.alreadyHandlesPattern(pattern);
    }

    bool isOriginOfFrame(
        const protocol::messages::DataFrame& frame) const override {
        return frame.traversed.size() == 0;
    }

    bool handlePublish(
        const protocol::messages::PublishMessage& publication) override {
        if (isOriginOfFrame(publication.frame)) {
            return false;
        }

        auto subscribers = utils::randomOrder(
            subscriptions.getTriggeredSubscriptions(publication.topic));
        if (subscribers.size() == 0) {
            return false;
        }

        if (publication.deliveryStrategy ==
            protocol::messages::DeliveryStrategy::AT_MOST_ONCE) {
            subscribers.front().handler(publication.topic, publication.payload);
            return true;
        }

        for (auto& subscriber : subscribers) {
            subscriber.handler(publication.topic, publication.payload);
        }

        return true;
    }

    void handleSubscribe(
        const protocol::messages::SubscribeMessage& subscription) override {
        (void)subscription;
        /* No-op */
    }

    void handleUnsubscribe(
        const protocol::messages::UnsubscribeMessage& unsubscription) override {
        (void)unsubscription;
        /* No-op */
    }

    void handleTerminateNetwork(
        const protocol::messages::TerminateNetworkMessage& termination)
        override {
        (void)termination;
        /* No-op */
    }
};
}  // namespace directmq::api
