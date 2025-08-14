package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	flam "github.com/happyhippyhippo/flam"
	envelope "github.com/happyhippyhippo/flam-envelope"
)

func Test_NewEnvelopeError(t *testing.T) {
	e := errors.New("test error")
	ctx := flam.Bag{"key": "value"}

	envelopeErr := envelope.NewEnvelopeError(http.StatusBadRequest, e, ctx)
	assert.Equal(t, http.StatusBadRequest, envelopeErr.Status)
	assert.Equal(t, "test error", envelopeErr.Context["message"])
	assert.Equal(t, "value", envelopeErr.Context["key"])
}

func Test_EnvelopeError_Composition(t *testing.T) {
	baseErr := envelope.NewEnvelopeError(http.StatusInternalServerError, errors.New("base"))

	testCases := []struct {
		name     string
		composer func(envelope.EnvelopeError) envelope.EnvelopeError
		expected string
	}{
		{"WithStatus", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithStatus(http.StatusNotFound) }, ""},
		{"WithServiceId", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithServiceId(1) }, "s:1"},
		{"WithEndpointId", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithEndpointId(2) }, "e:2"},
		{"WithParamId", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithParamId(3) }, "p:3"},
		{"WithErrorId (string)", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithErrorId("custom") }, "custom"},
		{"WithErrorId (int)", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithErrorId(123) }, "c:123"},
		{"WithContext", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithContext("user", "test") }, ""},
		{"Composed", func(e envelope.EnvelopeError) envelope.EnvelopeError {
			return e.WithServiceId(1).WithEndpointId(2).WithParamId(3).WithErrorId(4)
		}, "s:1.e:2.p:3.c:4"},
		{"Composed with gaps", func(e envelope.EnvelopeError) envelope.EnvelopeError { return e.WithServiceId(1).WithParamId(3) }, "s:1.p:3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updated := tc.composer(baseErr)

			switch {
			case tc.expected != "":
				assert.Equal(t, tc.expected, updated.Id)
			case tc.name == "WithStatus":
				assert.Equal(t, http.StatusNotFound, updated.Status)
			case tc.name == "WithContext":
				assert.Equal(t, "test", updated.Context["user"])
			}
		})
	}
}
