package main

import "os"

func main() {
	chips := ParseFile(os.Args[1])
	Run(chips)
}
