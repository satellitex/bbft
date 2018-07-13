package model

import "github.com/pkg/errors"

var (
	ErrStatefulValidate  = errors.Errorf("Failed StatefulValidator Validate")
	ErrStatelessValidate = errors.Errorf("Failed StatelessValidator Validate")
)

type StatefulValidator interface {
	Validate(block Block) error
}

type StatelessValidator interface {
	Validate(block Block) error
}
