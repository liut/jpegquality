package main

import (
	"fmt"
	"log"
	"os"

	"github.com/liut/jpegquality"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jpegq file.jpg")
		return
	}

	jpegquality.SetLogger(log.New(os.Stderr, "jpegq", log.LstdFlags|log.Lshortfile))

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	jr, err := jpegquality.New(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(jr.Quality())
}
