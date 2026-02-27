package ui

import (
	"fmt"
	"os"
	"time"
)

type Loader struct {
	char_set  []string
	delay     time.Duration
	stop_chan chan byte
}

func NewLoader(char_set []string, delay uint16) *Loader {
	return &Loader{
		char_set:  char_set,
		delay:     (time.Millisecond) * time.Duration(delay),
		stop_chan: make(chan byte),
	}
}

func (l *Loader) Start() {
	idx := 0
	n := len(l.char_set)
	go func() {
		for {
			i := idx % n
			select {
			case <-l.stop_chan:
				{
					fmt.Fprintf(os.Stdout, "\033[2K\r")
					return
				}
			default:
				{
					fmt.Fprintf(os.Stdout, "\033[2K\r")
					fmt.Fprintf(os.Stdout, "%s", l.char_set[i])
				}
			}
			idx++
			time.Sleep(l.delay)
		}
	}()
}

func (l *Loader) Stop() {
	l.stop_chan <- 'x'
	close(l.stop_chan)
}
