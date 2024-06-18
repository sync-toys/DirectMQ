#include <catch2/catch.hpp>
#include <queue>

#include "portal.hpp"
#include "protocol/decoder.hpp"
#include "protocol/embedded/decoder.hpp"
#include "protocol/embedded/encoder.hpp"
#include "protocol/encoder.hpp"
#include "protocol/messages.hpp"

SCENARIO("basic encoding and decoding", "[protocol]") {
    GIVEN("supported protocol versions message") {
        WHEN("should succeed") { REQUIRE(true); }

        WHEN("correct message provided") {
            using namespace directmq;
            using namespace directmq::protocol;

            class TestHandler : public DecodingHandler {
               private:
                std::function<void(
                    const messages::SupportedProtocolVersionsMessage&)>
                    onSupportedProtocolVersionsCallback;

                std::function<void(const messages::InitConnectionMessage&)>
                    onInitConnectionCallback;

                std::function<void(const messages::ConnectionAcceptedMessage&)>
                    onConnectionAcceptedCallback;

                std::function<void(const messages::GracefullyCloseMessage&)>
                    onGracefullyCloseCallback;

                std::function<void(const messages::TerminateNetworkMessage&)>
                    onTerminateNetworkCallback;

                std::function<void(const messages::PublishMessage&)>
                    onPublishCallback;

                std::function<void(const messages::SubscribeMessage&)>
                    onSubscribeCallback;

                std::function<void(const messages::UnsubscribeMessage&)>
                    onUnsubscribeCallback;

                std::function<void(const messages::MalformedMessage&)>
                    onMalformedMessageCallback;

               public:
                void setOnSupportedProtocolVersions(
                    std::function<
                        void(const messages::SupportedProtocolVersionsMessage&)>
                        onSupportedProtocolVersions) {
                    this->onSupportedProtocolVersionsCallback =
                        onSupportedProtocolVersions;
                }
                void onSupportedProtocolVersions(
                    const messages::SupportedProtocolVersionsMessage& message) {
                    if (this->onSupportedProtocolVersionsCallback == nullptr) {
                        return;
                    }

                    this->onSupportedProtocolVersionsCallback(message);
                }

                void setOnInitConnection(
                    std::function<void(const messages::InitConnectionMessage&)>
                        onInitConnection) {
                    this->onInitConnectionCallback = onInitConnection;
                }
                void onInitConnection(
                    const messages::InitConnectionMessage& message) {
                    if (this->onInitConnectionCallback == nullptr) {
                        return;
                    }

                    this->onInitConnectionCallback(message);
                }

                void setOnConnectionAccepted(
                    std::function<
                        void(const messages::ConnectionAcceptedMessage&)>
                        onConnectionAccepted) {
                    this->onConnectionAcceptedCallback = onConnectionAccepted;
                }
                void onConnectionAccepted(
                    const messages::ConnectionAcceptedMessage& message) {
                    if (this->onConnectionAcceptedCallback == nullptr) {
                        return;
                    }

                    this->onConnectionAcceptedCallback(message);
                }

                void setOnGracefullyClose(
                    std::function<void(const messages::GracefullyCloseMessage&)>
                        onGracefullyClose) {
                    this->onGracefullyCloseCallback = onGracefullyClose;
                }
                void onGracefullyClose(
                    const messages::GracefullyCloseMessage& message) {
                    if (this->onGracefullyCloseCallback == nullptr) {
                        return;
                    }

                    this->onGracefullyCloseCallback(message);
                }

                void setOnTerminateNetwork(
                    std::function<
                        void(const messages::TerminateNetworkMessage&)>
                        onTerminateNetwork) {
                    this->onTerminateNetworkCallback = onTerminateNetwork;
                }
                void onTerminateNetwork(
                    const messages::TerminateNetworkMessage& message) {
                    if (this->onTerminateNetworkCallback == nullptr) {
                        return;
                    }

                    this->onTerminateNetworkCallback(message);
                }

                void setOnPublish(
                    std::function<void(const messages::PublishMessage&)>
                        onPublish) {
                    this->onPublishCallback = onPublish;
                }
                void onPublish(const messages::PublishMessage& message) {
                    if (this->onPublishCallback == nullptr) {
                        return;
                    }

                    this->onPublishCallback(message);
                }

                void setOnSubscribe(
                    std::function<void(const messages::SubscribeMessage&)>
                        onSubscribe) {
                    this->onSubscribeCallback = onSubscribe;
                }
                void onSubscribe(const messages::SubscribeMessage& message) {
                    if (this->onSubscribeCallback == nullptr) {
                        return;
                    }

                    this->onSubscribeCallback(message);
                }

                void setOnUnsubscribe(
                    std::function<void(const messages::UnsubscribeMessage&)>
                        onUnsubscribe) {
                    this->onUnsubscribeCallback = onUnsubscribe;
                }
                void onUnsubscribe(
                    const messages::UnsubscribeMessage& message) {
                    if (this->onUnsubscribeCallback == nullptr) {
                        return;
                    }

                    this->onUnsubscribeCallback(message);
                }

                void setOnMalformedMessage(
                    std::function<void(const messages::MalformedMessage&)>
                        onMalformedMessage) {
                    this->onMalformedMessageCallback = onMalformedMessage;
                }
                void onMalformedMessage(
                    const messages::MalformedMessage& message) {
                    if (this->onMalformedMessageCallback == nullptr) {
                        return;
                    }

                    this->onMalformedMessageCallback(message);
                }
            };

            class TestWriter : public DataWriter {
               private:
                std::queue<Packet>* packets;
                const size_t packetSize;
                pb_byte_t* data;
                size_t currentPosition = 0;

               public:
                TestWriter(std::queue<Packet>* packets, const size_t packetSize)
                    : packets(packets),
                      packetSize(packetSize),
                      data(new pb_byte_t[packetSize]) {}

                bool write(bytes block, const size_t blockSize) {
                    for (size_t blockPos = 0; blockPos < blockSize;
                         blockPos++) {
                        data[this->currentPosition + blockPos] =
                            block[blockPos];
                    }

                    this->currentPosition += blockSize;
                    return true;
                }

                void end() {
                    Packet packet{
                        this->currentPosition,
                        this->data,
                    };

                    packets->push(packet);
                };
            };

            class TestPortal : public Portal {
               private:
                std::queue<Packet> packets;

               public:
                void close() {}

                DataWriter* beginWrite(const size_t packetSize) {
                    auto writer = new TestWriter(&this->packets, packetSize);
                    return writer;
                }

                Packet readPacket() {
                    if (this->packets.empty()) {
                        auto packet = Packet{0, nullptr};
                        return packet;
                    }

                    auto packet = packets.front();
                    packets.pop();

                    return packet;
                }
            };

            auto portal = TestPortal();
            auto handler = TestHandler();

            auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
            auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

            messages::DataFrame frame{
                .ttl = 32, .traversed = nullptr, .traversedCount = 0};

            uint32_t supportedVersions[] = {1};
            messages::SupportedProtocolVersionsMessage message{
                .frame = frame,
                .supportedVersions = supportedVersions,
                .supportedVersionsCount = 1};

            bool called = false;
            handler.setOnSupportedProtocolVersions(
                [&](const messages::SupportedProtocolVersionsMessage& message) {
                    REQUIRE(message.supportedVersionsCount == 1);
                    REQUIRE(message.supportedVersions[0] == 1);
                    called = true;
                });

            encoder.supportedProtocolVersions(message, portal);
            auto result = decoder.readMessage(&portal, &handler);

            std::cout << result.error << std::endl;
            REQUIRE(called);
        }
    }
}
