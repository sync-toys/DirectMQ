#pragma once

#include <cctype>
#include <cstring>
#include <list>
#include <memory>
#include <string>

namespace directmq::topics {
namespace internal {
bool isPatternSeparatorOperator(char c) { return c == '/'; }

bool isPatternWildcardOperator(char c) { return c == '*'; }

bool isPatternSuperWildcardOperator(const char* c) {
    return memcmp(c, "**", 2) == 0;
}

bool isPatternExtensionOperator(char c) { return c == '@'; }

bool isPatternSystemOperator(char c) { return c == '$'; }

bool isWordsJoiner(char c) { return c == '_'; }

bool isAllowedPatternCharacter(char c) {
    return std::isalnum(c) || isPatternSeparatorOperator(c) ||
           isPatternWildcardOperator(c) || isPatternExtensionOperator(c) ||
           isPatternSystemOperator(c) || isWordsJoiner(c);
}

bool containsOnlyAllowedPatternCharacters(const std::string& pattern) {
    for (char c : pattern) {
        if (!isAllowedPatternCharacter(c)) {
            return false;
        }
    }

    return true;
}

// dirty segment name example: "topic/segment*"
// segment with wildcard and normal characters mixed
// wildcard segments can be only "*", "/*", "*/", "/*/", "**", "/**", "**/",
// "/**/"
bool includesDirtySegmentNames(const std::string& pattern) {
    for (size_t i = 0; i < pattern.size(); ++i) {
        if (pattern[i] != '*') {
            continue;
        }

        if (i != 0 && (pattern[i - 1] != '/' && pattern[i - 1] != '*')) {
            return true;
        }

        if (i != pattern.size() - 1 &&
            (pattern[i + 1] != '/' && pattern[i + 1] != '*')) {
            return true;
        }
    }

    return false;
}

class TopicSegment {
   private:
    std::shared_ptr<std::string> pattern;

    size_t firstCharacterIndex;
    size_t lastCharacterIndex;

    TopicSegment(std::shared_ptr<std::string> pattern,
                 size_t firstCharacterIndex, size_t lastCharacterIndex)
        : pattern(pattern),
          firstCharacterIndex(firstCharacterIndex),
          lastCharacterIndex(lastCharacterIndex) {}

    friend class TopicPattern;

   public:
    size_t size() const { return lastCharacterIndex - firstCharacterIndex; }

    std::string operator*() const {
        return pattern->substr(firstCharacterIndex, size());
    }

    TopicSegment& operator++() {
        if (lastCharacterIndex == pattern->size()) {
            throw std::out_of_range(
                "Cannot increment past the end of the topic pattern");
        }

        firstCharacterIndex = lastCharacterIndex + 1;

        lastCharacterIndex = pattern->find_first_of('/', firstCharacterIndex);
        if (lastCharacterIndex == std::string::npos) {
            lastCharacterIndex = pattern->size();
        }

        return *this;
    }

    bool operator==(const TopicSegment& other) const {
        if (other.pattern == pattern) {
            return firstCharacterIndex == other.firstCharacterIndex &&
                   lastCharacterIndex == other.lastCharacterIndex;
        }

        if (size() != other.size()) {
            return false;
        }

        auto patternSegment = pattern->c_str() + firstCharacterIndex;
        auto otherSegment = other.pattern->c_str() + other.firstCharacterIndex;

        auto result = memcmp(patternSegment, otherSegment, size());

        return result == 0;
    }

    bool operator!=(const TopicSegment& other) const {
        return !(*this == other);
    }

    bool isWildcard() const {
        if (size() != 1) {
            return false;
        }

        return isPatternWildcardOperator((*pattern)[firstCharacterIndex]);
    }

    bool isSuperWildcard() const {
        if (size() != 2) {
            return false;
        }

        return isPatternSuperWildcardOperator(pattern->c_str() +
                                              firstCharacterIndex);
    }

    bool isFirst() const { return firstCharacterIndex == 0; }

