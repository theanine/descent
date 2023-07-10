package main

import (
	"fmt"
	"net/http"
	"os"
)

const version = "v4.0.0.181218"
const downloadEnabled = false

func usage() {
	fmt.Println(`Usage: heroes <OPTION>

Options:
  -s, --serve      Start local server on localhost:8000
  -h, --heroes     Generate heroes.html
  -c, --classes    Generate classes.html
  -o, --overlord   Generate overlord.html
  -p, --plot       Generate plot.html
  -i, --items      Generate items.html
  -k, --console    Generate console.html`)
	os.Exit(1)
}

func serve() {
	fmt.Printf("Server started on localhost:8000\n")
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 1 {
		usage()
	}
	startServer := false
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "-s" || arg == "--serve" {
			startServer = true
		}
		if arg == "-h" || arg == "--heroes" {
			heroesGen()
		}
		if arg == "-c" || arg == "--classes" {
			classesGen()
		}
		if arg == "-o" || arg == "--overlord" {
			overlordGen()
		}
		if arg == "-p" || arg == "--plot" {
			plotGen()
		}
		// if arg == "-i" || arg == "--items" {
		// 	itemsGen()
		// }
		if arg == "-k" || arg == "--console" {
			consoleGen()
		}
	}
	if startServer == true {
		serve()
	}
}
