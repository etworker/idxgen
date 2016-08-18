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
	n = 1000000
)

func TestBaseIdxGen(t *testing.T) {
	fmt.Println("test BaseIdxGen begin")

	ig := NewIdxGen(&BaseIdxGen{})
	t1 := time.Now()
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		k := GenIdx(ig)
		_, prs := m[k]
		if prs {
			fmt.Println("!!!") // duplicate!
		}
		m[k]++
	}
	fmt.Println("time:", time.Now().Sub(t1))
	Stop(ig)
	fmt.Println("test BaseIdxGen finished")
}

func TestMapIdxGen(t *testing.T) {
	fmt.Println("test MapIdxGen begin")

	ig := NewIdxGen(&MapIdxGen{})
	t1 := time.Now()
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		k := GenIdx(ig)
		_, prs := m[k]
		if prs {
			fmt.Println("!!!") // duplicate!
		}
		m[k]++
	}
	fmt.Println("time:", time.Now().Sub(t1))
	Stop(ig)
	fmt.Println("test MapIdxGen finished")
}

func TestBFIdxGen(t *testing.T) {
	fmt.Println("test BFIdxGen begin")

	ig := NewIdxGen(&BFIdxGen{})
	t1 := time.Now()
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		k := GenIdx(ig)
		_, prs := m[k]
		if prs {
			fmt.Println("!!!") // duplicate!
		}
		m[k]++
	}
	fmt.Println("time:", time.Now().Sub(t1))
	Stop(ig)
	fmt.Println("test BFIdxGen finished")
}
