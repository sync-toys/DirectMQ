#include <catch2/catch.hpp>

#include "protocol/embedded/decoder.hpp"
#include "protocol/embedded/encoder.hpp"
#include "utils/test_handler.hpp"
#include "utils/test_portal.hpp"
#include "utils/test_writer.hpp"

SCENARIO("basic encoding and decoding", "[protocol]") {
    GIVEN("supported protocol versions message") {
        WHEN("should succeed") { REQUIRE(true); }

        WHEN("correct message provided") {
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

            REQUIRE(result.error == nullptr);
            REQUIRE(called);
        }
    }
}
