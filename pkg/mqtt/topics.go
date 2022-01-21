package mqtt

import (
	"fmt"
	"strings"
)

// StatusTopic gives the status topic for the specified clientId. If clientId
// is nil, then a wildcard topic for all clients is returned.
func StatusTopic(clientId *string) string {
	return fmt.Sprintf("feeder/%s/status", wildcardOrClientId(clientId))
}

// FeedTopic gives the feed topic for the specified clientId. If clientId
// is nil, then a wildcard topic for all clients is returned.
func FeedTopic(clientId *string) string {
	return fmt.Sprintf("feeder/%s/feed", wildcardOrClientId(clientId))
}

// FeedLogTopic gives the feed log topic for the specified clientId. If clientId
// is nil, then a wildcard topic for all clients is returned.
func FeedLogTopic(clientId *string) string {
	return fmt.Sprintf("feeder/%s/feed_log", wildcardOrClientId(clientId))
}

// ClientIdFromTopic extracts the clientId from a topic. Panics if the topic
// format is invalid.
func ClientIdFromTopic(topic string) string {
	return strings.Split(topic, "/")[1]
}

func wildcardOrClientId(clientId *string) string {
	c := "+"
	if clientId != nil {
		c = *clientId
	}
	return c
}
