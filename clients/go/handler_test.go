package kubefuncs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	t.Run("NoTopic", func(t *testing.T) {
		_, err := NewEvent("", &Empty{})
		require.Error(t, err)
	})
}
