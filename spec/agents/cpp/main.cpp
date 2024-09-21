#include <chrono>
#include <directmq.hpp>
#include <iomanip>
#include <iostream>
#include <portals/websocket/client.hpp>
#include <portals/websocket/server.hpp>
#include <sstream>
#include <string>
#include <thread>
#include <vector>

std::string getCurrentTimeString() {
    auto now = std::chrono::system_clock::now();
    std::time_t now_time = std::chrono::system_clock::to_time_t(now);
    std::tm now_tm = *std::localtime(&now_time);

    std::ostringstream oss;
    oss << std::put_time(&now_tm, "%Y-%m-%d %H:%M:%S");

    return oss.str();
}

#include "commands.hpp"
#include "notifications.hpp"

directmq::DirectMQNode *node;
std::shared_ptr<directmq::portal::websocket::Runnable> runnablePortalHost;

void log(const std::string &message) {
    std::cout << message << std::endl;
    std::cout.flush();
}

void sendNotification(const UniversalNotification &notification) {
    std::cout << notification.toJson() << std::endl;
    std::cout.flush();
}

void fatal(const std::string &error) {
    sendNotification(UniversalNotification::makeFatal(error));
    exit(1);
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

    auto result = directmq::portal::websocket::WebsocketServer::create(
        command.address.c_str(), command.port,
        directmq::portal::websocket::WsDataFormat::BINARY, node);

    if (result.error) {
        fatal("Failed to listen as server: " + std::string(result.error));
        exit(1);
    }

    runnablePortalHost.reset();
    runnablePortalHost = result.instance;
}

void handleConnectCommand(ConnectCommand command) {
    log("Connecting as client to " + command.address + ":" +
        std::to_string(command.port));

    auto result = directmq::portal::websocket::WebsocketClient::connect(
        command.address.c_str(), command.port, "/",
        directmq::portal::websocket::WsDataFormat::BINARY, node);

    if (result.error) {
        fatal("Failed to connect as client: " + std::string(result.error));
        exit(1);
    }

    runnablePortalHost.reset();
    runnablePortalHost = result.client;
}

void handleStopCommand(StopCommand command) {
    log("Stopping DirectMQ node: " + command.reason);

    node->closeNode(command.reason, [command]() {
        sendNotification(UniversalNotification::makeStopped(command.reason));
    });

    log("Clean exit 0");
    exit(0);
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

std::string stdinBuffer;
std::string readCommandFromStdin() {
    std::cin.seekg(0, std::cin.end);
    int length = std::cin.tellg();
    if (length < 1) {
        return "";
    }

    std::cin.seekg(0, std::cin.beg);
    std::string data(length, ' ');
    std::cin.read(const_cast<char *>(data.c_str()), length);

    stdinBuffer += data;

    if (stdinBuffer.find('\n') != std::string::npos) {
        std::string line = stdinBuffer.substr(0, stdinBuffer.find('\n'));
        stdinBuffer = stdinBuffer.substr(stdinBuffer.find('\n') + 1);
        return line;
    }

    return "";
}

void runCommandLoop() {
    try {
        std::string rawCommand = readCommandFromStdin();
        if (rawCommand.empty()) {
            return;
        }

        UniversalCommand command = UniversalCommand::fromJson(rawCommand);
        handleIncomingCommand(command);
    } catch (const std::exception &e) {
        fatal(e.what());
    }
}

void runNodeLoop() {
    if (!node || !runnablePortalHost) {
        return;
    }

    runnablePortalHost->run(0);
}

int main() {
    log("Starting DirectMQ testing agent");
    sendNotification(UniversalNotification::makeReady(getCurrentTimeString()));

    log("Starting main loop");
    while (true) {
        try {
            auto currentMillis =
                std::chrono::duration_cast<std::chrono::milliseconds>(
                    std::chrono::system_clock::now().time_since_epoch());

            runCommandLoop();
            runNodeLoop();

            std::this_thread::sleep_for(
                std::chrono::milliseconds(100) -
                (std::chrono::duration_cast<std::chrono::milliseconds>(
                     std::chrono::system_clock::now().time_since_epoch()) -
                 currentMillis));
        } catch (const std::exception &e) {
            fatal("Main loop fatal failure: " + std::string(e.what()));
        }
    }

    return 0;
}
