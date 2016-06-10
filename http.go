package piroxiki

import "fmt"

type Http struct {
	Method     string
	URL        string
	PORT       string
	Body       string
	BodyFormat string
}

func (*Http) HandleInput(policy Policy, filter Filter, errc chan<- error) (<-chan Message, error) {
	return nil, nil
}

func (*Http) HandleOutput(policy Policy, message <-chan Message, errc chan<- error) error {
	for {
		msg := <-message
		fmt.Println(msg)
	}
	return nil
}
