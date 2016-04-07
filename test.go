package main

import (
  "fmt"
  "bufio"
  "sync"
  "os"
  "os/exec"
  "io"
)

//const command = "python"
const command = "printf"
//var args = []string{"-c", "print \"o\" * 100000"}
//var a1 = "-c"
var a1 = "foo\nbar\nbaz\nblah\n"

//var a2 = "\"print \\\"%dA\\\" * 10000; print \\\"\\n\\\"; print \\\"%dB\\\" * 1000\""
//var a2 = "a=1;print a;print a + 1;print a+2"
//var a2 = "print \"FACk\""
var i = 1

var stdoutMutex = new(sync.Mutex)
var stderrMutex = new(sync.Mutex)

type Scraper struct {
  uri string
  cmd *exec.Cmd
  done chan bool
}

func NewScraper() *Scraper {
  // cmd := exec.Command(command, args...)
  fmt.Printf("FMF")
  //fmt.Println(a2)
  //cmd := exec.Command(command, a1, fmt.Sprintf(a2, i, i))
  //cmd := exec.Command(command, a1, a2)
  cmd := exec.Command(command, a1)
  i += 1
  return &Scraper{"uri", cmd, make(chan bool)}
}

func (s *Scraper) Start(wg *sync.WaitGroup) {
  s.PipeStdout()
  //s.PipeStderr()
  s.cmd.Start()

  wg.Add(1)

  go func() {
    <-s.done
    s.Kill(wg)
  }()
}

func (s *Scraper) Kill(wg *sync.WaitGroup) {
  s.cmd.Process.Kill()
  wg.Done()
}

func (s *Scraper) PipeStdout() {
  var pipe io.Reader
  var err error

  pipe, err = s.cmd.StdoutPipe()

  if err != nil {
    panic(err)
  }

  scanner := bufio.NewScanner(pipe)

  go func() {
    for scanner.Scan() {
      fmt.Printf("stdout: %s", scanner.Text())
      stdoutMutex.Lock()
      defer stdoutMutex.Unlock()

      _, err = os.Stdout.Write([]byte(scanner.Text()))

      if err != nil {
        panic(err)
      }

      _, err = os.Stdout.Write([]byte("\n"))

      if err != nil {
        panic(err)
      }
    }

    fmt.Println("done stdout")
    s.done <- true
  }()
}

func (s *Scraper) PipeStderr() {
  var pipe io.Reader
  var err error

  pipe, err = s.cmd.StderrPipe()

  if err != nil {
    panic(err)
  }

  scanner := bufio.NewScanner(pipe)

  go func() {
    for scanner.Scan() {
      fmt.Printf("stderr: %s", scanner.Text())
      stderrMutex.Lock()
      defer stderrMutex.Unlock()

      _, err = os.Stderr.Write([]byte(scanner.Text()))

      if err != nil {
        panic(err)
      }

      _, err = os.Stderr.Write([]byte("\n"))

      if err != nil {
        panic(err)
      }
    }
    fmt.Println("done stderr")
    s.done <- true
  }()
}

func main() {
  wg := &sync.WaitGroup{}
  s := NewScraper()
  s.Start(wg)
  // block until waitgroup is complete
  wg.Wait()
}
