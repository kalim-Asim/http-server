package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("Valid single header with extra whitespace", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("     Host:    localhost:42069     \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, 36, n) // includes leading + trailing whitespace + only first crlf
		assert.False(t, done)
	})

	t.Run("Valid 2 headers with existing headers", func(t *testing.T) {
		headers := NewHeaders()

		// first header
		data := []byte("Host: localhost:42069\r\n")
		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.False(t, done)
		assert.Equal(t, "localhost:42069", headers["Host"])
		assert.Equal(t, 23, n)

		// second header (called again with new data)
		data = []byte("User-Agent: curl/8.0\r\n\r\n")
		n, done, err = headers.Parse(data)

		require.NoError(t, err)
		assert.False(t, done)
		assert.Equal(t, "curl/8.0", headers["User-Agent"])
		assert.Equal(t, 22, n)
	})

	t.Run("Valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.True(t, done)
		assert.Equal(t, len(SEPARATOR), n)
		assert.Len(t, headers, 0)
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069       \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})
}
