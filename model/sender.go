package model

import "github.com/pkg/errors"

var (
	ErrConsensusSenderPropagate = errors.Errorf("Failed ConsensusSender Propagate")
	ErrConsensusSenderPropose   = errors.Errorf("Failed ConsensusSender Propose")
	ErrConsensusSenderVote      = errors.Errorf("Failed ConsensusSender Vote")
	ErrConsensusSenderPreCommit = errors.Errorf("Failed ConsensusSender PreCommit")
)

type ConsensusSender interface {
	Propagate(tx Transaction) error
	Propose(proposal Proposal) error
	Vote(vote VoteMessage) error
	PreCommit(vote VoteMessage) error
}
