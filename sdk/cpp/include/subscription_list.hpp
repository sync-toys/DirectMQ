#pragma once

#include <list>
#include <memory>
#include <random>
#include <stdexcept>
#include <string>

#include "topic_pattern.hpp"

namespace directmq::subscriptions {
template <typename THandler = void>
struct Subscription {
    int id;
    std::shared_ptr<const std::string> topicPattern;
    THandler handler;
};

template <typename THandler = void*>
class SubscriptionList {
   private:
    std::list<Subscription<THandler>> subscriptions;

    int getRandomId() {
        ushort tries = 0;

        while (tries < 1000) {
            int id = std::rand();

            bool found = false;
            for (auto& subscription : subscriptions) {
                if (subscription.id == id) {
                    found = true;
                    break;
                }
            }

            if (!found) {
                return id;
            }

            tries++;
        }

        throw std::runtime_error("Failed to generate a unique subscription ID");
    }

   public:
    int addSubscription(std::shared_ptr<const std::string> topicPattern,
                        THandler handler) {
        Subscription<THandler> subscription;

        subscription.id = getRandomId();
        subscription.topicPattern = topicPattern;
        subscription.handler = handler;

        subscriptions.push_back(subscription);

        return subscription.id;
    }

    void removeSubscription(int id) {
        for (auto it = subscriptions.begin(); it != subscriptions.end(); it++) {
            if (it->id == id) {
                subscriptions.erase(it);
                return;
            }
        }
    }

    const std::list<Subscription<THandler>> getTriggeredSubscriptions(
        const std::string& topic) const {
        std::list<Subscription<THandler>> triggeredSubscriptions;

        for (auto& subscription : subscriptions) {
            if (topics::matchTopicPattern(*subscription.topicPattern, topic)) {
                triggeredSubscriptions.push_back(subscription);
            }
        }

        return triggeredSubscriptions;
    }

    bool willHandleTopic(const std::string& topic) const {
        for (auto& subscription : subscriptions) {
            if (topics::matchTopicPattern(*subscription.topicPattern, topic)) {
                return true;
            }
        }

        return false;
    }

    bool alreadyHandlesPattern(const std::string& pattern) const {
        for (auto& subscription : subscriptions) {
            if (topics::isSubtopicPattern(*subscription.topicPattern,
                                          pattern)) {
                return true;
            }
        }

        return false;
    }

    const std::list<Subscription<THandler>> getSubscriptions() const {
        return subscriptions;
    }

    const std::list<std::shared_ptr<const std::string>>
    getUniqueSubscribedTopics() const {
        std::list<std::shared_ptr<const std::string>> uniqueSubscribedTopics;

        for (auto& subscription : subscriptions) {
            bool found = false;
            for (auto& uniqueSubscribedTopic : uniqueSubscribedTopics) {
                if (*subscription.topicPattern == *uniqueSubscribedTopic) {
                    found = true;
                    break;
                }
            }

            if (!found) {
                uniqueSubscribedTopics.push_back(subscription.topicPattern);
            }
        }

        return uniqueSubscribedTopics;
    }

    const std::list<std::shared_ptr<const std::string>>
    getOnlyTopLevelSubscribedTopics() const {
        auto uniqueSubscribedTopics = getUniqueSubscribedTopics();
        auto topLevelTopics =
            topics::deduplicateOverlappingTopics(uniqueSubscribedTopics);

        return topLevelTopics;
    }

    void clear() { subscriptions.clear(); }
};
}  // namespace directmq::subscriptions
