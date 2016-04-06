package main

import (
  "time"
  "math/rand"
  "sync"
  "os"
  "io"
  "fmt"
  "strings"
)

const bufferSize = 4096

var wg = sync.WaitGroup{}

var mutex = &sync.Mutex{}

// main question is regarding this 'WriteStdout' function; is this done in the
// correct way? How should I improve this?
// WriteStdout takes a io.Reader, locks until the reader has been read from,
// then unlocks.
func WriteStdout(r io.Reader) {
  // lock the mutex
  mutex.Lock()
  // defer unlocking the mutex
  defer mutex.Unlock()

  // create a new buffer of size bufferSize (4096)
  buffer := make([]byte, bufferSize)

  // read all of the bytes from the reader and output
  for {
    // read from the reader into the buffer
    bytesRead, readErr := r.Read(buffer)

    // if bytesRead is 0 then we are done
    if bytesRead == 0 {
      break
    }

    // write to Stdout
    _, writeErr := os.Stdout.Write(buffer)

    // check for a read error AFTER writing (go docs said to do this)
    if readErr != nil {
      panic(readErr)
    }

    // check for a write error
    if writeErr != nil {
      panic(writeErr)
    }
  }
}

func output(index int) {
  // to add some randomness to the order of the operations
  time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
  // write to stdout
  WriteStdout(strings.NewReader(fmt.Sprintf("Hello from %d\n", index)))
  // decrement the waitgroup
  wg.Done()
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  for i := 0; i < 20; i++ {
    // create a new go routine to output
    go output(i)

    // increment waitgroup
    wg.Add(1)
  }

  // block until waitgroups is complete
  wg.Wait()
}