    bool isLast() const { return lastCharacterIndex == pattern->size(); }

    bool matches(const TopicSegment& other) const {
        if (isWildcard() || isSuperWildcard()) {
            return true;
        }

        // just compare contents
        return *this == other;
    }
};

class TopicPattern {
   private:
    std::shared_ptr<std::string> pattern;

   public:
    TopicPattern(const std::string& pattern)
        : pattern(std::make_shared<std::string>(pattern)) {}

    TopicSegment begin() {
        if (pattern->find('/') == std::string::npos) {
            return TopicSegment(pattern, 0, pattern->size());
        }

        size_t firstTopicEnd = pattern->find_first_of('/');
        return TopicSegment(pattern, 0, firstTopicEnd);
    }

    TopicSegment end() {
        if (pattern->find('/') == std::string::npos) {
            return TopicSegment(pattern, 0, pattern->size());
        }

        size_t lastTopicStart = pattern->find_last_of('/');
        return TopicSegment(pattern, lastTopicStart + 1, pattern->size());
    }
};

void skipTargetSegmentsUntilExpected(TopicSegment& segmentToSkip,
                                     const TopicSegment& expectedSegment) {
    while (!expectedSegment.matches(segmentToSkip) && !segmentToSkip.isLast()) {
        ++segmentToSkip;
    }
}
}  // namespace internal

// bool isCorrectTopic(const std::string& topic) {
//     // TODO
//     throw std::runtime_error("Not implemented");
// }

bool isCorrectTopicPattern(const std::string& topicPattern) {
    if (!internal::containsOnlyAllowedPatternCharacters(topicPattern)) {
        return false;
    }

    if (topicPattern.empty()) {
        return false;
    }

    if (topicPattern[0] == '/' ||
        topicPattern[topicPattern.size() - 1] == '/') {
        return false;
    }

    if (topicPattern.find("//") != std::string::npos) {
        return false;
    }

    if (topicPattern.find("**/**") != std::string::npos) {
        // forbidden, because it makes no sense
        return false;
    }

    if (topicPattern.find("***") != std::string::npos) {
        // forbidden, because it makes no sense
        return false;
    }

    if (internal::includesDirtySegmentNames(topicPattern)) {
        return false;
    }

    return true;
}

bool matchTopicPattern(const std::string& patternTopic,
                       const std::string& targetTopic) {
    internal::TopicPattern pattern(patternTopic);
    internal::TopicPattern target(targetTopic);

    auto patternSegment = pattern.begin();
    auto targetSegment = target.begin();

    while (patternSegment != pattern.end() && targetSegment != target.end()) {
        if (patternSegment.isSuperWildcard() && patternSegment.isLast()) {
            return true;
        }

        if (patternSegment.isSuperWildcard()) {
            ++patternSegment;
            skipTargetSegmentsUntilExpected(targetSegment, patternSegment);
            continue;
        }

        if (!patternSegment.matches(targetSegment)) {
            return false;
        }

        ++patternSegment;
        ++targetSegment;
    }

    if (!patternSegment.isLast()) {
        return false;
    }

    if (patternSegment.isWildcard() && targetSegment.isLast()) {
        return true;
    }

    if (patternSegment.isSuperWildcard()) {
        return true;
    }

    if (targetSegment.isLast()) {
        return patternSegment.matches(targetSegment);
    }

    return false;
}

// bool isSubtopicPattern(const std::string& topLevelPattern,
//                        const std::string& subtopicPattern) {
//     return false;
// }

// std::list<std::string> deduplicateOverlappingTopics(
//     const std::list<std::string>& topics) {
//     return std::list<std::string>();
// }

// struct OverlappingTopicsDiff {
//     std::list<std::string> added;
//     std::list<std::string> removed;
// };

// OverlappingTopicsDiff getDeduplicatedOverlappingTopicsDiff(
//     const std::list<std::string>& oldTopics,
//     const std::list<std::string>& newTopics) {
//     return OverlappingTopicsDiff();
// }
}  // namespace directmq::topics
