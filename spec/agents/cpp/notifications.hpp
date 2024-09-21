#pragma once

#include <cstdint>
#include <nlohmann/json.hpp>
#include <string>

using json = nlohmann::json;

/**
 * Agent API
 */

struct ConnectionEstablishedNotification {
    std::string bridgedNodeId;
};

struct ConnectionLostNotification {
    std::string bridgedNodeId;
    std::string reason;
};

struct ReadyNotification {
    std::string time;
};

struct FatalNotification {
    std::string err;
};

struct MessageReceivedNotification {
    std::string topic;
    std::string payload;
};

struct SubscribedNotification {
    int32_t subscriptionId;
};

struct StoppedNotification {
    std::string reason;
};

/**
 * Diagnostics API
 */

struct OnPublicationNotification {
    int32_t ttl;
    std::vector<std::string> traversed;
    std::string topic;
    uint8_t deliveryStrategy;
    std::string payload;
};

struct OnSubscriptionNotification {
    int32_t ttl;
    std::vector<std::string> traversed;
    std::string topic;
};

struct OnUnsubscribeNotification {
    int32_t ttl;
    std::vector<std::string> traversed;
    std::string topic;
};

struct OnNetworkTerminationNotification {
    int32_t ttl;
    std::vector<std::string> traversed;
    std::string reason;
};

/**
 * Universal notification
 */

struct UniversalNotification {
    // Agent API
    ReadyNotification *ready = nullptr;
    FatalNotification *fatal = nullptr;
    MessageReceivedNotification *messageReceived = nullptr;
    SubscribedNotification *subscribed = nullptr;
    StoppedNotification *stopped = nullptr;

    // Diagnostics API
    ConnectionEstablishedNotification *connectionEstablished = nullptr;
    ConnectionLostNotification *connectionLost = nullptr;
    OnPublicationNotification *onPublication = nullptr;
    OnSubscriptionNotification *onSubscription = nullptr;
    OnUnsubscribeNotification *onUnsubscribe = nullptr;
    OnNetworkTerminationNotification *onNetworkTermination = nullptr;

    std::string toJson() const {
        json j;
        if (ready) {
            j["ready"] = {{"time", ready->time}};
        }

        if (fatal) {
            j["fatal"] = {{"err", fatal->err}};
        }

        if (messageReceived) {
            j["messageReceived"] = {{"topic", messageReceived->topic},
                                    {"payload", messageReceived->payload}};
        }

        if (subscribed) {
            j["subscribed"] = {{"subscriptionId", subscribed->subscriptionId}};
        }

        if (stopped) {
            j["stopped"] = {{"reason", stopped->reason}};
        }

        if (connectionEstablished) {
            j["connectionEstablished"] = {
                {"bridgedNodeId", connectionEstablished->bridgedNodeId}};
        }

        if (connectionLost) {
            j["connectionLost"] = {
                {"bridgedNodeId", connectionLost->bridgedNodeId},
                {"reason", connectionLost->reason}};
        }

        if (onPublication) {
            j["onPublication"] = {
                {"ttl", onPublication->ttl},
                {"traversed", onPublication->traversed},
                {"topic", onPublication->topic},
                {"deliveryStrategy", onPublication->deliveryStrategy},
                {"payload", onPublication->payload}};
        }

        if (onSubscription) {
            j["onSubscription"] = {{"ttl", onSubscription->ttl},
                                   {"traversed", onSubscription->traversed},
                                   {"topic", onSubscription->topic}};
        }

        if (onUnsubscribe) {
            j["onUnsubscribe"] = {{"ttl", onUnsubscribe->ttl},
                                  {"traversed", onUnsubscribe->traversed},
                                  {"topic", onUnsubscribe->topic}};
        }

        if (onNetworkTermination) {
            j["onNetworkTermination"] = {
                {"ttl", onNetworkTermination->ttl},
                {"traversed", onNetworkTermination->traversed},
                {"reason", onNetworkTermination->reason}};
        }

        return j.dump();
    }

    static UniversalNotification makeReady(const std::string &time) {
        UniversalNotification n;
        n.ready = new ReadyNotification{.time = time};
        return n;
    }

    static UniversalNotification makeStopped(const std::string &reason) {
        UniversalNotification n;
        n.stopped = new StoppedNotification{.reason = reason};
        return n;
    }

    static UniversalNotification makeFatal(const std::string &err) {
        UniversalNotification n;
        n.fatal = new FatalNotification{.err = err};
        return n;
    }

    static UniversalNotification makeMessageReceived(
        const std::string &topic, const std::string &payload) {
        UniversalNotification n;
        n.messageReceived =
            new MessageReceivedNotification{.topic = topic, .payload = payload};
        return n;
    }

    static UniversalNotification makeSubscribed(int32_t subscriptionId) {
        UniversalNotification n;
        n.subscribed =
            new SubscribedNotification{.subscriptionId = subscriptionId};
        return n;
    }

    static UniversalNotification makeConnectionEstablished(
        const std::string &bridgedNodeId) {
        UniversalNotification n;
        n.connectionEstablished = new ConnectionEstablishedNotification{
            .bridgedNodeId = bridgedNodeId};
        return n;
    }

    static UniversalNotification makeConnectionLost(
        const std::string &bridgedNodeId, const std::string &reason) {
        UniversalNotification n;
        n.connectionLost = new ConnectionLostNotification{
            .bridgedNodeId = bridgedNodeId, .reason = reason};
        return n;
    }

    static UniversalNotification makeOnPublication(
        int32_t ttl, const std::vector<std::string> &traversed,
        const std::string &topic, uint8_t deliveryStrategy,
        const std::string &payload) {
        UniversalNotification n;
        n.onPublication =
            new OnPublicationNotification{.ttl = ttl,
                                          .traversed = traversed,
                                          .topic = topic,
                                          .deliveryStrategy = deliveryStrategy,
                                          .payload = payload};
        return n;
    }

    static UniversalNotification makeOnSubscription(
        int32_t ttl, const std::vector<std::string> &traversed,
        const std::string &topic) {
        UniversalNotification n;
        n.onSubscription = new OnSubscriptionNotification{
            .ttl = ttl, .traversed = traversed, .topic = topic};
        return n;
    }

    static UniversalNotification makeOnUnsubscribe(
        int32_t ttl, const std::vector<std::string> &traversed,
        const std::string &topic) {
        UniversalNotification n;
        n.onUnsubscribe = new OnUnsubscribeNotification{
            .ttl = ttl, .traversed = traversed, .topic = topic};
        return n;
    }

    static UniversalNotification makeOnNetworkTermination(
        int32_t ttl, const std::vector<std::string> &traversed,
        const std::string &reason) {
        UniversalNotification n;
        n.onNetworkTermination = new OnNetworkTerminationNotification{
            .ttl = ttl, .traversed = traversed, .reason = reason};
        return n;
    }
};
