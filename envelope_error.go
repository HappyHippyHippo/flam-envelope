package rest

import (
	"fmt"
	"strconv"
	"strings"

	flam "github.com/happyhippyhippo/flam"
)

type EnvelopeError struct {
	Status     int      `json:"status"`
	Id         string   `json:"id"`
	Context    flam.Bag `json:"context"`
	serviceId  int
	endpointId int
	paramId    int
	errorId    string
}

func NewEnvelopeError(
	status int,
	e error,
	ctx ...flam.Bag,
) EnvelopeError {
	context := flam.Bag{"message": e.Error()}
	for _, c := range ctx {
		context.Merge(c)
	}

	return EnvelopeError{Status: status, Context: context}.compose()
}

func (envelopeError EnvelopeError) WithStatus(
	val int,
) EnvelopeError {
	envelopeError.Status = val

	return envelopeError
}

func (envelopeError EnvelopeError) WithServiceId(
	val int,
) EnvelopeError {
	envelopeError.serviceId = val

	return envelopeError.compose()
}

func (envelopeError EnvelopeError) WithEndpointId(
	val int,
) EnvelopeError {
	envelopeError.endpointId = val

	return envelopeError.compose()
}

func (envelopeError EnvelopeError) WithParamId(
	param int,
) EnvelopeError {
	envelopeError.paramId = param

	return envelopeError.compose()
}

func (envelopeError EnvelopeError) WithErrorId(
	id any,
) EnvelopeError {
	envelopeError.errorId = fmt.Sprintf("%v", id)

	return envelopeError.compose()
}

func (envelopeError EnvelopeError) WithContext(
	key string,
	value any,
) EnvelopeError {
	envelopeError.Context.Set(key, value)

	return envelopeError
}

func (envelopeError EnvelopeError) compose() EnvelopeError {
	cb := strings.Builder{}
	if envelopeError.serviceId != 0 {
		cb.WriteString(fmt.Sprintf("s:%d", envelopeError.serviceId))
	}

	if envelopeError.endpointId != 0 {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}
		cb.WriteString(fmt.Sprintf("e:%d", envelopeError.endpointId))
	}

	if envelopeError.paramId != 0 {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}
		cb.WriteString(fmt.Sprintf("p:%d", envelopeError.paramId))
	}

	if envelopeError.errorId != "" {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}
		if i, ee := strconv.Atoi(envelopeError.errorId); ee != nil {
			cb.WriteString(envelopeError.errorId)
		} else {
			cb.WriteString(fmt.Sprintf("c:%d", i))
		}
	}

	envelopeError.Id = cb.String()

	return envelopeError
}
