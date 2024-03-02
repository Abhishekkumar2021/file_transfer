package main

import (
	"bufio"
	"os"
	"strings"
)

const (
	// the directory to store the uploaded files
	uploadDir string = "files"
	UPLOAD    string = "1"
	LIST      string = "2"
	DOWNLOAD  string = "3"
	EXIT      string = "4"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	println("Choose an option:")
	println("1. Upload")
	println("2. List")
	println("3. Download")
	println("4. Exit")

	choice, _ := reader.ReadString('\n')

	// remove the newline character
	choice = strings.TrimSpace(choice)
	choice = choice[:len(choice)-1]

	// handle the choice
	if choice == UPLOAD {
		println("Upload")
	} else if choice == LIST {
		println("List")
	} else if choice == DOWNLOAD {
		println("Download")
	} else if choice == EXIT {
		println("Exit")
	} else {
		println("Invalid choice")
	}
}