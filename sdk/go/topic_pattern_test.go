package directmq

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Topic pattern util", func() {
	Context("IsCorrectTopicPattern", func() {
		var forbiddenCharacters []string = []string{"!", "#", "%", "^", "&", "(", ")", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", "?", " ", "\t", "\n", "\r", "\f", "\v", "`", "~", "-"}

		DescribeTable(
			"returns correct result for topic",
			func(topic string, expected bool) {
				Expect(IsCorrectTopicPattern(topic)).To(Equal(expected))
			},
			Entry("single level topic", "topic", true),
			Entry("multi level topic", "topic/level", true),
			Entry("multi level topic with wildcard", "topic/level/*", true),
			Entry("multi level topic with wildcard in the middle", "topic/*/somewhere", true),
			Entry("empty topic", "", false),
			Entry("topic without any characters but space", " ", false),
			Entry("topic without any characters but separator", "/", false),
			Entry("topic without any characters but separators", "//", false),
			Entry("topic that starts with a separator", "/topic", false),
			Entry("topic that ends with a separator", "topic/", false),
			Entry("allows single level wildcard only", "*", false),
			Entry("allows super wildcard only", "**", false),
			Entry("topic with empty subtopic segment", "topic//subtopic", false),
		)

		It("should return false for topics with forbidden characters", func() {
			for _, forbiddenCharacter := range forbiddenCharacters {
				Expect(IsCorrectTopicPattern("topic" + forbiddenCharacter)).To(BeFalse())
			}
		})
	})

	Context("MatchTopicPattern", func() {
		describe := func(desc string) func(pattern string, topic string, expected bool) string {
			return func(pattern string, topic string, expected bool) string {
				operator := "NOT IN"
				if expected {
					operator = "IN"
				}

				return fmt.Sprintf("%s: %s %s %s", desc, topic, operator, pattern)
			}
		}

		DescribeTable(
			"returns correct result for topic",
			func(pattern string, topic string, expected bool) {
				Expect(MatchTopicPattern(pattern, topic)).To(Equal(expected))
			},

			Entry(describe("single level topic"), "topic", "topic", true),
			Entry(describe("single level topic"), "topic", "other", false),

			Entry(describe("multi level topic with wildcard end"), "topic/*", "topic/level", true),
			Entry(describe("multi level topic with wildcard end"), "topic/*", "topic/level/sublevel", false),
			Entry(describe("multi level topic with multiple wildcards"), "topic/*/*/sublevel", "topic/level1/level2/sublevel", true),

			Entry(describe("multi level topic with super wildcard"), "topic/**/sublevel", "topic/level/something/sublevel", true),
			Entry(describe("multi level topic with multiple super wildcards"), "topic/**/sublevel/**", "topic/level/something/sublevel/other", true),
			Entry(describe("multi level topic with multiple super wildcards"), "topic/**/sublevel/**", "topic/level/something/x/other", false),

			Entry(describe("multi level with mixed wildcards and super wildcards"), "topic/*/sublevel/**", "topic/level/sublevel/other", true),
			Entry(describe("multi level with mixed wildcards and super wildcards"), "topic/*/sublevel/**", "topic/level1/level2/sublevel/other/levels", false),
		)
	})

	Context("GetDeduplicatedOverlappingTopicsDiff", func() {
		It("should return the correct removed and added topics", func() {
			oldTopics := []string{"topic1", "topic2", "topic3"}
			newTopics := []string{"topic2", "topic3", "topic4"}

			removed, added := GetDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics)

			Expect(removed).To(ConsistOf("topic1"))
			Expect(added).To(ConsistOf("topic4"))
		})

		It("should handle empty input topics", func() {
			oldTopics := []string{}
			newTopics := []string{"topic1", "topic2"}

			removed, added := GetDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics)

			Expect(removed).To(BeEmpty())
			Expect(added).To(ConsistOf("topic1", "topic2"))
		})

		It("should handle overlapping topics with duplicates", func() {
			oldTopics := []string{"topic1", "topic2", "topic2", "topic3"}
			newTopics := []string{"topic2", "topic3", "topic3", "topic4"}

			removed, added := GetDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics)

			Expect(removed).To(ConsistOf("topic1"))
			Expect(added).To(ConsistOf("topic4"))
		})

		It("should handle no overlapping topics", func() {
			oldTopics := []string{"topic1", "topic2", "topic3"}
			newTopics := []string{"topic4", "topic5", "topic6"}

			removed, added := GetDeduplicatedOverlappingTopicsDiff(oldTopics, newTopics)

			Expect(removed).To(ConsistOf("topic1", "topic2", "topic3"))
			Expect(added).To(ConsistOf("topic4", "topic5", "topic6"))
		})
	})
})
