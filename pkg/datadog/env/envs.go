package env

var (
	DATADOG_DOGSTATSD_PORT int = 8125
	DATADOG_AGENT_PORT     int = 8126
)

type DatadogEnvironment struct {
	DATADOG_ENABLED        bool
	DATADOG_AGENT_ADDR     string
	DATADOG_DOGSTATSD_ADDR string

	DD_SERVICE    string
	DD_ENV        string
	DD_VERSION    string
	DD_AGENT_HOST string
}
