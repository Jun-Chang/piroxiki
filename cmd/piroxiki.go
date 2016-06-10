package main

import (
	"fmt"

	"github.com/Jun-Chang/piroxiki"
)

func main() {
	cnf, err := piroxiki.LoadConf("./test.toml")
	if err != nil {
		panic(err)
	}
	errc := make(chan error)
	for k, v := range cnf {
		fmt.Println(k)
		in, err := piroxiki.NewIn(v.In)
		if err != nil {
			panic(err)
		}
		message, err := in.HandleInput(v.In.Policy, v.In.Filter, errc)
		if err != nil {
			panic(err)
		}

		out, err := piroxiki.NewOut(v.Out)
		if err != nil {
			panic(err)
		}
		if err := out.HandleOutput(v.Out.Policy, message, errc); err != nil {
			panic(err)
		}
	}

	for {
		err := <-errc
		fmt.Println(err)
	}
}
