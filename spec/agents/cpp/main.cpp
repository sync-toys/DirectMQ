#include <atomic>
#include <chrono>
#include <directmq.hpp>
#include <iomanip>
#include <iostream>
#include <mutex>
#include <sstream>
#include <string>
#include <thread>

#include "commands.hpp"
#include "notifications.hpp"

std::string getCurrentTimeString() {
    auto now = std::chrono::system_clock::now();
    std::time_t now_time = std::chrono::system_clock::to_time_t(now);
    std::tm now_tm = *std::localtime(&now_time);

    std::ostringstream oss;
    oss << std::put_time(&now_tm, "%Y-%m-%d %H:%M:%S");

    return oss.str();
}

std::mutex nodeMutex;

std::mutex exitMutex;
const int NO_EXIT = -1;
std::atomic<int> exitFlag(NO_EXIT);

directmq::DirectMQNode *node;

void log(const std::string &message) {
    std::cout << message << std::endl;
    std::cout.flush();
}

void sendNotification(const UniversalNotification &notification) {
    std::cout << notification.toJson() << std::endl;
    std::cout.flush();
}

void exitAgent(int exitCode) {
    exitFlag = exitCode;
    exitMutex.unlock();
}

void fatal(const std::string &error) {
    sendNotification(UniversalNotification::makeFatal(error));
    exitAgent(1);
}

void registerDiagnosticsHandlers() {
    node->setOnConnectionEstablishedHandler(
        [](const std::string &bridgedNodeId, directmq::portal::Portal &portal) {
            sendNotification(UniversalNotification::makeConnectionEstablished(
                bridgedNodeId));
        });

    node->setOnEdgeDisconnectionHandler(
        [](const std::string &bridgedNodeId, const std::string &reason,
           directmq::network::edge::NetworkEdge &edge,
           directmq::portal::Portal &portal) {
            sendNotification(UniversalNotification::makeConnectionLost(
                bridgedNodeId, reason));
        });

    node->setOnPublicationHandler(
        [](directmq::protocol::messages::PublishMessage publication) {
            sendNotification(UniversalNotification::makeOnPublication(
                publication.frame.ttl,
                std::vector<std::string>(publication.frame.traversed.begin(),
                                         publication.frame.traversed.end()),
                publication.topic, publication.deliveryStrategy,
                std::string(publication.payload.begin(),
                            publication.payload.end())));
        });

    node->setOnSubscriptionHandler(
        [](const directmq::protocol::messages::SubscribeMessage &subscription) {
            sendNotification(UniversalNotification::makeOnSubscription(
                subscription.frame.ttl,
                std::vector<std::string>(subscription.frame.traversed.begin(),
                                         subscription.frame.traversed.end()),
                subscription.topic));
        });

    node->setOnUnsubscribeHandler(
        [](const directmq::protocol::messages::UnsubscribeMessage
               &unsubscription) {
            sendNotification(UniversalNotification::makeOnUnsubscribe(
                unsubscription.frame.ttl,
                std::vector<std::string>(unsubscription.frame.traversed.begin(),
                                         unsubscription.frame.traversed.end()),
                unsubscription.topic));
        });

    node->setOnTerminateNetworkHandler(
        [](const directmq::protocol::messages::TerminateNetworkMessage
               &termination) {
            sendNotification(UniversalNotification::makeOnNetworkTermination(
                termination.frame.ttl,
                std::vector<std::string>(termination.frame.traversed.begin(),
                                         termination.frame.traversed.end()),
                termination.reason));
        });
}

void handleSetupCommand(SetupCommand command) {
    log("Setting up DirectMQ node");

    directmq::network::NetworkNodeConfig config{
        .hostTTL = command.ttl,
        .hostMaxIncomingMessageSize = command.maxMessageSize,
        .hostID = command.nodeId};

    node = new directmq::DirectMQNode(config);

    registerDiagnosticsHandlers();

    log("Setup complete");
}

void handleListenCommand(ListenCommand command) {
    log("Listening as server at " + command.address + ":" +
        std::to_string(command.port));
}

void handleConnectCommand(ConnectCommand command) {
    log("Connecting as client to " + command.address + ":" +
        std::to_string(command.port));
}

void handleStopCommand(StopCommand command) {
    log("Stopping DirectMQ node: " + command.reason);

    node->closeNode(command.reason, [command]() {
        sendNotification(UniversalNotification::makeStopped(command.reason));
    });

    log("Clean exit 0");
    exitAgent(0);
}

void handlePublishCommand(PublishCommand command) {
    log("Publishing message to topic " + command.topic);

    std::vector<uint8_t> payload(command.payload.begin(),
                                 command.payload.end());
    node->publish(command.topic, payload, command.deliveryStrategy);
}

void handleSubscribeCommand(SubscribeTopicCommand command) {
    log("Subscribing to topic " + command.topic);

    auto subscriptionId = node->subscribe(
        command.topic, [command](const std::string &topic,
                                 const std::vector<uint8_t> &payload) {
            sendNotification(UniversalNotification::makeMessageReceived(
                topic, std::string(payload.begin(), payload.end())));
        });

    log("Subscription ID: " + std::to_string(subscriptionId));

    sendNotification(UniversalNotification::makeSubscribed(subscriptionId));
}

void handleUnsubscribeCommand(UnsubscribeTopicCommand command) {
    log("Unsubscribing from subscription ID " +
        std::to_string(command.subscriptionId));

    node->unsubscribe(command.subscriptionId);
}

void handleIncomingCommand(const UniversalCommand &command) {
    std::lock_guard<std::mutex> lock(nodeMutex);

    if (command.setup) {
        handleSetupCommand(*command.setup);
    }

    if (command.listen) {
        handleListenCommand(*command.listen);
    }

    if (command.connect) {
        handleConnectCommand(*command.connect);
    }

    if (command.stop) {
        handleStopCommand(*command.stop);
    }

    if (command.publish) {
        handlePublishCommand(*command.publish);
    }

    if (command.subscribeTopic) {
        handleSubscribeCommand(*command.subscribeTopic);
    }

    if (command.unsubscribeTopic) {
        handleUnsubscribeCommand(*command.unsubscribeTopic);
    }
}

std::string readCommandFromStdin() {
    std::string rawCommand;
    std::getline(std::cin, rawCommand);
    return rawCommand;
}

void runCommandLoop() {
    while (exitFlag == NO_EXIT) {
        try {
            std::string rawCommand = readCommandFromStdin();
            if (rawCommand.empty()) {
                return;
            }

            UniversalCommand command = UniversalCommand::fromJson(rawCommand);
            handleIncomingCommand(command);
        } catch (const std::exception &e) {
            fatal("Command loop fatal failure: " + std::string(e.what()));
        }
    }
}

void runNodeLoop() {
    while (exitFlag == NO_EXIT) {
        try {
            std::lock_guard<std::mutex> lock(nodeMutex);
        } catch (const std::exception &e) {
            fatal("Node loop fatal failure: " + std::string(e.what()));
        }
    }
}

int main() {
    log("Starting DirectMQ testing agent");

    // initialize exit mutex
    exitMutex.lock();

    log("Starting node loop");
    std::thread nodeLoop(runNodeLoop);

    log("Starting command loop");
    std::thread commandLoop(runCommandLoop);

    log("Agent ready");
    sendNotification(UniversalNotification::makeReady(getCurrentTimeString()));

    // exit mutex is unlocked when exitAgent is called
    exitMutex.lock();

    log("Exiting agent with code " + std::to_string(exitFlag));
    return exitFlag;
}
