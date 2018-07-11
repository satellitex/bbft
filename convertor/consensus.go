package convertor

import "github.com/satellitex/bbft/proto"

type ProposalTx struct {
	*bbft.ProposalTx
}

type VoteMessage struct {
	*bbft.VoteMessage
}