package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

const LOADING_ALIAS = "..."

func statisticOfDatapoint(datapoint types.Datapoint, statistic types.Statistic) *float64 {
	switch statistic {
	case types.StatisticAverage:
		return datapoint.Average
	case types.StatisticMaximum:
		return datapoint.Maximum
	case types.StatisticMinimum:
		return datapoint.Minimum
	case types.StatisticSampleCount:
		return datapoint.SampleCount
	case types.StatisticSum:
		return datapoint.Sum
	}
	return nil
}

func formatBool(b bool) string {
	if b {
		return "Yes"
	} else {
		return "No"
	}
}

func formatTime(t *time.Time) string {
	return t.Format("02/01/2006 15:04:05")
}

func formatSeconds(seconds int) string {
	m := seconds / 60
	s := seconds % 60
	str := fmt.Sprintf("%dm %ds", m, s)
	return str
}

func formatJson(jsonString string) string {
	buf := new(bytes.Buffer)
	err := json.Indent(buf, []byte(jsonString), "", "    ")
	if err != nil {
		return jsonString
		// log.Fatalf("Unable to format json document %s", jsonString)
	}
	return buf.String()
}

// YYYY-MM-DDThh:mm:ss.sTZD
var ISO_8601 = "2006-01-02T15:04:05.000-0700"
