package httpclient

import (
	"fmt"
	"net/http"

	"github.com/andre/code-styles-golang/pkg/datadog/env"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var ddEnvs *env.DatadogEnvironment

func ConfigureHttpClient(e *env.DatadogEnvironment) {
	ddEnvs = e
}

func New() *http.Client {
	if !ddEnvs.DATADOG_ENABLED {
		return &http.Client{}
	}

	return httptrace.WrapClient(&http.Client{}, httptrace.RTWithResourceNamer(func(req *http.Request) string {
		return fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	}))
}

func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	if !ddEnvs.DATADOG_ENABLED {
		return rt
	}

	return httptrace.WrapRoundTripper(rt, httptrace.RTWithServiceName(ddEnvs.DD_SERVICE))
}
