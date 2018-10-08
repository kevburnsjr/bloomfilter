package main

import (
	"encoding/base64"
	"fmt"

	"github.com/httpimp/bloomfilter"
)

func main() {
	decoded, err := base64.StdEncoding.DecodeString("iCCACAiAACAACIgAQAIIAqSIEKgowKEKAowIGBIgyomAAIqI")
	if err != nil {
		panic(err)
	}
	bf := bloomfilter.NewFromBytes(decoded, 20)
	fmt.Println(bf.Test([]byte("foo")))
	fmt.Println(bf.Test([]byte("bar")))
	fmt.Println(bf.Test([]byte("baz")))
	fmt.Println(bf.Test([]byte("bork")))
}
