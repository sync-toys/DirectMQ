#pragma once

#include <cstdint>
#include <nlohmann/json.hpp>
#include <string>

using json = nlohmann::json;

/**
 * Connection API
 */

struct SetupCommand {
    int32_t ttl;
    std::string nodeId;
    uint64_t maxMessageSize;
};

struct ListenCommand {
    std::string address;
    uint16_t port;
};

struct ConnectCommand {
    std::string address;
    uint16_t port;
};

struct StopCommand {
    std::string reason;
};

/**
 * Native API
 */

struct PublishCommand {
    std::string topic;
    directmq::protocol::messages::DeliveryStrategy deliveryStrategy;
    std::string payload;
};

struct SubscribeTopicCommand {
    std::string topic;
};

struct UnsubscribeTopicCommand {
    int32_t subscriptionId;
};

/**
 * Universal command
 */

struct UniversalCommand {
    // Connection API
    SetupCommand *setup = nullptr;
    ListenCommand *listen = nullptr;
    ConnectCommand *connect = nullptr;
    StopCommand *stop = nullptr;

    // Native API
    PublishCommand *publish = nullptr;
    SubscribeTopicCommand *subscribeTopic = nullptr;
    UnsubscribeTopicCommand *unsubscribeTopic = nullptr;

    static UniversalCommand fromJson(std::string rawCommand) {
        json c = json::parse(rawCommand);
        UniversalCommand r;

        if (c.contains("setup")) {
            r.setup = new SetupCommand();
            r.setup->ttl = c["setup"]["ttl"];
            r.setup->nodeId = c["setup"]["nodeId"];
            r.setup->maxMessageSize = c["setup"]["maxMessageSize"];
        }

        if (c.contains("listen")) {
            r.listen = new ListenCommand();
            r.listen->address = c["listen"]["address"];
            r.listen->port = c["listen"]["port"];
        }

        if (c.contains("connect")) {
            r.connect = new ConnectCommand();
            r.connect->address = c["connect"]["address"];
            r.connect->port = c["connect"]["port"];
        }

        if (c.contains("stop")) {
            r.stop = new StopCommand();
            r.stop->reason = c["stop"]["reason"];
        }

        if (c.contains("publish")) {
            r.publish = new PublishCommand();
            r.publish->topic = c["publish"]["topic"];
            r.publish->deliveryStrategy =
                (directmq::protocol::messages::DeliveryStrategy)(
                    uint8_t)c["publish"]["deliveryStrategy"];
            r.publish->payload = c["publish"]["payload"];
        }

        if (c.contains("subscribeTopic")) {
            r.subscribeTopic = new SubscribeTopicCommand();
            r.subscribeTopic->topic = c["subscribeTopic"]["topic"];
        }

        if (c.contains("unsubscribeTopic")) {
            r.unsubscribeTopic = new UnsubscribeTopicCommand();
            r.unsubscribeTopic->subscriptionId =
                c["unsubscribeTopic"]["subscriptionId"];
        }

        return r;
    }

    ~UniversalCommand() {
        delete setup;
        delete listen;
        delete connect;
        delete stop;
        delete publish;
        delete subscribeTopic;
        delete unsubscribeTopic;
    }
};
