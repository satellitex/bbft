package main

import (
	"encoding/base64"
	"fmt"
	"github.com/satellitex/bbft/config"
	"github.com/satellitex/bbft/convertor"
	"github.com/satellitex/bbft/proto"
	. "github.com/satellitex/bbft/test_utils"
	"github.com/satellitex/bbft/usecase"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

func NewTxGateClient(conf *config.BBFTConfig) bbft.TxGateClient {
	conn, err := grpc.Dial(conf.Host+":"+conf.Port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return bbft.NewTxGateClient(conn)
}

func main() {
	conf := GetTestConfig()

	fmt.Println(base64.StdEncoding.EncodeToString(conf.PublicKey))
	fmt.Println(base64.StdEncoding.EncodeToString(conf.SecretKey))

	client := NewTxGateClient(conf)
	rand.Seed(usecase.Now())

	for i := 0; ; i++ {
		tx, err := convertor.NewTxModelBuilder().Message(fmt.Sprintf(RandomStr()+"Messageid: %d", i)).Sign(conf.PublicKey, conf.SecretKey).Build()
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = client.Write(context.TODO(), tx.(*convertor.Transaction).Transaction)
		if err != nil {
			fmt.Println("failed!  ", i, err)
		} else {
			fmt.Println("success! ", i)
		}
		time.Sleep(time.Millisecond)
	}
}
