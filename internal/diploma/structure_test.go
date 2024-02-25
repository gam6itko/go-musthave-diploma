package diploma

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrderStatusFromString(t *testing.T) {
	t.Run("unknown", func(t *testing.T) {
		s, err := OrderStatusFromString("foo")
		assert.Error(t, err)
		assert.EqualError(t, err, "unknown status: foo")
		assert.Equal(t, StatusUndefined, s)
	})

	t.Run("StatusRegistered", func(t *testing.T) {
		s, err := OrderStatusFromString("REGISTERED")
		assert.NoError(t, err)
		assert.Equal(t, StatusRegistered, s)
	})

	t.Run("StatusInvalid", func(t *testing.T) {
		s, err := OrderStatusFromString("INVALID")
		assert.NoError(t, err)
		assert.Equal(t, StatusInvalid, s)
	})

	t.Run("StatusProcessing", func(t *testing.T) {
		s, err := OrderStatusFromString("PROCESSING")
		assert.NoError(t, err)
		assert.Equal(t, StatusProcessing, s)
	})

	t.Run("StatusProcessed", func(t *testing.T) {
		s, err := OrderStatusFromString("PROCESSED")
		assert.NoError(t, err)
		assert.Equal(t, StatusProcessed, s)
	})
}
