package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type ProposalTx struct {
	*bbft.ProposalTx
}

type VoteMessage struct {
	*bbft.VoteMessage
}

func (p *ProposalTx) GetTransaction() model.Transaction {
	return &Transaction{p.Tx}
}

func (p *ProposalTx) GetSignature() model.Signature {
	return &Signature{p.Signature}
}

func (v *VoteMessage) GetSignature() model.Signature {
	return &Signature{v.Signature}
}

type GrpcConsensusSender struct {
	client bbft.ConsensusGateClient
}

func NewConsensusSender() model.ConsensusSender {
	return &GrpcConsensusSender{nil}
}

func (s *GrpcConsensusSender) Propagate(ptx model.ProposalTx) error {
	return nil
}

func (s *GrpcConsensusSender) Propose(proposal model.Proposal) error {
	return nil
}

func (s *GrpcConsensusSender) Vote(vote model.VoteMessage) error {
	return nil
}

func (s *GrpcConsensusSender) PreCommit(vote model.VoteMessage) error {
	return nil
}
