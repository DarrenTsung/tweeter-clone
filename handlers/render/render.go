package render

import (
	"encoding/json"
	"net/http"
	"tweeter/handlers/responses"
	"tweeter/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Response renders a response from the resp provided
func Response(endpointName string, w http.ResponseWriter, statusCode int, resp interface{}) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		// This is an error because resp is controlled by the programmer and
		// should be correct in all situations
		logrus.WithFields(logrus.Fields{"err": err}).Error("Resp passed to render was not json.Marshalable")
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(respBytes); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("Failed to write bytes to http.ResponseWriter")
	}

	metrics.APIResponses.With(prometheus.Labels{"endpointName": endpointName, "code": string(statusCode)}).Inc()
}

// ErrorResponse renders the error response with the status code provided
func ErrorResponse(endpointName string, w http.ResponseWriter, statusCode int, errors ...responses.Error) {
	Response(endpointName, w, statusCode, responses.NewErrorResponse(errors...))
}
