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
}

func NewScraper() *Scraper {
  cmd := exec.Command(command, args...)
  return &Scraper{"uri", cmd}
}

func (s *Scraper) Start(wg *sync.WaitGroup) {
  stdoutDone := s.PipeStdout()
  stderrDone := s.PipeStderr()

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
    <-stdoutDone
    <-stderrDone
    s.Kill(wg)
  }()

  go func() {
    <-signalStop
    s.Kill(wg)
  }()
}

func (s *Scraper) Kill(wg *sync.WaitGroup) {
  s.cmd.Process.Kill()
  wg.Done()
}

func (s *Scraper) PipeStdout() <-chan bool {
  pipe, pipeErr := s.cmd.StdoutPipe()

  if pipeErr != nil {
    panic(pipeErr)
  }

  reader := bufio.NewReader(pipe)

  done := make(chan bool)

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

    done <- true
  }()

  return done
}

func (s *Scraper) PipeStderr() <-chan bool {
  pipe, pipeErr := s.cmd.StderrPipe()

  if pipeErr != nil {
    panic(pipeErr)
  }

  reader := bufio.NewReader(pipe)

  done := make(chan bool)

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

    done <- true
  }()

  return done
}

func main() {
  wg := &sync.WaitGroup{}

  s1 := NewScraper()
  s2 := NewScraper()

  s1.Start(wg)
  s2.Start(wg)

  wg.Wait()
}
