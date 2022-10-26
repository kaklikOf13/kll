package main

import (
	"os"

	"github.com/kaklikOf13/kll"
)

func main() {
	i := kll.Interpreter{}
	if len(os.Args) == 1 {
		print("comandos:\n   run:\n      roda o programa\n")
		return
	}
	if os.Args[1] == "run" || os.Args[1] == "debug-run" {
		if os.Args[1] == "debug-run" {
			i.Debug = true
		}
		if len(os.Args) >= 3 {
			i.Exec_Main(os.Args[2])
		} else {
			i.Exec_Main("main.kll")
		}
	}

}

//go run main.go run
//go run main.go debug-run
//go build -buildmode=c-shared -o kll.so main.go
//go build main.go
//GOOS=windows GOARCH=amd64 go build -o kll.exe main.go
