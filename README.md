# idxgen

## Introduction

If you want a unique id, the traditional way is to generate a GUID such as `43944275-7a4c-4bc6-4b47-5f2dc9e2f97e`.

GUID is good, but it cost 36 chars to save one GUID.

idx is one way to convert `43944275-7a4c-4bc6-4b47-5f2dc9e2f97e` to `4lISTtOe` with some loss.

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

## Example

I have deployed this service at [here](http://idxgen-etworker.alaudacn.me:49662/), every time you refresh it, you will get an idx generated, please check it.