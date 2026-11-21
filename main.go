package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	tenv "github.com/muesli/termenv"
)

func main() {
	var s Stockfish
	err := s.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("Stockfish init successfully")

	output := tenv.NewOutput(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("1. Newgame")
		fmt.Print("what to do?: ")
		var option int
		scanner.Scan()
		option, err = strconv.Atoi(scanner.Text())
		fmt.Println()
		if err != nil {
			continue
		}
		switch option {
		case 1:
			s.Send("ucinewgame")
			s.Send("isready")
			_, err = s.Receive("readyok")
			if err != nil {
				str := output.String("couldn't start new game: " + err.Error())
				fmt.Println(str.Foreground(output.Color("#FFCCCC")))
			} else {
				str := output.String("done")
				fmt.Println(str.Foreground(output.Color("#AAFFAA")))
			}
		}
	}
}
