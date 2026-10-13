package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefPatterns(t *testing.T) {
	// RefPattern: #ZHERO#<target>#{<json meta>}#
	t.Run("RefPattern extracts target and meta", func(t *testing.T) {
		token := `#ZHERO#Article/01H123#{"linkText":"Click me","altText":"image alt"}#`
		matches := RefPattern.FindStringSubmatch(token)

		require := assert.New(t)
		require.Len(matches, 3) // full match, target, meta
		require.Equal(token, matches[0])
		require.Equal("Article/01H123", matches[1])
		require.Equal(`"linkText":"Click me","altText":"image alt"`, matches[2])
	})

	// ShortRefPattern: #refN
	t.Run("ShortRefPattern extracts index", func(t *testing.T) {
		token := "#ref12"
		matches := ShortRefPattern.FindStringSubmatch(token)

		require := assert.New(t)
		require.Len(matches, 2) // full match, index
		require.Equal(token, matches[0])
		require.Equal("12", matches[1])
	})

	t.Run("RefPattern rejects malformed token without target", func(t *testing.T) {
		token := `#ZHERO##{"linkText":"x"}#`
		matches := RefPattern.FindStringSubmatch(token)
		assert.Nil(t, matches)
	})

	t.Run("RefPattern allows empty meta", func(t *testing.T) {
		token := `#ZHERO#Article/01H123#{}#`
		matches := RefPattern.FindStringSubmatch(token)
		require := assert.New(t)
		require.Len(matches, 3)
		require.Equal("Article/01H123", matches[1])
		require.Equal("", matches[2])
	})
}
