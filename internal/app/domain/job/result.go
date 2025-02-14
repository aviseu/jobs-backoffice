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
	id    uuid.UUID
}

type ResultOptional func(*Result)

func WithError(err string) ResultOptional {
	return func(r *Result) {
		r.error = err
	}
}

func NewResult(id uuid.UUID, t ResultType, opts ...ResultOptional) *Result {
	r := &Result{
		id: id,
		t:  t,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Result) JobID() uuid.UUID {
	return r.id
}

func (r *Result) Type() ResultType {
	return r.t
}
