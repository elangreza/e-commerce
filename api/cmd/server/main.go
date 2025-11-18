package main

import (
	"log"
)

func main() {

}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
