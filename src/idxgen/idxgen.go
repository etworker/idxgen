// Intro:
// generate uid in 8 char, e.g. Sn4YhXzU
// each char is one of "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"
// so totally the space is 64^8 = 2^48 = 281474976710656
//
// There are 4 types IdxGen for use
// BaseIdxGen, no check for idx unique
// MapIdxGen, check idx unique with map, append new idx to file when generated
// Map2IdxGen, check idx unique with map, save all new idx to file when stop
// BFIdxGen, check idx unique with Bloom Filter, fast to check while slow to load and save
//
// Example:
//
// ig := idxgen.NewIdxGen(&idxgen.MapIdxGen{})
// idx := idxgen.GenIdx(ig)
// idxgen.Stop(ig)
//
// please check idxgen_test.go for more example

package idxgen

import (
	"bufio"
	"fmt"
	"github.com/AndreasBriese/bbloom"
	"github.com/nu7hatch/gouuid"
	"io"
	"io/ioutil"
	"os"
)

type IIdxGen interface { // Index Generator Interface
	SetName(name string)         // Set name
	GetName() string             // Return name
	Init()                       // Init
	GenIdx() string              // Generate unique index
	AppendIdx(idx string)        // Append index in storage
	LoadIdx() uint64             // Load index from storage
	SaveIdx()                    // Save index to storage
	IsUniqueIdx(idx string) bool // Is index unique
	Start()                      // Start
	Stop()                       // Stop
}

// ------------------------------------------------------------
// Public function
// ------------------------------------------------------------

func NewIdxGen(iig IIdxGen) IIdxGen {
	iig.Init()
	go Start(iig)
	return iig
}

func GenIdx(iig IIdxGen) string {
	idx := iig.GenIdx()
	iig.AppendIdx(idx)
	return idx
}

func Start(iig IIdxGen) {
	n := iig.LoadIdx()
	if n > 0 {
		fmt.Println("LoadIdx =", n)
	}
	iig.Start()
}

func Stop(iig IIdxGen) {
	iig.Stop()
	iig.SaveIdx()
}

// ------------------------------------------------------------
// IdxGen
// Basic IdxGen, not check unique
// ------------------------------------------------------------

type BaseIdxGen struct {
	name     string
	chIdx    chan string
	chIdxInd chan uint
}

func (ig *BaseIdxGen) GetName() string {
	return ig.name
}

func (ig *BaseIdxGen) SetName(name string) {
	ig.name = name
}

func (ig *BaseIdxGen) Init() {
	ig.name = "BaseIdxGen"
	ig.chIdx = make(chan string, 32)
	ig.chIdxInd = make(chan uint)

	fmt.Println(ig.name, "Init")
}

func (ig *BaseIdxGen) Start() {
	fmt.Println(ig.name, "Start")
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

func (ig *BaseIdxGen) GenIdx() string {
	// fmt.Println(ig.name, "GenIdx")
	idx := <-ig.chIdx

	for {
		if ig.IsUniqueIdx(idx) {
			break
		} else {
			fmt.Println(idx, "!!!") // duplicate!
		}

		idx = <-ig.chIdx
	}

	return idx
}

func (ig *BaseIdxGen) IsUniqueIdx(idx string) bool {
	return true
}

func (ig *BaseIdxGen) AppendIdx(idx string) {
	// fmt.Println(ig.name, "AppendIdx:", idx)
}

func (ig *BaseIdxGen) LoadIdx() uint64 {
	fmt.Println(ig.name, "LoadIdx")
	return 0
}

func (ig *BaseIdxGen) SaveIdx() {
	fmt.Println(ig.name, "SaveIdx")
}

func (ig *BaseIdxGen) Stop() {
	fmt.Println(ig.name, "Stop")
	ig.chIdxInd <- 0
}

// ------------------------------------------------------------
// MapIdxGen
// Using Map
// ------------------------------------------------------------

type MapIdxGen struct {
	BaseIdxGen

	m           map[string]int
	idxFilename string
}

func (ig *MapIdxGen) Init() {
	ig.BaseIdxGen.Init()

	ig.name = "MapIdxGen"
	ig.idxFilename = "idx.txt"
	ig.m = make(map[string]int)

	fmt.Println(ig.name, "Init")
}

func (ig *MapIdxGen) IsUniqueIdx(idx string) bool {
	_, prs := ig.m[idx]
	return !prs
}

func (ig *MapIdxGen) AppendIdx(idx string) {
	// fmt.Println(ig.name, "AppendIdx:", idx)

	file, err := os.OpenFile(ig.idxFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(idx)
	file.WriteString("\n")
}

func (ig *MapIdxGen) LoadIdx() uint64 {
	fmt.Println(ig.name, "LoadIdx")

	// Check if idx file exist
	if _, err := os.Stat(ig.idxFilename); os.IsNotExist(err) {
		fmt.Println("idx file not exist")
		return 0
	}

	file, err := os.Open(ig.idxFilename)
	if err != nil {
		return 0
	}
	defer file.Close()

	var n uint64
	n = 0
	bfRd := bufio.NewReader(file)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return n
			}
			return n
		}
		ig.m[string(line)] = 0
		n++
	}
}

