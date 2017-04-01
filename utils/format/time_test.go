package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	unixTS   string
	rfcTS    string
	prettyTS string
)

func init() {
	unixTS = "Mon Jan  2 16:33:18 UTC 2017"
	rfcTS = "2017-01-02T16:33:18.511281Z"
	prettyTS = "01-02-2017 4:33:18 PM"
}

func TestGetUnixTimeFromDBTimestamp(t *testing.T) {
	unixTime := GetUnixTimeFromDBTimestamp(rfcTS)

	assert.Equal(t, unixTS, unixTime)
}

func TestGetPrettyTimeFromDBTimestamp(t *testing.T) {
	prettyTime := GetPrettyTimeFromDBTimestamp(rfcTS)

	assert.Equal(t, prettyTS, prettyTime)
}
