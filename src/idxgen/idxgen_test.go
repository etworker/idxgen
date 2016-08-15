// Speed:
// generate 10,000,000 idx in about 30 seconds with very low duplicate rate
// load 2202040 idx from file and generate 10,000,000 idx while appending to file
// one by one, totally cost 157 seconds

package idxgen

import (
	"fmt"
	"testing"
)

func TestIdxgen(t *testing.T) {
	ig := NewIdxGen()
	n := 1000
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		k := ig.GenIdx()
		//fmt.Println(k)
		_, prs := m[k]
		if prs {
			fmt.Println("!!!") // duplicate!
		}
		m[k]++
	}

	ig.Stop()
	fmt.Println("test finished")
}
