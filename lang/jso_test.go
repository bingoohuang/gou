package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint gochecknoglobals
var awesomeJSON = []byte(`{
  "id": "123456789",
  "message": "Total awesomeness",
  "score": 9.99,
  "confirmed": true
}`)

func TestAwesomeToJSON(t *testing.T) {
	awesome := Awesome{"123456789", "Total awesomeness", 9.99, true}

	testJSON, err := JSONMarshalIndent(awesome, "", "  ")

	assert.Nil(t, err)
	assert.Equal(t, testJSON, awesomeJSON)
}

func TestAwesomeFromJSON(t *testing.T) {
	var awesome Awesome

	assert.Nil(t, JSONUnmarshal(awesomeJSON, &awesome))
	assert.Equal(t, Awesome{"123456789", "Total awesomeness", 9.99, true}, awesome)
}

type Awesome struct {
	ID        string
	Message   string
	Score     float64
	Confirmed bool
}
