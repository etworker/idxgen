# idxgen

## Introduction

If you want a unique id, the traditional way is to generate a GUID such as `43944275-7a4c-4bc6-4b47-5f2dc9e2f97e`.

GUID is good, but it cost 36 chars to save one GUID.

idx is one way to map `43944275-7a4c-4bc6-4b47-5f2dc9e2f97e` to `4lISTtOe` with some loss.

idx only cost 8 chars to save, so it is better for some services which need short unique id.

## Algorithm

Ignore the `-`, char space for GUID is only "0123456789abcdef", totally 16.

For idx, char space is "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-", totally 64.

So we map each 2 chars of GUID to 1 char of idx, just like below:

```golang
// Combine 2 bytes -> m[x], for 8 times
idx := ""
for i := 0; i < 16; i += 2 {
  u := u4[i]<<4 + u4[i+1]
  j := u % 64
  idx += string(m[j])
}
```

## Limitation

Of course there is large loss in this map, because we zip the possible space from 16^32 to 64^8, that is 2^128 to 2^48. 

So if you need strict unique id, just try GUID, or implement your own function named `IsUniqueIdx(idx string) bool`.

I have tested keep generate idx for 10,000,000 times, only 1 time duplicated. But I suggest user to implement the `IsUniqueIdx()` function for sure.

I have implemented a simple mechanics for avoiding duplicate id by using a pure text file named idx.txt, each time the service start, it will load the file into a map. When generate a new id, firstly check if it in map, then append it to the file if unique. 

**Note**: This implimentation is very simple but not perfect because it will occupy much more memory when keep so many id in map, so please think about it.

## Usage

### 0. Config environment

```shell
$ ./addenv.sh
```

### 1. Install `gouuid` for generating uuid.

```shell
$ go get github.com/nu7hatch/gouuid
```

### 2. Install `idxgen`

Go to `idxgen` folder and run `go install`

### 3. Test `idxgen`

Go to `idxgen` folder and run `go test`, result will be similar to below:

```
=== RUN   TestIdxgen
LoadIdxFile( idx.txt )
load idx file ok, exist idx count = 5087235
close chIdx
test finished
--- PASS: TestIdxgen (13.82s)
PASS
ok      _/home/ubuntu/work/go/idxgen/src/idxgen 13.846s
```

### 4. Run HTTP service
 
Go to `idxgen_demo` folder and run `go run main.go`, then open browser with the address is `IP`:5678, e.g. `http://localhost:5678`

### 5. Run `idxgen` in docker

Run script as below to build docker image:

```shell
cd src/idxgen_demo
./build.sh
cd ../..
./build_docker.sh
```

Now you can run this docker to provide service:

```
./run_docker.sh
```

## Example

I have deployed this service at [here](http://52.78.112.213:5678), every time you refresh it, you will get an idx generated, please check it.
