package job

import "github.com/google/uuid"

type ResultType int

const (
	ResultTypeNew ResultType = iota
	ResultTypeUpdated
	ResultTypeNoChange
	ResultTypeMissing
	ResultTypeFailed
)

type Result struct {
	error string
	t     ResultType
	jobID uuid.UUID
}

type ResultOptional func(*Result)

func WithError(err string) ResultOptional {
	return func(r *Result) {
		r.error = err
	}
}

func NewResult(jobID uuid.UUID, t ResultType, opts ...ResultOptional) *Result {
	r := &Result{
		jobID: jobID,
		t:     t,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Result) JobID() uuid.UUID {
	return r.jobID
}

func (r *Result) Type() ResultType {
	return r.t
}
