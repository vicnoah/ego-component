package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego-component/ek8s"
	"github.com/gotomicro/ego/core/elog"
)

func main() {
	if err := ego.New().Invoker(
		invokerGrpc,
	).Run(); err != nil {
		elog.Error("startup", elog.FieldErr(err))
	}
}

func invokerGrpc() error {
	obj := ek8s.Load("k8s").Build()
	list, err := obj.ListPod("svc-oss")
	if err != nil {
		panic(err)
	}

	spew.Dump(list)
	return nil
}
