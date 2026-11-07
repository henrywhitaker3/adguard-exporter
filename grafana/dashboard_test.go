package grafana

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardDoesNotTreatRollingStatsAsCounters(t *testing.T) {
	dashboard, err := os.ReadFile("dashboard.json")
	require.NoError(t, err)

	var decoded any
	require.NoError(t, json.Unmarshal(dashboard, &decoded))

	rollingStatCounterFunction := regexp.MustCompile(`(?:rate|increase)\(\s*adguard_(?:queries(?:_blocked)?|replaced_(?:safesearch|safebrowsing|parental)|top_(?:queried_domains|blocked_domains|clients|upstreams))`)
	assert.Empty(t, matchingExpressions(decoded, rollingStatCounterFunction))
}

func matchingExpressions(value any, pattern *regexp.Regexp) []string {
	matches := []string{}
	collectMatchingExpressions(value, pattern, &matches)

	return matches
}

func collectMatchingExpressions(value any, pattern *regexp.Regexp, matches *[]string) {
	switch value := value.(type) {
	case map[string]any:
		for key, item := range value {
			if key == "expr" {
				if expression, ok := item.(string); ok && pattern.MatchString(expression) {
					*matches = append(*matches, expression)
				}
			}
			collectMatchingExpressions(item, pattern, matches)
		}
	case []any:
		for _, item := range value {
			collectMatchingExpressions(item, pattern, matches)
		}
	}
}
