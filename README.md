# wave

[![GoDoc](https://godoc.org/github.com/bake/wave?status.svg)](https://godoc.org/github.com/bake/wave)
[![Go Report Card](https://goreportcard.com/badge/github.com/bake/wave)](https://goreportcard.com/report/github.com/bake/wave)
[![codecov](https://codecov.io/gh/bake/wave/branch/master/graph/badge.svg)](https://codecov.io/gh/bake/wave)

Package `wave` offers a simplified API to read and write WAVE files by only
allowing to work with samples directly. The underlying package `riff` contains
implementations for a reader and a writer for RIFF files.

## Reader example

Create a new WAVE reader by wrapping it around an `io.Reader`, optionally a
buffered one.

```go
r, err := os.Open("audio.wav")
if err != nil {
  log.Fatalf("could not open file: %v", err)
}
buf := bufio.NewReader(r)
wavr, err := wave.NewReader(buf)
if err != nil {
  log.Fatalf("could not create wave reader: %v", err)
}
```

Read the samples one by one or into a slice of integers. The `wave.Reader` skips
all non-data chunks.

```go
for {
  sample, err := wavr.Sample()
  if err == io.EOF {
    break
  }
  if err != nil {
    log.Fatalf("could not read sample: %v", err)
  }
  fmt.Println(sample)
}
```

```go
samples, err := wavr.Samples()
if err != nil {
  log.Fatalf("could not read samples: %v", err)
}
```
