package directmq

import (
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

func IsCorrectTopicPattern(pattern string) bool {
	allowedTopicPatternCharactersRegex := "^[a-zA-Z0-9_*/@]+$"

	if len(pattern) == 0 {
		return false
	}

	if pattern[0] == '/' || pattern[len(pattern)-1] == '/' {
		return false
	}

	if strings.Contains(pattern, "//") {
		return false
	}

	patternWithoutOperatorsReplacer := strings.NewReplacer(
		"*", "",
		"/", "",
	)

	if len(patternWithoutOperatorsReplacer.Replace(pattern)) == 0 {
		return false
	}

	if !regexp.MustCompile(allowedTopicPatternCharactersRegex).Match([]byte(pattern)) {
		return false
	}

	_, err := glob.Compile(pattern)

	return err == nil
}

func MatchTopicPattern(pattern string, topic string) bool {
	if !IsCorrectTopicPattern(pattern) {
		return false
	}

	if !IsCorrectTopicPattern(topic) {
		return false
	}

	g := glob.MustCompile(pattern, '/')

	return g.Match(topic)
}

func IsSubtopicPattern(topLevelPattern string, target string) bool {
	return MatchTopicPattern(topLevelPattern, target)
}

func DeduplicateOverlappingTopics(topics []string) []string {
	topLevelTopics := make([]string, 0)

	for _, topic := range unique(topics) {
		isSubtopic := false

		for _, topLevelTopic := range topLevelTopics {
			if topic != topLevelTopic && IsSubtopicPattern(topLevelTopic, topic) {
				isSubtopic = true
				break
			}
		}

		if !isSubtopic {
			topLevelTopics = append(topLevelTopics, topic)
		}
	}

	for _, topic := range topLevelTopics {
		for _, otherTopic := range topLevelTopics {
			if topic != otherTopic && IsSubtopicPattern(topic, otherTopic) {
				return DeduplicateOverlappingTopics(reverse(topLevelTopics))
			}
		}
	}

	return topLevelTopics
}

func GetDeduplicatedOverlappingTopicsDiff(oldTopics []string, newTopics []string) (removed []string, added []string) {
	oldDeduplicatedTopics := DeduplicateOverlappingTopics(oldTopics)
	newDeduplicatedTopics := DeduplicateOverlappingTopics(newTopics)

	removed = make([]string, 0)
	added = make([]string, 0)

	for _, oldTopic := range oldDeduplicatedTopics {
		if !contains(newDeduplicatedTopics, oldTopic) {
			removed = append(removed, oldTopic)
			continue
		}
	}

	for _, newTopic := range newDeduplicatedTopics {
		if !contains(oldDeduplicatedTopics, newTopic) {
			added = append(added, newTopic)
		}
	}

	return
}
