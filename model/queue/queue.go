package queue

import "github.com/satellitex/bbft/model/proto"

type ProposalTxQueue interface {
	Push(tx proto.ProposalTx)
	Pop() proto.ProposalTx
}