package result

import (
	"chasqi-go/types"
)

func GetTotalStats(results []*types.AgentResult) (errorCount, okCount int, totalDurationInSeconds float64) {
	for _, result := range results {
		errorCount += result.ErrorCount
		okCount += result.SuccessCount
		// TODO: this is not correct, gives bad value
		totalDurationInSeconds += result.Duration().Seconds()
	}
	return errorCount, okCount, totalDurationInSeconds / float64(len(results))
}
