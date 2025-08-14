package rest

import (
	"net/http"

	flam "github.com/happyhippyhippo/flam"
)

type Envelope struct {
	Status     EnvelopeStatus `json:"status,omitempty"`
	Data       any            `json:"data,omitempty"`
	Pagination any            `json:"pagination,omitempty"`
}

func NewEnvelope(
	data ...any,
) Envelope {
	env := Envelope{}
	env.Status.Success = true
	env.Status.Status = http.StatusOK
	if len(data) > 0 {
		env.Data = data[0]
	}

	return env
}

func (envelope Envelope) WithNewError(
	status int,
	e error,
	ctx ...flam.Bag,
) Envelope {
	return envelope.WithError(NewEnvelopeError(status, e, ctx...))
}

func (envelope Envelope) WithError(
	e EnvelopeError,
) Envelope {
	if envelope.Status.Success {
		envelope.Status.Success = false
		envelope.Status.Status = e.Status
	}

	envelope.Status.Errors = append(envelope.Status.Errors, e)

	return envelope
}

func (envelope Envelope) WithServiceId(
	val int,
) Envelope {
	for id, e := range envelope.Status.Errors {
		envelope.Status.Errors[id] = e.WithServiceId(val)
	}

	return envelope
}

func (envelope Envelope) WithEndpointId(
	val int,
) Envelope {
	for id, e := range envelope.Status.Errors {
		envelope.Status.Errors[id] = e.WithEndpointId(val)
	}

	return envelope
}

func (envelope Envelope) WithNewEnvelopePagination(
	search string,
	start,
	count,
	total int64,
) Envelope {
	envelope.Pagination = NewEnvelopePagination(search, start, count, total)

	return envelope
}

func (envelope Envelope) WithPagination(
	pagination any,
) Envelope {
	envelope.Pagination = pagination

	return envelope
}
