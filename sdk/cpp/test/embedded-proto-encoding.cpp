#include <catch2/catch.hpp>

#include "protocol/embedded/decoder.hpp"
#include "protocol/embedded/encoder.hpp"
#include "utils/payload_from_string.hpp"
#include "utils/test_handler.hpp"
#include "utils/test_portal.hpp"
#include "utils/test_writer.hpp"

TEST_CASE(
    "supported protocol versions message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::SupportedProtocolVersionsMessage message{
        .frame = frame, .supportedVersions = {1}};

    bool called = false;
    handler.setOnSupportedProtocolVersions(
        [&](const messages::SupportedProtocolVersionsMessage& message) {
            REQUIRE(message.supportedVersions.size() == 1);
            REQUIRE(*message.supportedVersions.begin() == 1);
            called = true;
        });

    auto encodingResult = encoder.supportedProtocolVersions(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("init connection message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::InitConnectionMessage message{.frame = frame,
                                            .maxMessageSize = 1024};

    bool called = false;
    handler.setOnInitConnection(
        [&](const messages::InitConnectionMessage& message) {
            REQUIRE(message.maxMessageSize == 1024);
            called = true;
        });

    auto encodingResult = encoder.initConnection(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("connection accepted message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::ConnectionAcceptedMessage message{.frame = frame,
                                                .maxMessageSize = 1024};

    bool called = false;
    handler.setOnConnectionAccepted(
        [&](const messages::ConnectionAcceptedMessage& message) {
            REQUIRE(message.maxMessageSize == 1024);
            called = true;
        });

    auto encodingResult = encoder.connectionAccepted(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("gracefully close message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::GracefullyCloseMessage message{.frame = frame,
                                             .reason = "some reason"};

    bool called = false;
    handler.setOnGracefullyClose(
        [&](const messages::GracefullyCloseMessage& message) {
            REQUIRE(message.reason == "some reason");
            called = true;
        });

    auto encodingResult = encoder.gracefullyClose(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("terminate network message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::TerminateNetworkMessage message{.frame = frame,
                                              .reason = "some reason"};

    bool called = false;
    handler.setOnTerminateNetwork(
        [&](const messages::TerminateNetworkMessage& message) {
            REQUIRE(message.reason == "some reason");
            called = true;
        });

    auto encodingResult = encoder.terminateNetwork(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("publish message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::PublishMessage message{
        .frame = frame,
        .topic = "some/topic",
        .deliveryStrategy = messages::DeliveryStrategy::AT_MOST_ONCE,
        .payload = payloadFromString("some payload")};

    bool called = false;
    handler.setOnPublish([&](const messages::PublishMessage& message) {
        REQUIRE(message.topic == "some/topic");
        REQUIRE(message.deliveryStrategy ==
                messages::DeliveryStrategy::AT_MOST_ONCE);
        REQUIRE(message.payload == payloadFromString("some payload"));
        called = true;
    });

    auto encodingResult = encoder.publish(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("subscribe message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::SubscribeMessage message{.frame = frame, .topic = "some/topic"};

    bool called = false;
    handler.setOnSubscribe([&](const messages::SubscribeMessage& message) {
        REQUIRE(message.topic == "some/topic");
        called = true;
    });

    auto encodingResult = encoder.subscribe(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("unsubscribe message embedded encoding and decoding") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{.ttl = 32, .traversed = std::list<std::string>()};

    messages::UnsubscribeMessage message{.frame = frame, .topic = "some/topic"};

    bool called = false;
    handler.setOnUnsubscribe([&](const messages::UnsubscribeMessage& message) {
        REQUIRE(message.topic == "some/topic");
        called = true;
    });

    auto encodingResult = encoder.unsubscribe(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
    REQUIRE(called);
}

TEST_CASE("embedded decoder should report malformed message") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();
    auto writer = portal.beginWrite(4);

    auto malformed = new uint8_t[4]{1, 2, 3, 4};
    writer->write(malformed, 4);
    writer->end();

    delete[] malformed;
    delete writer;

    bool callbackCalled = false;

    handler.setOnMalformedMessage(
        [&](const messages::MalformedMessage& message) {
            std::vector<uint8_t> expectedMessage({1, 2, 3, 4});
            REQUIRE(message.bytes == expectedMessage);
            REQUIRE(message.error.empty() == false);
            callbackCalled = true;
        });

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error != nullptr);
    REQUIRE(callbackCalled == true);
}

TEST_CASE("embedded encoder correctly encodes traversed nodes correctly") {
    auto portal = TestPortal();
    auto handler = TestHandler();

    auto encoder = embedded::EmbeddedProtocolEncoderImplementation();
    auto decoder = embedded::EmbeddedProtocolDecoderImplementation();

    messages::DataFrame frame{
        .ttl = 32, .traversed = {"node1", "node2", "node3", "node4"}};

    messages::SupportedProtocolVersionsMessage message{
        .frame = frame, .supportedVersions = {1}};

    handler.setOnSupportedProtocolVersions(
        [&](const messages::SupportedProtocolVersionsMessage& message) {
            REQUIRE(message.frame.traversed == frame.traversed);
        });

    auto encodingResult = encoder.supportedProtocolVersions(message, portal);
    REQUIRE(encodingResult.error == nullptr);

    auto decodingResult = decoder.readMessage(&portal, &handler);

    REQUIRE(decodingResult.error == nullptr);
}