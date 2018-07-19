package controller

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
	"github.com/satellitex/bbft/usecase"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsensusController struct {
	receiver usecase.ConsensusReceiver
	author   *convertor.Author
}

func NewConsensusController(receiver usecase.ConsensusReceiver, author *convertor.Author) *ConsensusController {
	return &ConsensusController{
		receiver: receiver,
		author:   author,
	}
}

func (c *ConsensusController) Propagate(ctx context.Context, tx *bbft.Transaction) (*bbft.ConsensusResponse, error) {
	ctx, err := c.author.ProtoAurhorize(ctx, tx)
	if err != nil { // Unauthenticated ( code = 16 )
		return nil, err
	}

	proposalTx := &convertor.Transaction{tx}
	err = c.receiver.Propagate(proposalTx)
	if err != nil {
		cause := errors.Cause(err)
		if cause == model.ErrStatelessTxValidate {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if cause == usecase.ErrAlradyReceivedSameObject {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
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
