
#include "topic_pattern.hpp"

#include <catch2/catch.hpp>
#include <vector>

SCENARIO("topic patterns") {
    GIVEN("isCorrectTopicPattern") {
        WHEN("topic pattern is correct") {
            auto pattern = GENERATE(
                "topic", "topic/level", "topic/level/*", "topic/*/somewhere",
                "*", "**", "topic/**", "topic/**/somewhere",
                "topic/**/somewhere/*", "*/topic", "*/topic/*", "*/topic/**",
                "*/topic/**/somewhere", "@extension/topic", "$system/topic",
                "topic_with___and_123");

            THEN("it should return true in case of correct pattern") {
                INFO("Pattern: " << pattern);
                REQUIRE(directmq::topics::isCorrectTopicPattern(pattern) ==
                        true);
            }
        }

        WHEN("topic pattern is incorrect") {
            auto pattern =
                GENERATE("", " ", "/", "//", "/topic", "topic/", "topic//level",
                         "topic/**/**", "some/topic*", "some/topic**",
                         "some/topic/***", "some/*topic", "some/**topic",
                         "some/*topic*/", "some/**topic**");

            THEN("it should return false in case of incorrect pattern") {
                INFO("Pattern: " << pattern);
                REQUIRE(directmq::topics::isCorrectTopicPattern(pattern) ==
                        false);
            }
        }

        WHEN("forbidded character is used") {
            auto forbiddenCharacter = GENERATE(
                "!", "#", "%", "^", "&", "(", ")", "+", "=", "{", "}", "[", "]",
                "|", "\\", ":", ";", "\"", "'", "<", ">", ",", "?", " ", "\t",
                "\n", "\r", "\f", "\v", "`", "~", "-");

            THEN("it should return false in case of forbidden character") {
                INFO("Pattern: topic/" << forbiddenCharacter);
                REQUIRE(directmq::topics::isCorrectTopicPattern(
                            std::string("topic/") + forbiddenCharacter) ==
                        false);
            }
        }
    }

    GIVEN("matchTopicPattern") {
        struct TopicPatternTestCase {
            std::string pattern;
            std::string topic;
        };

        WHEN("simple matching topic patterns provided") {
            auto testCase = GENERATE(
                TopicPatternTestCase{"topic", "topic"},
                TopicPatternTestCase{"topic/level", "topic/level"},
                TopicPatternTestCase{"topic/level/1", "topic/level/1"});

            THEN("it should return true") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == true);
            }
        }

        // WHEN("simple not matching topic patterns provided") {
        //     auto testCase =
        //         GENERATE(TopicPatternTestCase{"topic", "another"},
        //                  TopicPatternTestCase{"topic", "topic/level"},
        //                  TopicPatternTestCase{"topic/level", "topic"},
        //                  TopicPatternTestCase{"topic/level", "level"},
        //                  TopicPatternTestCase{"topic/level", "topic/level/1"});

        //     THEN("it should return false") {
        //         INFO("Pattern: " << testCase.pattern);
        //         INFO("Topic: " << testCase.topic);
        //         REQUIRE(directmq::topics::matchTopicPattern(
        //                     testCase.pattern, testCase.topic) == false);
        //     }
        // }

        // WHEN("single wildcard matching topic patterns provided") {
        //     auto testCase = GENERATE(
        //         TopicPatternTestCase{"*", "topic"},
        //         TopicPatternTestCase{"topic/*", "topic/level"},
        //         TopicPatternTestCase{"topic/*/1", "topic/level/1"},
        //         TopicPatternTestCase{"topic/*/*", "topic/level/1"},
        //         TopicPatternTestCase{"*/topic", "test/topic"},
        //         TopicPatternTestCase{"*/topic/*", "test/topic/level"},
        //         TopicPatternTestCase{"*/topic/*/*", "test/topic/level/1"},
        //         TopicPatternTestCase{"*/topic/*/another",
        //                              "test/topic/1/another"});

        //     THEN("it should return true") {
        //         INFO("Pattern: " << testCase.pattern);
        //         INFO("Topic: " << testCase.topic);
        //         REQUIRE(directmq::topics::matchTopicPattern(
        //                     testCase.pattern, testCase.topic) == true);
        //     }
        // }

        // WHEN("single wildcard not matching topic patterns provided") {
        //     auto testCase =
        //         GENERATE(TopicPatternTestCase{"*", "topic/level"},
        //                  TopicPatternTestCase{"topic/*", "topic"},
        //                  TopicPatternTestCase{"topic/*", "topic/level/1"},
        //                  TopicPatternTestCase{"topic/*/*", "topic/level"},
        //                  TopicPatternTestCase{"*/topic", "topic"},
        //                  TopicPatternTestCase{"*/topic/*", "topic/123"},
        //                  TopicPatternTestCase{"*/topic/*/*", "topic"});

        //     THEN("it should return false") {
        //         INFO("Pattern: " << testCase.pattern);
        //         INFO("Topic: " << testCase.topic);
        //         REQUIRE(directmq::topics::matchTopicPattern(
        //                     testCase.pattern, testCase.topic) == false);
        //     }
        // }

        // WHEN("double wildcard matching topic patterns provided") {
        //     auto testCase = GENERATE(
        //         TopicPatternTestCase{"**", "topic"},
        //         TopicPatternTestCase{"**", "topic/level/three"},
        //         TopicPatternTestCase{"topic/**/1", "topic/level/1"},
        //         TopicPatternTestCase{"**/topic", "test/topic"},
        //         TopicPatternTestCase{"**/topic/**", "test/1/topic/level/2"},
        //         TopicPatternTestCase{"part1/**/part2",
        //                              "part1/some/other/segments/part2"});

        //     THEN("it should return true") {
        //         INFO("Pattern: " << testCase.pattern);
        //         INFO("Topic: " << testCase.topic);
        //         REQUIRE(directmq::topics::matchTopicPattern(
        //                     testCase.pattern, testCase.topic) == true);
        //     }
        // }

        // WHEN("double wildcard not matching topic patterns provided") {
        //     auto testCase = GENERATE(
        //         TopicPatternTestCase{"topic/**/1", "topic/level"},
        //         TopicPatternTestCase{"**/topic", "topic/pattern"},
        //         TopicPatternTestCase{"**/topic/**", "test/1/level/2"},
        //         TopicPatternTestCase{"part1/**/part2", "part1/part2"});

        //     THEN("it should return false") {
        //         INFO("Pattern: " << testCase.pattern);
        //         INFO("Topic: " << testCase.topic);
        //         REQUIRE(directmq::topics::matchTopicPattern(
        //                     testCase.pattern, testCase.topic) == false);
        //     }
        // }
    }
}