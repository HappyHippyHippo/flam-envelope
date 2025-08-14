package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	flam "github.com/happyhippyhippo/flam"
	envelope "github.com/happyhippyhippo/flam-envelope"
)

func Test_NewProblemEnvelope(t *testing.T) {
	ctx := flam.Bag{"extra_field": "extra_value"}
	problem := envelope.NewProblemEnvelope(http.StatusTeapot, "type", "title", "detail", "instance", ctx)
	assert.Equal(t, http.StatusTeapot, problem.GetStatus())
	assert.Equal(t, "type", problem.GetType())
	assert.Equal(t, "title", problem.GetTitle())
	assert.Equal(t, "detail", problem.GetDetail())
	assert.Equal(t, "instance", problem.GetInstance())
	assert.Equal(t, "extra_value", problem.Get("extra_field"))
}

func Test_NewProblemEnvelopeFrom(t *testing.T) {
	t.Run("FromEnvelopeWithErrors", func(t *testing.T) {
		envErr := envelope.NewEnvelopeError(http.StatusConflict, errors.New("conflict"))
		env := envelope.NewEnvelope().WithError(envErr)
		problem := envelope.NewProblemEnvelopeFrom(env)
		assert.Equal(t, http.StatusConflict, problem.GetStatus())
		assert.Equal(t, http.StatusText(http.StatusConflict), problem.GetTitle())
	})

	t.Run("FromEnvelopeWithoutErrors", func(t *testing.T) {
		env := envelope.NewEnvelope()
		env.Status.Status = http.StatusAccepted
		problem := envelope.NewProblemEnvelopeFrom(env)
		assert.Equal(t, http.StatusAccepted, problem.GetStatus())
		assert.Equal(t, "unknown error", problem.GetDetail())
	})
}

func Test_ProblemEnvelope_Getters(t *testing.T) {
	problem := envelope.NewProblemEnvelope(
		http.StatusNotFound,
		"https://example.com/probs/not-found",
		"Not Found",
		"The resource was not found",
		"/resource/123",
	)
	envErr := envelope.NewEnvelopeError(http.StatusNotFound, errors.New("err")).WithServiceId(1)
	assert.NoError(t, problem.Set("id", envErr.Id))

	assert.Equal(t, http.StatusNotFound, problem.GetStatus())
	assert.Equal(t, "s:1", problem.GetId())
	assert.Equal(t, "https://example.com/probs/not-found", problem.GetType())
	assert.Equal(t, "Not Found", problem.GetTitle())
	assert.Equal(t, "The resource was not found", problem.GetDetail())
	assert.Equal(t, "/resource/123", problem.GetInstance())
}

func Test_ProblemEnvelope_With(t *testing.T) {
	problem := envelope.NewProblemEnvelope(200, "", "", "", "")
	updated := problem.With("key", "value")
	assert.Equal(t, "value", updated.Get("key"))
	assert.Nil(t, problem.Get("key"))
}
