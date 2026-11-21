package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Stockfish struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	Writer *bufio.Writer
	Reader *bufio.Scanner
	cmd    *exec.Cmd
}

func (s *Stockfish) Init() error {
	s.cmd = exec.Command("stockfish")

	var err error
	s.stdin, err = s.cmd.StdinPipe()
	if err != nil {
		return err
	}
	s.stdout, err = s.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = s.cmd.Start(); err != nil {
		return err
	}

	s.Reader = bufio.NewScanner(s.stdout)
	s.Writer = bufio.NewWriter(s.stdin)

	s.Send("uci")
	s.Receive("uciok")
	s.Send("isready")
	s.Receive("readyok")
	return nil
}

func (s *Stockfish) Quit() error {
	fmt.Fprint(s.stdin, "quit")
	s.stdin.Close()
	err := s.cmd.Wait()
	if err != nil {
		return fmt.Errorf("stockfish didn't exit successfully: %w", err)
	}
	return nil
}

func (s Stockfish) Send(command string) {
	fmt.Fprintln(s.Writer, command)
	s.Writer.Flush()
}

func (s Stockfish) Receive(expected string) (string, error) {
	for s.Reader.Scan() {
		line := s.Reader.Text()

		if strings.HasPrefix(line, expected) {
			return line, nil
		}
	}
	return "", fmt.Errorf("stockfish didn't answer the expected")
}
