package mqtt

import "fmt"

func statusTopic(clientId string) string {
	return fmt.Sprintf("feeder/%s/status", clientId)
}

func feedTopic(clientId string) string {
	return fmt.Sprintf("feeder/%s/feed", clientId)
}

func feedLogTopic(clientId string) string {
	return fmt.Sprintf("feeder/%s/feed_log", clientId)
}
