// Speed:
// generate 10,000,000 idx in about 30 seconds with very low duplicate rate

package idxgen

import (
	"fmt"
	"testing"
)

func TestIdxgen(t *testing.T) {
	ig := NewIdxGen()
	n := 10000000
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
