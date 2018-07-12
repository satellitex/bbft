package controller

import (
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase"
	"golang.org/x/net/context"
)

type ConsensusController struct {
	receiver usecase.ConsensusReceiver
}

func (c *ConsensusController) Propagate(_ context.Context, ptx *bbft.ProposalTx) (*bbft.ConsensusResponse, error) {
	proposalTx := &convertor.ProposalTx{ptx}
	err := c.receiver.Propagate(proposalTx)
	if err != nil {
		return nil, err
	}
	return &bbft.ConsensusResponse{}, nil
}

func (c *ConsensusController) Propose(_ context.Context, p *bbft.Proposal) (*bbft.ConsensusResponse, error) {
	proposal := &convertor.Proposal{p}
	err := c.receiver.Propose(proposal)
	if err != nil {
		return nil, err
	}
	return &bbft.ConsensusResponse{}, nil
}
func (c *ConsensusController) Vote(_ context.Context, v *bbft.VoteMessage) (*bbft.ConsensusResponse, error) {
	vote := &convertor.VoteMessage{v}
	err := c.receiver.Vote(vote)
	if err != nil {
		return nil, err
	}
	return &bbft.ConsensusResponse{}, nil
}
func (c *ConsensusController) PreCommit(_ context.Context, v *bbft.VoteMessage) (*bbft.ConsensusResponse, error) {
	preCommit := &convertor.VoteMessage{v}
	err := c.receiver.PreCommit(preCommit)
	if err != nil {
		return nil, err
	}
	return &bbft.ConsensusResponse{}, nil
}
