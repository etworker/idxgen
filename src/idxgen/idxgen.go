// Intro:
// generate uid in 8 char, e.g. Sn4YhXzU
// each char is one of "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"
// so totally the space is 64^8 = 2^48 = 281474976710656
//
// Note:
// user should override IsUniqueIdx() for unique verification
//
// Example:
// please check idxgen_test.go for example

package idxgen

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

type IIdxGen interface {
	GenIdx() string
	IsUniqueIdx(idx string) bool
	Stop()
}

type IdxGen struct {
	chIdx    chan string
	chIdxInd chan uint
}

func NewIdxGen() *IdxGen {
	ig := &IdxGen{}

	ig.chIdx = make(chan string, 32)
	ig.chIdxInd = make(chan uint)
	go ig.start()

	return ig
}

func (ig *IdxGen) GenIdx() string {
	idx := <-ig.chIdx

	for {
		if ig.IsUniqueIdx(idx) {
			break
		}

		idx = <-ig.chIdx
	}

	return idx
}

func (ig *IdxGen) IsUniqueIdx(idx string) bool {
	return true
}

func (ig *IdxGen) start() {
	m := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-")

Loop:
	for {
		// Generate new idx
		u4, err := uuid.NewV4()
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		// Combine 2 bytes -> m[x], for 8 times
		idx := ""
		for i := 0; i < 16; i += 2 {
			u := u4[i]<<4 + u4[i+1]
			j := u % 64
			idx += string(m[j])
		}

		select {
		case ig.chIdx <- idx:
		case _ = <-ig.chIdxInd:
			fmt.Println("close chIdx")
			close(ig.chIdx)
			break Loop
		}
	}
}

func (ig *IdxGen) Stop() {
	ig.chIdxInd <- 0
}
