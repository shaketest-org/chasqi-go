package result

import (
	"chasqi-go/types"
)

func GetTotalErrorCount(results []*types.AgentResult) (errorCount, okCount int, totalDurationInSeconds float64) {
	for _, result := range results {
		errorCount += result.ErrorCount
		okCount += result.SuccessCount
		totalDurationInSeconds += result.Duration().Seconds()
	}
	return errorCount, okCount, totalDurationInSeconds
}
