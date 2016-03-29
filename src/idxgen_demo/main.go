// Intro:
// idx generate server
//
// Usage:
// run the exe, it will listen on 5678
// if open in browser, e.g. "127.0.0.1:5678", it will get a random idx, e.g. g66bQnpx

package main

import (
	"fmt"
	"idxgen"
	"io"
	"net/http"
)

var (
	ig idxgen.IIdxGen = idxgen.NewIdxGen()
)

func getNewIdx(rw http.ResponseWriter, req *http.Request) {
	idx := ig.GenIdx()
	fmt.Println(idx)
	io.WriteString(rw, idx)
}

func main() {
	http.HandleFunc("/", getNewIdx)
	http.ListenAndServe(":5678", nil) // user could change the port
	ig.Stop()
}
