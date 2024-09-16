
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

        WHEN("simple not matching topic patterns provided") {
            auto testCase =
                GENERATE(TopicPatternTestCase{"topic", "another"},
                         TopicPatternTestCase{"topic", "topic/level"},
                         TopicPatternTestCase{"topic/level", "topic"},
                         TopicPatternTestCase{"topic/level", "level"},
                         TopicPatternTestCase{"topic/level/1", "topic/level/2"},
                         TopicPatternTestCase{"topic/level", "topic/level/1"});

            THEN("it should return false") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == false);
            }
        }

        WHEN("single wildcard matching topic patterns provided") {
            auto testCase = GENERATE(
                TopicPatternTestCase{"*", "topic"},
                TopicPatternTestCase{"topic/*", "topic/level"},
                TopicPatternTestCase{"topic/*/1", "topic/level/1"},
                TopicPatternTestCase{"topic/*/*", "topic/level/1"},
                TopicPatternTestCase{"*/topic", "test/topic"},
                TopicPatternTestCase{"*/topic/*", "test/topic/level"},
                TopicPatternTestCase{"*/topic/*/*", "test/topic/level/1"},
                TopicPatternTestCase{"*/topic/*/another",
                                     "test/topic/1/another"});

            THEN("it should return true") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == true);
            }
        }

        WHEN("single wildcard not matching topic patterns provided") {
            auto testCase =
                GENERATE(TopicPatternTestCase{"*", "topic/level"},
                         TopicPatternTestCase{"topic/*", "topic"},
                         TopicPatternTestCase{"topic/*", "topic/level/1"},
                         TopicPatternTestCase{"topic/*/*", "topic/level"},
                         TopicPatternTestCase{"*/topic", "topic"},
                         TopicPatternTestCase{"*/topic/*", "topic/123"},
                         TopicPatternTestCase{"*/topic/*/*", "topic"});

            THEN("it should return false") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == false);
            }
        }

        WHEN("double wildcard matching topic patterns provided") {
            auto testCase = GENERATE(
                TopicPatternTestCase{"**", "topic"},
                TopicPatternTestCase{"**", "topic/level/three"},
                TopicPatternTestCase{"topic/**/1", "topic/level/1"},
                TopicPatternTestCase{"**/topic", "test/topic"},
                TopicPatternTestCase{"**/topic/**", "test/1/topic/level/2"},
                TopicPatternTestCase{"part1/**/part2",
                                     "part1/some/other/segments/part2"});

            THEN("it should return true") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == true);
            }
        }

        WHEN("double wildcard not matching topic patterns provided") {
            auto testCase = GENERATE(
                TopicPatternTestCase{"topic/**/1", "topic/level"},
                TopicPatternTestCase{"**/topic", "topic/pattern"},
                TopicPatternTestCase{"**/topic/**", "test/1/level/2"},
                TopicPatternTestCase{"part1/**/part2", "part1/part2"});

            THEN("it should return false") {
                INFO("Pattern: " << testCase.pattern);
                INFO("Topic: " << testCase.topic);
                REQUIRE(directmq::topics::matchTopicPattern(
                            testCase.pattern, testCase.topic) == false);
            }
        }
    }

    GIVEN("isSubtopicPattern") {
        class SubtopicPatternTestCase {
            public:
                std::string topLevelPattern;
                std::string subtopicPattern;
        };

        auto testCase = GENERATE(
            // simple single wildcard
            SubtopicPatternTestCase{"*", "topic"},
            SubtopicPatternTestCase{"*/*", "topic/subtopic"},
            SubtopicPatternTestCase{"topic/*", "topic/subtopic"},

            // complex single wildcard
            SubtopicPatternTestCase{"*/*", "topic/*"},
            SubtopicPatternTestCase{"*/*", "*/topic"},
            SubtopicPatternTestCase{"topic/*/*/*", "topic/*/level/*"},

            // simple double wildcard
            SubtopicPatternTestCase{"**", "topic"},
            SubtopicPatternTestCase{"**", "topic/subtopic"},
            SubtopicPatternTestCase{"part1/**/part2", "part1/some/topics/inside/part2"},
            SubtopicPatternTestCase{"**/topic/**", "some/topic/inside"},
            SubtopicPatternTestCase{"**/topic/**", "some/random/topic/inside/string"},
            SubtopicPatternTestCase{"part1/**/part2/**/part3", "part1/x1/x2/part2/x3/x4/part3"},

            // complex double wildcard
            SubtopicPatternTestCase{"topic/**", "topic/**/subtopic"},
            SubtopicPatternTestCase{"**/subtopic", "topic/**/subtopic"},
            SubtopicPatternTestCase{"topic/**/subtopic", "topic/**/other/subtopic"},
            SubtopicPatternTestCase{"topic/**/subtopic", "topic/**/other/**/subtopic"},

            // mixed single and double wildcard
            SubtopicPatternTestCase{"topic/**/subtopic/*", "topic/**/subtopic/level"},
            SubtopicPatternTestCase{"topic/**/subtopic", "topic/*/subtopic"},
            SubtopicPatternTestCase{"topic/**/subtopic", "topic/*/something/subtopic"},
            SubtopicPatternTestCase{"topic/**/subtopic/*", "topic/*/subtopic/*"});

        WHEN("real subtopic patterns provided") {
            THEN("it should return true") {
                INFO("Top level pattern: " << testCase.topLevelPattern);
                INFO("Subtopic pattern: " << testCase.subtopicPattern);
                REQUIRE(directmq::topics::isSubtopicPattern(
                            testCase.topLevelPattern, testCase.subtopicPattern) ==
                        true);
            }
        }

        WHEN("fake subtopic patterns provided") {
            THEN("it should return false") {
                INFO("Top level pattern: " << testCase.subtopicPattern);
                INFO("Subtopic pattern: " << testCase.topLevelPattern);
                REQUIRE(directmq::topics::isSubtopicPattern(
                            testCase.subtopicPattern, testCase.topLevelPattern) ==
                        false);
            }
        }
    }

    GIVEN("determineTopLevelTopicPattern") {
        WHEN("left topic pattern is the top level") {
            THEN("it should return LEFT") {
                REQUIRE(directmq::topics::determineTopLevelTopicPattern(
                            "topic/*", "topic/level") ==
                        directmq::topics::DetermineTopLevelTopicPatternResult::LEFT);
            }
        }

        WHEN("right topic pattern is the top level") {
            THEN("it should return RIGHT") {
                REQUIRE(directmq::topics::determineTopLevelTopicPattern(
                            "topic/level", "topic/*") ==
                        directmq::topics::DetermineTopLevelTopicPatternResult::RIGHT);
            }
        }

        WHEN("no top level topic pattern") {
            THEN("it should return NONE") {
                REQUIRE(directmq::topics::determineTopLevelTopicPattern(
                            "topic/level", "topic/another") ==
                        directmq::topics::DetermineTopLevelTopicPatternResult::NONE);
            }
        }
    }

    GIVEN("deduplicateOverlappingTopics") {
        WHEN("no overlapping topics provided") {
            auto topics = directmq::topics::internal::makeTopicsList({"topic/1", "topic/2", "topic/3"});

            THEN("it should return the same list") {
                auto deduplicatedTopics = directmq::topics::deduplicateOverlappingTopics(topics);
                auto result = directmq::topics::internal::compareTopicLists(topics, deduplicatedTopics);
                REQUIRE(result == true);
            }
        }

        WHEN("overlapping topics provided") {
            auto topics = directmq::topics::internal::makeTopicsList({"topic/1", "topic/1/level", "topic/2", "topic/*", "topic/3", "topic/1/level"});

            THEN("it should return deduplicated list") {
                auto deduplicatedTopics = directmq::topics::deduplicateOverlappingTopics(topics);
                auto result = directmq::topics::internal::compareTopicLists(
                    directmq::topics::internal::makeTopicsList({"topic/*", "topic/1/level"}), deduplicatedTopics);

                REQUIRE(result == true);
            }
        }
    }

    GIVEN("getDeduplicatedOverlappingTopicsDiff") {
        WHEN("no overlapping topics provided") {
            auto oldTopics = directmq::topics::internal::makeTopicsList({"topic/1", "topic/2", "topic/3"});
            auto newTopics = directmq::topics::internal::makeTopicsList({"topic/4", "topic/5", "topic/6"});

            THEN("it should return empty diff") {
                auto diff = directmq::topics::getDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics);

                REQUIRE(directmq::topics::internal::compareTopicLists(diff.added, newTopics) == true);
                REQUIRE(directmq::topics::internal::compareTopicLists(diff.removed, oldTopics) == true);
            }
        }

        WHEN("overlapping topics provided") {
            auto oldTopics = directmq::topics::internal::makeTopicsList({"topic/0", "topic/1", "topic/2", "topic/3"});
            auto newTopics = directmq::topics::internal::makeTopicsList({"topic/2", "topic/3", "topic/4", "topic/5"});

            THEN("it should return diff with added and removed topics") {
                auto diff = directmq::topics::getDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics);

                std::vector<std::string> addedTopics = {"topic/4", "topic/5"};
                std::vector<std::string> removedTopics = {"topic/0", "topic/1"};

                auto addedResult = directmq::topics::internal::compareTopicLists(
                    directmq::topics::internal::makeTopicsList(addedTopics), diff.added);

                auto removedResult = directmq::topics::internal::compareTopicLists(
                    directmq::topics::internal::makeTopicsList(removedTopics), diff.removed);

                REQUIRE(addedResult == true);
                REQUIRE(removedResult == true);
            }
        }
    }
}
