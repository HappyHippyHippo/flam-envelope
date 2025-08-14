package rest

type EnvelopeStatus struct {
	Status  int             `json:"status"`
	Success bool            `json:"success"`
	Errors  []EnvelopeError `json:"error"`
}
