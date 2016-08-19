// Speed:
// generate 10,000,000 idx in about 30 seconds with very low duplicate rate
// load 2202040 idx from file and generate 10,000,000 idx while appending to file
// one by one, totally cost 157 seconds

package idxgen

import (
	"fmt"
	"testing"
	"time"
)

const (
	n = 10000
)

func IdxGenTest(iig IIdxGen, n int) {
	fmt.Println("Begin test:", iig.GetName())

	ig := NewIdxGen(iig)
	t0 := time.Now()
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		k := GenIdx(ig)
		_, prs := m[k]
		if prs {
			fmt.Println("!!!") // duplicate!
		}
		m[k]++
	}
	fmt.Println("Cost time:", time.Now().Sub(t0))
	Stop(ig)

	fmt.Println("End test:", iig.GetName())
}

func TestBaseIdxGen(t *testing.T) {
	IdxGenTest(&BaseIdxGen{}, n)
}

func TestMapIdxGen(t *testing.T) {
	IdxGenTest(&MapIdxGen{}, 10000)
}

func TestMap2IdxGen(t *testing.T) {
	IdxGenTest(&Map2IdxGen{}, 10000)
}

func TestBFIdxGen(t *testing.T) {
	IdxGenTest(&BFIdxGen{}, 10000)
}
