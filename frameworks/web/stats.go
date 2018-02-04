package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/thoas/stats"
)

// HandleStatistics - Exports stats as a REST Endpoint in either JSON or prometheus string
func (r *RestAPI) HandleStatistics(w http.ResponseWriter, req *http.Request) {
	requestType := req.Header.Get("Accept")
	requestType = strings.ToLower(requestType)

	// Get Statistics
	metrics := r.Statistics.Data()

	if strings.Contains(requestType, "text/plain") {
		w.Header().Set("Content-Type", "text/plain")
		metricsResponse := r.metricsToPrometheus(metrics)
		w.Write([]byte(metricsResponse))
	} else {
		//Return JSON Version
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(metrics)
		w.Write(b)
	}

}

// Takes the metrics data structure and converts it to string
func (r *RestAPI) metricsToPrometheus(metrics *stats.Data) string {
	var buffer bytes.Buffer

	buffer.WriteString("# HELP figura_uptime_sec How many seconds app has been up.\n")
	buffer.WriteString("# TYPE figura_uptime_sec counter\n")
	buffer.WriteString(fmt.Sprintf("figura_uptime_sec %v\n", metrics.UpTimeSec))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP figura_total_response_time_sec Total time spent in handling requests.\n")
	buffer.WriteString("# TYPE figura_total_response_time_sec counter\n")
	buffer.WriteString(fmt.Sprintf("figura_total_response_time_sec %v\n", metrics.TotalResponseTimeSec))
	buffer.WriteString("\n")

	buffer.WriteString("# HELP figura_average_response_time_sec Average time spent in handling requests.\n")
	buffer.WriteString("# TYPE figura_average_response_time_sec guage\n")
	buffer.WriteString(fmt.Sprintf("figura_average_response_time_sec %v\n", metrics.AverageResponseTimeSec))
	buffer.WriteString("\n")

	// Work around for bug in underlying stats library code
	calls := 0

	// Iterate through individual request counts
	if len(metrics.TotalStatusCodeCount) > 0 {
		for statuskey := range metrics.TotalStatusCodeCount {
			tally := metrics.TotalStatusCodeCount[statuskey]
			buffer.WriteString(fmt.Sprintf("# HELP figura_response_status_%v Total Number of Requests returning http status %v\n", statuskey, statuskey))
			buffer.WriteString(fmt.Sprintf("# TYPE figura_response_status_%v counter\n", statuskey))
			buffer.WriteString(fmt.Sprintf("figura_response_status_%v %v\n", statuskey, tally))
			buffer.WriteString("\n")

			calls = calls + tally
		}

	}
	buffer.WriteString("# HELP figura_response_total_count Total Number of Requests.\n")
	buffer.WriteString("# TYPE figura_response_total_count counter\n")
	buffer.WriteString(fmt.Sprintf("figura_response_total_count %v\n", calls))
	buffer.WriteString("\n")

	return buffer.String()
}
