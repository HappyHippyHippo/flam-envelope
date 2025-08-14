package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	envelope "github.com/happyhippyhippo/flam-envelope"
)

func Test_NewEnvelopePagination(t *testing.T) {
	t.Run("FirstPage", func(t *testing.T) {
		p := envelope.NewEnvelopePagination("search", 0, 10, 100)
		assert.Empty(t, p.Prev)
		assert.Equal(t, "?search=search&start=10&count=10", p.Next)
	})

	t.Run("MiddlePage", func(t *testing.T) {
		p := envelope.NewEnvelopePagination("search", 20, 10, 100)
		assert.Equal(t, "?search=search&start=10&count=10", p.Prev)
		assert.Equal(t, "?search=search&start=30&count=10", p.Next)
	})

	t.Run("LastPage", func(t *testing.T) {
		p := envelope.NewEnvelopePagination("search", 90, 10, 100)
		assert.Equal(t, "?search=search&start=80&count=10", p.Prev)
		assert.Empty(t, p.Next)
	})

	t.Run("PrevLinkToStart", func(t *testing.T) {
		p := envelope.NewEnvelopePagination("search", 5, 10, 100)
		assert.Equal(t, "?search=search&start=0&count=10", p.Prev)
	})
}
