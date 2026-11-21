package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("stockfish")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	fmt.Println("stdin:", stdin)
	fmt.Println("stdout:", stdout)

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	defer func() {
		fmt.Fprint(stdin, "quit")
		stdin.Close()
		err = cmd.Wait()
		if err != nil {
			fmt.Println("stockfish didn't exit successfully:", err)
		}
	}()

	reader := bufio.NewScanner(stdout)
	writer := bufio.NewWriter(stdin)

	send := func(command string) {
		fmt.Fprintln(writer, command)
		writer.Flush()
		fmt.Println("done with sending")
	}

	receive := func(expected string) string {
		for reader.Scan() {
			line := reader.Text()

			if strings.HasPrefix(line, expected) {
				return line
			}
		}
		panic("stockfish didn't return something expected")
	}

	send("uci")
	receive("uciok")
	send("isready")
	receive("readyok")
	fmt.Println("stockfish init successfully")
	//cli_reader := bufio.NewReader(os.Stdin)
	//cli_writer := bufio.NewWriter(os.Stdout)

	for {
		fmt.Println("enter command: ")
		var command string
		fmt.Scanf("%s", command)
		fmt.Println()
		if err != nil {
			continue
		}
		switch command {
		case "ucinewgame":
			send("ucinewgame")
			send("isready")
			receive("readyok")
		}
	}
}
