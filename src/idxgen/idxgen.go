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
  "bufio"
  "fmt"
  "github.com/nu7hatch/gouuid"
  "io"
  "os"
)

type IIdxGen interface {
  GenIdx() string
  IsUniqueIdx(idx string) bool
  Stop()
}

type IdxGen struct {
  chIdx       chan string
  chIdxInd    chan uint
  m           map[string]int
  idxFilename string
}

func NewIdxGen() *IdxGen {
  ig := &IdxGen{}

  ig.chIdx = make(chan string, 32)
  ig.chIdxInd = make(chan uint)
  ig.m = make(map[string]int)
  ig.idxFilename = "idx.txt"
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

  err := ig.AppendIdxToFile(ig.idxFilename, idx)
  if err != nil {
    fmt.Println("append idx to file failed.", err.Error())
  }

  return idx
}

func (ig *IdxGen) IsUniqueIdx(idx string) bool {
  // return true

  _, prs := ig.m[idx]
  if prs {
    fmt.Println("!!!") // duplicate!
    return false
  }

  return true
}

func (ig *IdxGen) LoadIdxFile(filename string) (n int, err error) {
  fmt.Println("LoadIdxFile(", filename, ")")

  // Check if idx file exist
  if _, err := os.Stat(filename); os.IsNotExist(err) {
    fmt.Println("idx file not exist")
    return 0, nil
  }

  file, err := os.Open(filename)
  if err != nil {
    return 0, err
  }
  defer file.Close()

  n = 0
  bfRd := bufio.NewReader(file)
  for {
    line, err := bfRd.ReadBytes('\n')
    if err != nil {
      if err == io.EOF {
        return n, nil
      }
      return n, err
    }
    ig.m[string(line)] = 0
    n++
  }
}

func (ig *IdxGen) AppendIdxToFile(filename, idx string) (err error) {
  //fmt.Println("AppendIdxToFile(", filename, idx, ")")

  file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
  if err != nil {
    return err
  }
  defer file.Close()
  file.WriteString(idx)
  file.WriteString("\n")
  return nil
}

func (ig *IdxGen) start() {
  m := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-")

  // Load idx file
  n, err := ig.LoadIdxFile(ig.idxFilename)
  if err != nil {
    fmt.Println("load idx file failed.", err.Error())
  }
  fmt.Println("load idx file ok, exist idx count =", n)

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
