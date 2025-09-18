package messaging

import "regexp"

func replacePassword(connectionString string) (string, error) {

	re := regexp.MustCompile(`amqp://[^:/]+:([^@]+)@`)

	return re.ReplaceAllString(connectionString, "amqp://user:******@"), nil
}
