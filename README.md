# wave

[![GoDoc](https://godoc.org/github.com/bake/wave?status.svg)](https://godoc.org/github.com/bake/wave)
[![Go Report Card](https://goreportcard.com/badge/github.com/bake/wave)](https://goreportcard.com/report/github.com/bake/wave)
[![codecov](https://codecov.io/gh/bake/wave/branch/master/graph/badge.svg)](https://codecov.io/gh/bake/wave)

Package [`wave`](https://godoc.org/github.com/bake/wave) offers a simplified API
to read and write WAVE files by only allowing to work with samples directly. The
underlying package [`riff`](https://godoc.org/github.com/bake/wave/riff)
contains implementations for a reader and a writer for RIFF files.

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

## Writer example

Create a new WAVE writer by wrapping it around an `io.WriteSeeker`. This one is
automatically buffered.

```go
format := wave.Format{
  AudioFormat:   1,
  NumChans:      2,
  SampleRate:    44100,
  ByteRate:      176400,
  BlockAlign:    4,
  BitsPerSample: 16,
}
samples := []int{
  0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
  -6348, -23005, -3524, -3548, -12783, 3354,
  0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
}
w, err := os.Create("audio.wav")
if err != nil {
  log.Fatalf("could not create file: %v", err)
}
defer w.Close()
wavw, err := wave.NewWriter(w, format)
if err != nil {
  log.Fatalf("could not create wave writer: %v", err)
}
defer wavw.Close()
```

Write the samples one by one or as a slice of integers.

```go
for _, s := range samples {
  if err := wavw.Sample(s); err != nil {
    log.Fatalf("could not write sample %d: %v", s, err)
  }
}
```

```go
if err := wavw.Samples(samples); err != nil {
  log.Fatalf("could not write samples: %v", err)
}
```

Before creating a new chunk, the current one has to be closed which
automatically writes its size.
