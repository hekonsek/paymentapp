package payments

import (
	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
	"time"
)

func TestDateFormatSerialization(t *testing.T) {
	// Given
	t.Parallel()
	formattedToday := time.Now().Format("2006-01-02")
	payment := Payment{
		Attributes: Attributes{
			ProcessingDate: ProcessingDate(time.Now()),
		},
	}

	// When
	output, err := json.Marshal(&payment)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, formattedToday, gjson.GetBytes(output, "attributes.processing_date").Str)
}
