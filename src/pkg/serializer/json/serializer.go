/**
 * @Time : 2019-11-15 14:27
 * @Author : solacowa@gmail.com
 * @File : serializer
 * @Software: GoLand
 */

package json

import (
	"encoding/json"
	"github.com/icowan/shorter/src/pkg/shortener"
	"github.com/pkg/errors"
)

type Redirect struct {
}

func (r *Redirect) Decode(input []byte) (redirect *shortener.Redirect, err error) {
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return
}

func (r *Redirect) Encode(redirect *shortener.Redirect) (b []byte, err error) {
	b, err = json.Marshal(redirect)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return
}
