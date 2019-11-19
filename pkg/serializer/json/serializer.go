package json

import (
	"encoding/json"
	"github.com/icowan/shorter/pkg/service"
	"github.com/pkg/errors"
)

type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*service.Redirect, error) {
	redirect := &service.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *service.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
