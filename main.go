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
  -h, --heroes     Generate heroes.html
  -c, --classes    Generate classes.html
  -o, --overlord   Generate overlord.html
  -p, --plot       Generate plot.html
  -i, --items      Generate items.html
  -k, --console    Generate console.html`)
	os.Exit(1)
}

func serve() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 1 || len(os.Args) > 2 {
		usage()
	}
	if len(os.Args) == 2 && (os.Args[1] == "-s" || os.Args[1] == "--serve") {
		serve()
	}
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--heroes" {
		heroesGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-c" || os.Args[1] == "--classes" {
		classesGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-o" || os.Args[1] == "--overlord" {
		overlordGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-p" || os.Args[1] == "--plot" {
		plotGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-i" || os.Args[1] == "--items" {
		itemsGen()
	}
	if len(os.Args) == 1 || os.Args[1] == "-k" || os.Args[1] == "--console" {
		consoleGen()
	}
}
