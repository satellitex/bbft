package lock

import "github.com/satellitex/bbft/model/proto"

type Lock interface {
	AddVoteMessage(vote proto.VoteMessage)
	GetLockedProposal() proto.Proposal
}