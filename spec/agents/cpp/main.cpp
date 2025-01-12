#include <atomic>
#include <chrono>
#include <cstdint>
#include <directmq.hpp>
#include <iomanip>
#include <iostream>
#include <portals/streams/tcp_portal_client.hpp>
#include <portals/streams/tcp_portal_server.hpp>
#include <sstream>
#include <string>

#include <boost/archive/iterators/base64_from_binary.hpp>
#include <boost/archive/iterators/binary_from_base64.hpp>
#include <boost/archive/iterators/transform_width.hpp>

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

const int NO_EXIT = -1;
std::atomic<int> exitFlag(NO_EXIT);

std::shared_ptr<directmq::DirectMQNode> node;
directmq::portal::streams::TcpPortalClient::Pointer client = nullptr;
directmq::portal::streams::TcpPortalServer::Pointer server = nullptr;

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
}

void fatal(const std::string &error) {
    sendNotification(UniversalNotification::makeFatal(error));
    exitAgent(1);
}

std::string fixBase64Padding(const std::string &input) {
    std::string fixed = input;
    size_t len = fixed.length();

    if (len >= 2 && fixed.substr(len - 2) == "AA") {
        fixed.replace(len - 2, 2, "==");
    } else if (len >= 1 && fixed.substr(len - 1) == "A") {
        fixed.replace(len - 1, 1, "=");
    }

    return fixed;
}

std::string base64Encode(const std::string &input) {
    using namespace boost::archive::iterators;
    using It = base64_from_binary<transform_width<std::string::const_iterator, 6, 8>>;

    std::string encoded(It(std::begin(input)), It(std::end(input)));
    size_t padding = (3 - input.length() % 3) % 3;
    encoded.append(padding, '=');
    return encoded;
}

std::string base64Decode(const std::string &input) {
    using namespace boost::archive::iterators;
    using It = transform_width<binary_from_base64<std::string::const_iterator>, 8, 6>;

    // Remove padding characters
    std::string decoded(input);
    decoded.erase(std::remove(decoded.begin(), decoded.end(), '='), decoded.end());

    try {
        std::string output(It(std::begin(decoded)), It(std::end(decoded)));
        return output;
    } catch (const std::exception &e) {
        throw std::runtime_error("Invalid base64 input");
    }
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
            auto encoded = base64Encode(std::string(publication.payload.begin(),
                            publication.payload.end()));

            sendNotification(UniversalNotification::makeOnPublication(
                publication.frame.ttl,
                std::vector<std::string>(publication.frame.traversed.begin(),
                                         publication.frame.traversed.end()),
                publication.topic, publication.deliveryStrategy,
                encoded));
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

    node = std::shared_ptr<directmq::DirectMQNode>(
        new directmq::DirectMQNode(config));

    registerDiagnosticsHandlers();

    log("Setup complete");
}

std::pair<std::string, uint_least16_t> parseAddress(
    const std::string &address) {
    std::size_t colonPos = address.rfind(':');
    if (colonPos == std::string::npos) {
        throw std::invalid_argument(
            "Invalid address format. Expected format: host:port");
    }

    const std::string TCP_PROTOCOL = "tcp://";
    std::string host = address.substr(TCP_PROTOCOL.length(), colonPos - TCP_PROTOCOL.length());
    std::string rawPort = address.substr(colonPos + 1, address.length() - colonPos - 2);

    uint_least16_t port = std::stoi(rawPort);

    return {host, port};
}

void handleListenCommand(ListenCommand command) {
    log("Listening as server at " + command.address);

    auto [host, port] = parseAddress(command.address);

    server = directmq::portal::streams::TcpPortalServer::create(node, port, 0);

    server->start();
}

void handleConnectCommand(ConnectCommand command) {
    log("Connecting as client to " + command.address);

    auto [host, port] = parseAddress(command.address);

    client =
        directmq::portal::streams::TcpPortalClient::connect(node, host, port);
}

void handleStopCommand(StopCommand command) {
    log("Stopping DirectMQ node: " + command.reason);

    node->closeNode(command.reason, [command]() {
        sendNotification(UniversalNotification::makeStopped(command.reason));
    });

    if (server) {
        server->stop();
    }

    if (client) {
        client->close();
    }

    log("Clean exit 0");
    exitAgent(0);
}

void handlePublishCommand(PublishCommand command) {
    log("Publishing message to topic " + command.topic);

    auto decodedPayload = base64Decode(command.payload);
    std::vector<uint8_t> payload(decodedPayload.begin(),
                                 decodedPayload.end());
    node->publish(command.topic, payload, command.deliveryStrategy);
}

void handleSubscribeCommand(SubscribeTopicCommand command) {
    log("Subscribing to topic " + command.topic);

    auto subscriptionId = node->subscribe(
        command.topic, [](const std::string &topic,
                                 const std::vector<uint8_t> &payload) {
            sendNotification(UniversalNotification::makeMessageReceived(
                topic, base64Encode(std::string(payload.begin(), payload.end()))));
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
                continue;
            }

            UniversalCommand command = UniversalCommand::fromJson(rawCommand);
            handleIncomingCommand(command);
        } catch (const std::exception &e) {
            fatal("Command loop fatal failure: " + std::string(e.what()));
        }
    }
}

int main() {
    try {
        log("Starting DirectMQ C++ SDK testing agent");

        log("Agent ready");
        sendNotification(
            UniversalNotification::makeReady(getCurrentTimeString()));

        log("Starting command loop");
        runCommandLoop();

        log("Exiting agent with code " + std::to_string(exitFlag));
        return exitFlag;
    } catch (const std::exception &e) {
        fatal("Fatal failure: " + std::string(e.what()));
        return 1;
    }
}
