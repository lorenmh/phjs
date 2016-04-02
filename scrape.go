package main

import (
  "net/http"
//  "net"
//  "bufio"
  "os/exec"
  "os"
  "fmt"
  "strings"
)

const num_scraper_proc = 5
const scraper_host = "127.0.0.1"
const program = "phantomjs"
const script = "scraper.js"

var current_scraper_port = 5000

type Scraper struct {
  uri string
  cmd *exec.Cmd
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
  _, error := http.Post(s.uri, "application/json", buffer)
  return error
}

func main() {
  s := NewScraper()
  fmt.Println(program)
  fmt.Println(s.uri)
  s.cmd.Stdout = os.Stdout
  s.cmd.Stderr = os.Stderr
  go s.cmd.Run()
  s.Scrape("http://www.audiosf.com/events/")
  s.Scrape("http://www.audiosf.com/events/")
  s.Scrape("http://www.audiosf.com/events/")
  s.Scrape("http://www.audiosf.com/events/")
  s.Scrape("http://www.audiosf.com/events/")
  s.Scrape("http://www.audiosf.com/events/")
  s.cmd.Wait()
}
