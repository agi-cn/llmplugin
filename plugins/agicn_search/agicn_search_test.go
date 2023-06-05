package agicn_search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgicnSearch(t *testing.T) {

	// TODO(zy): fix agi.cn search
	t.Skip("agi.cn search not valid NOW")

	ts := []struct {
		testname string
		query    string
	}{
		{
			"Search in english",
			"NBA schedule today",
		},
		{
			"Search in chinese",
			"今天nba有哪些比赛",
		},
	}

	s := NewAgicnSearch()

	for _, tc := range ts {

		t.Run(tc.testname, func(t *testing.T) {

			answer, err := s.Do(context.Background(), tc.query)
			assert.NoError(t, err)

			assert.NotEmpty(t, answer)

			t.Logf("query=%s\nresult:%s", tc.query, answer)
		})
	}
}
