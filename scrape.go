package main

import (
  "net/http"
  "os/exec"
  "os/signal"
  "syscall"
  "sync"
  "os"
  "fmt"
  "strings"
)

const num_scraper_proc = 5
const scraper_host = "127.0.0.1"
const program = "phantomjs"
const script = "scraper.js"

const buffer_size = 4096

var current_scraper_port = 5000

var mutex = &sync.Mutex{}
func WriteStdout(r *Reader) {
  mutex.Lock()
  buf := make([]byte, buffer_size)
  os.Stdout.Write(
}

type Scraper struct {
  uri string
  cmd *exec.Cmd
}

type ScraperQueue struct {
  scrapers []Scraper
  wg sync.WaitGroup
  current int
}
asd
func NewScraperQueue() *ScraperQueue {
  scrapers := make([]Scraper, num_scraper_proc)
  var wg sync.WaitGroup

  for i := range scrapers {
    s := *NewScraper()

    s.cmd.Stdout = os.Stdout
    s.cmd.Stderr = os.Stderr

    wg.Add(1)
    go s.cmd.Run()

    scrapers[i] = s
  }

  stopSig := make(chan os.Signal, 1)

  signal.Notify(
    stopSig,
    syscall.SIGHUP,
    syscall.SIGINT,
    syscall.SIGTERM,
    syscall.SIGQUIT,
  )

  go func() {
    <-stopSig
    for _, s := range scrapers {
      s.cmd.Process.Kill()
    }
    os.Exit(0)
  }()

  return &ScraperQueue{scrapers, wg, 0}
}

func (q *ScraperQueue) Scrape(path string) error {
  q.current %= num_scraper_proc

  fmt.Println(q.scrapers[q.current].uri)
  err := q.scrapers[q.current].Scrape(path)

  q.current++

  return err
}

func uri(host string, port int) string {
  return fmt.Sprintf("http://%s:%d", host, port)
}

func NewScraper() *Scraper {
  port := current_scraper_port
  current_scraper_port++

  uri := uri(scraper_host, port)
  cmd := exec.Command(program, script, fmt.Sprintf("%d", port))

  return &Scraper{uri, cmd}
}

func (s *Scraper) Scrape(path string) error {
  buffer := strings.NewReader(path)
  _, err := http.Post(s.uri, "application/json", buffer)
  return err
}

func main() {
  q := NewScraperQueue()
  q.Scrape("http://www.audiosf.com/events/")
  q.wg.Wait()
  //q.scrapers[0].cmd.Wait()
  //s.cmd.Wait()
}