func (ig *MapIdxGen) SaveIdx() {
	fmt.Println(ig.name, "SaveIdx")
}

// ------------------------------------------------------------
// Map2IdxGen
// Using Map
// ------------------------------------------------------------

type Map2IdxGen struct {
	MapIdxGen

	m2 map[string]int
}

func (ig *Map2IdxGen) Init() {
	ig.MapIdxGen.Init()

	ig.name = "Map2IdxGen"
	ig.m2 = make(map[string]int)

	fmt.Println(ig.name, "Init")
}

func (ig *Map2IdxGen) IsUniqueIdx(idx string) bool {
	_, prs := ig.m[idx]
	if prs {
		_, prs = ig.m2[idx]
	}
	return !prs
}

func (ig *Map2IdxGen) AppendIdx(idx string) {
	// fmt.Println(ig.name, "AppendIdx:", idx)

	ig.m2[idx] = 0
}

func (ig *Map2IdxGen) SaveIdx() {
	fmt.Println(ig.name, "SaveIdx")

	file, err := os.OpenFile(ig.idxFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	for idx := range ig.m2 {
		file.WriteString(idx)
		file.WriteString("\n")
	}
}

// ------------------------------------------------------------
// BFIdxGen
// Using Bloom Filter (2^26, 1%)
// ------------------------------------------------------------

type BFIdxGen struct {
	BaseIdxGen

	bf         bbloom.Bloom
	bfFilename string
}

func (ig *BFIdxGen) Init() {
	ig.BaseIdxGen.Init()

	ig.name = "BFIdxGen"
	ig.bfFilename = "bf.json"
	ig.bf = bbloom.New(float64(1<<26), float64(0.01))

	fmt.Println(ig.name, "Init")
}

func (ig *BFIdxGen) IsUniqueIdx(idx string) bool {
	return !ig.bf.HasTS([]byte(idx))
}

func (ig *BFIdxGen) AppendIdx(idx string) {
	// fmt.Println(ig.name, "AppendIdx:", idx)

	ig.bf.AddTS([]byte(idx))
}

func (ig *BFIdxGen) LoadIdx() uint64 {
	fmt.Println(ig.name, "LoadIdx")

	// Check if idx file exist
	if _, err := os.Stat(ig.bfFilename); os.IsNotExist(err) {
		fmt.Println("idx json file not exist")
		return 0
	}

	file, err := os.Open(ig.bfFilename)
	if err != nil {
		return 0
	}
	defer file.Close()

	Json, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}
	ig.bf = bbloom.JSONUnmarshal(Json)
	return ig.bf.ElemNum
}

func (ig *BFIdxGen) SaveIdx() {
	fmt.Println(ig.name, "SaveIdx")

	ioutil.WriteFile(ig.bfFilename, ig.bf.JSONMarshal(), 0666)
}
