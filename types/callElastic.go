package types

import (
	"net/rpc"
)

type CallElastic interface {
	FuncElastic() (*rpc.Client, error)
}

type Elastic struct{}

func (elastic *Elastic) FuncElastic() (*rpc.Client, error) {

	client, err := rpc.DialHTTP("tcp", "0.0.0.0:1234")
	if err != nil {
		return client, err
	}

	return client, nil

}
