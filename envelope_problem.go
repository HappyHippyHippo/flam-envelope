package rest

import (
	"fmt"
	"net/http"

	flam "github.com/happyhippyhippo/flam"
)

type ProblemEnvelope struct {
	flam.Bag
}

func NewProblemEnvelope(
	status int,
	typeURI,
	title,
	details,
	instance string,
	ctx ...flam.Bag,
) ProblemEnvelope {
	problem := ProblemEnvelope{Bag: flam.Bag{}}
	for _, c := range ctx {
		problem.Merge(c)
	}
	_ = problem.Bag.Set("status", status)
	_ = problem.Bag.Set("type", typeURI)
	_ = problem.Bag.Set("title", title)
	_ = problem.Bag.Set("detail", details)
	_ = problem.Bag.Set("instance", instance)
	return problem
}

func NewProblemEnvelopeFrom(
	envelope Envelope,
) ProblemEnvelope {
	problem := ProblemEnvelope{Bag: flam.Bag{}}
	_ = problem.Bag.Set("status", envelope.Status.Status)
	_ = problem.Bag.Set("message", "unknown error")

	if len(envelope.Status.Errors) > 0 {
		e := envelope.Status.Errors[0]
		problem.Merge(e.Context)

		_ = problem.Bag.Set("status", e.Status)
		_ = problem.Bag.Set("id", e.Id)
	}

	status := problem.Bag.Int("status")
	_ = problem.Bag.Set("type", fmt.Sprintf(
		"https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/%d",
		status))
	_ = problem.Bag.Set("title", http.StatusText(status))
	_ = problem.Bag.Set("detail", problem.Bag.String("message"))

	return problem
}

func (problem ProblemEnvelope) GetStatus() int {
	return problem.Bag.Int("status", 0)
}

func (problem ProblemEnvelope) GetId() string {
	return problem.Bag.String("id", "")
}

func (problem ProblemEnvelope) GetType() string {
	return problem.Bag.String("type", "")
}

func (problem ProblemEnvelope) GetTitle() string {
	return problem.Bag.String("title", "")
}

func (problem ProblemEnvelope) GetDetail() string {
	return problem.Bag.String("detail", "")
}

func (problem ProblemEnvelope) GetInstance() string {
	return problem.Bag.String("instance", "")
}

func (problem ProblemEnvelope) Get(
	key string,
) any {
	return problem.Bag.Get(key)
}

func (problem ProblemEnvelope) With(
	key string,
	value any,
) ProblemEnvelope {
	cloned := ProblemEnvelope{Bag: problem.Bag.Clone()}
	_ = cloned.Bag.Set(key, value)
	return cloned
}
