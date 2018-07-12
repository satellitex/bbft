package convertor

import (
	"github.com/pkg/errors"
	"github.com/satellitex/bbft/model"
	"github.com/satellitex/bbft/proto"
)

type VoteMessage struct {
	*bbft.VoteMessage
}

func (v *VoteMessage) GetSignature() model.Signature {
	return &Signature{v.Signature}
}

func (v *VoteMessage) Sign(pubKey []byte, privKey []byte) error {
	signature, err := Sign(privKey, v.GetBlockHash())
	if err != nil {
		return err
	}
	if !Verify(pubKey, v.GetBlockHash(), signature) {
		return errors.Errorf("Failed Verify")
	}
	v.Signature = &bbft.Signature{Pubkey: pubKey, Signature: signature}
	return nil
}

func (v *VoteMessage) Verify() bool {
	if v.Signature == nil {
		return false
	}
	return Verify(v.Signature.Pubkey, v.GetBlockHash(), v.Signature.Signature)
}

type GrpcConsensusSender struct {
	client bbft.ConsensusGateClient
}

func NewConsensusSender() model.ConsensusSender {
	return &GrpcConsensusSender{nil}
}

func (s *GrpcConsensusSender) Propagate(tx model.Transaction) error {
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
