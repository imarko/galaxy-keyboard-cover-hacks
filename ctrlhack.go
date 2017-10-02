package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/termios/win"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func mapper(dst io.Writer, src io.Reader) {
	buf := make([]byte, 1)
	dollar := false
	for {
		n, err := src.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n != 1 {
			log.Println("zero input")
			return
		}
		c := buf[0]
		out := buf
		if dollar {
			switch {
			case c == ' ':
				out[0] = 0
			case c >= 'a' && c <= 'z':
				out[0] = out[0] - 96
			case c == '3':
				// f3
				out = []byte{0x1b, 'O', 'R'}
			case c == '4':
				// f4
				out = []byte{0x1b, 'O', 'S'}
			case c == '9':
				// f9
				out = []byte{0x1b, '[', '2', '0', '~'}
			default:
				// everything else just gets passed
				// through
			}
			dollar = false
		} else if c == '$' {
			dollar = true
			continue
		}

		n, err = dst.Write(out)
		if err != nil {
			log.Fatal(err)
		}
		if n != len(out) {
			log.Fatal("partial output")
			return
		}

	}
}

func setws(fd uintptr, w, h int) {
	ws := win.Winsize{Width: uint16(w), Height: uint16(h)}
	win.SetWinsize(fd, &ws)

}

func main() {
	cmd := "/data/data/com.termux/files/usr/bin/bash"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	fmt.Printf("Starting %s\n", cmd)
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	c := exec.Command(os.Args[1], os.Args[2:]...)
	pty, err := pty.Start(c)

	if err != nil {
		log.Fatal(err)
	}

	w, h, _ := terminal.GetSize(0)
	setws(pty.Fd(), w, h)

	winch := make(chan os.Signal)
	signal.Notify(winch, syscall.SIGWINCH)
	go func() {
		for range winch {
			w, h, _ := terminal.GetSize(0)
			setws(pty.Fd(), w, h)
		}
	}()

	go io.Copy(os.Stdout, pty)
	go mapper(pty, os.Stdin)
	c.Wait()
	log.Println("Exiting...")
}
