package model

type StatefulValidator interface {
	Validate(block Block) error
}

type StatelessValidator interface {
	Validate(block Block) error
}
