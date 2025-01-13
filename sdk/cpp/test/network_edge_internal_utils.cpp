#include <catch2/catch.hpp>

#include "network/edge/utilities.hpp"

SCENARIO("network edge internal utils") {
    GIVEN("checkForNetworkLoops") {
        WHEN("no network loops detected") {
            auto frame = directmq::protocol::messages::DataFrame();
            frame.traversed = {"host1", "host2", "host3"};

            THEN("it should return false") {
                REQUIRE(directmq::network::edge::internal::checkForNetworkLoops(
                            frame) == false);
            }
        }

        WHEN("network loops detected") {
            auto frame = directmq::protocol::messages::DataFrame();
            frame.traversed = {"host1", "host2", "host3", "host2"};

            THEN("it should return true") {
                REQUIRE(directmq::network::edge::internal::checkForNetworkLoops(
                            frame) == true);
            }
        }
    }
}
