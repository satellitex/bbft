package model

type StatefulValidator interface {
	Validate(block Block) bool
}

type StatelessValidator interface {
	Validate(block Block) bool
}
