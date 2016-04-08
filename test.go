package main

import (
  "bufio"
  "strings"
  "sync"
  "os"
  "os/exec"
  "os/signal"
  "io"
  "syscall"
  "fmt"
  "net/http"
  "time"
)

// phantomjs --proxy=ip:port script.js
const host = "127.0.0.1"
const command = "phantomjs"
var arg1 = "scraper.js"
var port = 5000

var stdoutMutex = new(sync.Mutex)
var stderrMutex = new(sync.Mutex)

var _ = fmt.Printf

type Scraper struct {
  uri string
  cmd *exec.Cmd
}

func createUri(port int) string {
  return fmt.Sprintf("http://%s:%d", host, port)
}

func NewScraper() *Scraper {
  uri := createUri(port)
  cmd := exec.Command(command, arg1, fmt.Sprintf("%d", port))

  port += 1

  return &Scraper{uri, cmd}
}

// message can be a path, a JSON object, whatever.  The scrape script will
// handle whatever is passed to it
func (s *Scraper) Scrape(message string) {
  buffer := strings.NewReader(message)

  response, _ := http.Post(s.uri, "application/json", buffer)

  fmt.Println(response)
  //fmt.Println(err)

  //if err != nil {
    //panic(err)
  //}
}

func (s *Scraper) Start(wg *sync.WaitGroup) {
  stdoutDone := s.PipeStdout()
  stderrDone := s.PipeStderr()

  startErr := s.cmd.Start()
  if startErr != nil {
    panic(startErr)
  }

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
  fmt.Println("KILLING")
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
  //s2 := NewScraper()

  s1.Start(wg)
  //s2.Start(wg)

  time.Sleep(1000 * time.Millisecond)
  str := "http://www.audiosf.com/events/"
  s1.Scrape(str)

  wg.Wait()
}
