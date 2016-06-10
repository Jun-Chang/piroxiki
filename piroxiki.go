package piroxiki

import (
	"github.com/pkg/errors"
)

type In interface {
	HandleInput(Policy, Filter, chan<- error) (<-chan Message, error)
}

type Out interface {
	HandleOutput(Policy, <-chan Message, chan<- error) error
}

type Policy struct {
	Count           int
	IntervalSeconds int
}

type Message map[string]interface{}

type Filter map[string]interface{}

func NewIn(cnf InConf) (In, error) {
	if mail := cnf.Mail; mail != nil {
		return mail, nil

	}
	if http := cnf.Http; http != nil {
		return http, nil

	}
	return nil, errors.New("Input Handler not exits.")
}

func NewOut(cnf OutConf) (Out, error) {
	if mail := cnf.Mail; mail != nil {
		return mail, nil

	}
	if http := cnf.Http; http != nil {
		return http, nil

	}
	return nil, errors.New("Input Handler not exits.")
}
