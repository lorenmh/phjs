package main

import (
  "bufio"
  "sync"
  "os"
  "os/exec"
  "os/signal"
  "io"
  "syscall"
  "fmt"
)

const command = "node"
var args = []string{"a.js"}

var stdoutMutex = new(sync.Mutex)
var stderrMutex = new(sync.Mutex)

var _ = fmt.Printf

type Scraper struct {
  uri string
  cmd *exec.Cmd
  done chan bool
}

func NewScraper() *Scraper {
  cmd := exec.Command(command, args...)
  return &Scraper{"uri", cmd, make(chan bool)}
}

func (s *Scraper) Start(wg *sync.WaitGroup) {
  s.PipeStdout()
  s.PipeStderr()
  s.cmd.Start()

  wg.Add(1)

  signalStop := make(chan os.Signal, 1)

  signal.Notify(
    signalStop,
    syscall.SIGHUP,
    syscall.SIGINT,
    syscall.SIGTERM,
    syscall.SIGQUIT,
  )

  go func() {
    <-s.done
    <-s.done
    s.Kill(wg)
  }()

  go func() {
    <-signalStop
    fmt.Println("killed")
    s.Kill(wg)
  }()
}

func (s *Scraper) Kill(wg *sync.WaitGroup) {
  s.cmd.Process.Kill()
  wg.Done()
}

func (s *Scraper) PipeStdout() {
  pipe, err := s.cmd.StdoutPipe()
  reader := bufio.NewReader(pipe)

  if err != nil {
    panic(err)
  }

  go func() {
    for {
      line, readErr := reader.ReadBytes('\n');

      if readErr == io.EOF {
        break
      }

      if readErr != nil {
        panic(readErr)
      }

      stdoutMutex.Lock()

      _, writeErr := os.Stdout.Write(line)

      if writeErr != nil {
        panic(writeErr)
      }

      stdoutMutex.Unlock()
    }

    s.done <- true
  }()
}

func (s *Scraper) PipeStderr() {
  pipe, err := s.cmd.StderrPipe()
  reader := bufio.NewReader(pipe)

  if err != nil {
    panic(err)
  }

  go func() {
    for {
      line, readErr := reader.ReadBytes('\n');

      if readErr == io.EOF {
        break
      }

      if readErr != nil {
        panic(readErr)
      }

      stderrMutex.Lock()

      _, writeErr := os.Stderr.Write(line)

      if writeErr != nil {
        panic(writeErr)
      }

      stderrMutex.Unlock()
    }

    s.done <- true
  }()
}

func main() {
  wg := &sync.WaitGroup{}

  s1 := NewScraper()
  s2 := NewScraper()

  s1.Start(wg)
  s2.Start(wg)

  wg.Wait()
}
