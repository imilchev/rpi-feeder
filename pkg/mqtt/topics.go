package mqtt

import "fmt"

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

func wildcardOrClientId(clientId *string) string {
	c := "+"
	if clientId != nil {
		c = *clientId
	}
	return c
}
