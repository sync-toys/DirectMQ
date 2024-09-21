#pragma once

#include <memory>

#include "../api/diagnostics.hpp"
#include "../protocol/messages.hpp"
#include "../topic_pattern.hpp"
#include "../utils/only_unique_strings.hpp"
#include "../utils/random_order.hpp"
#include "config.hpp"
#include "participant.hpp"

namespace directmq::network {
class GlobalNetwork {
   private:
    std::list<std::shared_ptr<NetworkParticipant>> participants;

    std::shared_ptr<NetworkParticipant> nativeAPI;
    std::shared_ptr<api::DiagnosticsAPIAdapter> diagnostics;

    std::list<std::shared_ptr<const std::string>>
    getSubscribedTopicsExcludingOriginOfMessage(
        const protocol::messages::DataFrame& frame) const {
        std::list<std::shared_ptr<const std::string>> result;
        for (auto& participant : participants) {
            if (participant->isOriginOfFrame(frame)) {
                continue;
            }

            for (auto& topic : participant->getSubscribedTopics()) {
                result.push_back(topic);
            }
        }
        return result;
    }

   public:
    const NetworkNodeConfig config;

    GlobalNetwork(const NetworkNodeConfig& config,
                  std::shared_ptr<NetworkParticipant> nativeAPI,
                  std::shared_ptr<api::DiagnosticsAPIAdapter> diagnostics)
        : config(config),
          nativeAPI(nativeAPI),
          diagnostics(diagnostics),
          participants({nativeAPI}) {}

    std::shared_ptr<api::DiagnosticsAPIAdapter> getDiagnostics() {
        return diagnostics;
    }

    void published(const protocol::messages::PublishMessage& publication) {
        diagnostics->handlePublish(publication);

        for (auto& participant : utils::randomOrder(participants)) {
            bool handled = participant->handlePublish(publication);
            if (publication.deliveryStrategy ==
                    protocol::messages::DeliveryStrategy::AT_MOST_ONCE &&
                handled) {
                return;
            }
        }
    }

    void subscribed(const protocol::messages::SubscribeMessage& subscription) {
        diagnostics->handleSubscribe(subscription);

        auto topLevelSubscriptions =
            getSubscribedTopicsExcludingOriginOfMessage(subscription.frame);

        for (auto& participant : participants) {
            participant->handleSubscribe(subscription);
        }

        if (subscription.frame.traversed.size() == 0) {
            // we are skipping the unsubscribe synchronization part
            // because current node is origin of the message
            // and synchronization will be done by native api itself
            return;
        }

        auto updatedTopLevelSubscriptions = getAllSubscribedTopics();
        auto diff = topics::getDeduplicatedOverlappingTopicsDiff(
            topLevelSubscriptions, updatedTopLevelSubscriptions);

        for (auto& topic : diff.removed) {
            protocol::messages::UnsubscribeMessage unsubscription;
            unsubscription.frame = subscription.frame;
            unsubscription.topic = *topic;

            for (auto& participant : participants) {
                if (participant->alreadyHandlesPattern(*topic)) {
                    continue;
                }

                participant->handleUnsubscribe(unsubscription);
            }
        }
    }

    void unsubscribed(
        const protocol::messages::UnsubscribeMessage& unsubscription) {
        diagnostics->handleUnsubscribe(unsubscription);

        for (auto& participant : participants) {
            if (participant->alreadyHandlesPattern(unsubscription.topic)) {
                return;
            }
        }

        // if every other participants are already not handling the topic
        // we can safely send unsubscribe message to other network participants

        for (auto& participant : participants) {
            participant->handleUnsubscribe(unsubscription);
        }
    }

    void terminated(
        const protocol::messages::TerminateNetworkMessage& termination) {
        diagnostics->handleTerminateNetwork(termination);

        for (auto& participant : participants) {
            participant->handleTerminateNetwork(termination);
        }
    }

    std::list<std::shared_ptr<const std::string>> getAllSubscribedTopics()
        const {
        std::list<std::shared_ptr<const std::string>> result;

        for (auto& participant : participants) {
            for (auto& topic : participant->getSubscribedTopics()) {
                result.push_back(topic);
            }
        }

        return utils::onlyUniqueStrings(result);
    }

    void addParticipant(std::shared_ptr<NetworkParticipant> participant) {
        participants.push_back(participant);
    }

    void removeParticipant(std::shared_ptr<NetworkParticipant> participant) {
        participants.remove(participant);
    }
};
}  // namespace directmq::network
