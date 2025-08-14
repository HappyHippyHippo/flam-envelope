package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	envelope "github.com/happyhippyhippo/flam-envelope"
)

func Test_NewEnvelope(t *testing.T) {
	t.Run("WithoutData", func(t *testing.T) {
		env := envelope.NewEnvelope()
		assert.True(t, env.Status.Success)
		assert.Equal(t, http.StatusOK, env.Status.Status)
		assert.Nil(t, env.Data)
	})

	t.Run("WithData", func(t *testing.T) {
		data := "test data"
		env := envelope.NewEnvelope(data)
		assert.True(t, env.Status.Success)
		assert.Equal(t, http.StatusOK, env.Status.Status)
		assert.Equal(t, data, env.Data)
	})
}

func Test_Envelope_WithMethods(t *testing.T) {
	baseErr := errors.New("base error")

	t.Run("WithError", func(t *testing.T) {
		envErr := envelope.NewEnvelopeError(http.StatusBadRequest, baseErr)
		env := envelope.NewEnvelope().WithError(envErr)
		assert.False(t, env.Status.Success)
		assert.Len(t, env.Status.Errors, 1)
		assert.Equal(t, envErr, env.Status.Errors[0])
	})

	t.Run("WithNewError", func(t *testing.T) {
		env := envelope.NewEnvelope().WithNewError(http.StatusBadRequest, baseErr)
		assert.False(t, env.Status.Success)
		require.Len(t, env.Status.Errors, 1)
		assert.Equal(t, http.StatusBadRequest, env.Status.Errors[0].Status)
		assert.Equal(t, baseErr.Error(), env.Status.Errors[0].Context["message"])
	})

	t.Run("WithServiceId and WithEndpointId", func(t *testing.T) {
		env := envelope.NewEnvelope().WithNewError(http.StatusBadRequest, baseErr)
		env = env.WithServiceId(10).WithEndpointId(20)
		require.Len(t, env.Status.Errors, 1)
		assert.Equal(t, "s:10.e:20", env.Status.Errors[0].Id)
	})

	t.Run("WithNewEnvelopePagination", func(t *testing.T) {
		env := envelope.NewEnvelope().WithNewEnvelopePagination("search", 5, 10, 25)
		p, ok := env.Pagination.(envelope.EnvelopePagination)
		require.True(t, ok)
		assert.Equal(t, "search", p.Search)
		assert.Equal(t, int64(5), p.Start)
		assert.Equal(t, int64(10), p.Count)
		assert.Equal(t, int64(25), p.Total)
	})

	t.Run("WithPagination", func(t *testing.T) {
		p := "pagination_data"
		env := envelope.NewEnvelope().WithPagination(p)
		assert.Equal(t, p, env.Pagination)
	})
}
