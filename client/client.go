package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	// the directory to store the uploaded files
	uploadDir   string = "uploads"
	downloadDir string = "downloads"
	UPLOAD      string = "1"
	LIST        string = "2"
	DOWNLOAD    string = "3"
	EXIT        string = "4"
)

var conn net.Conn

func main() {
	setupDir()
	connectToServer()
	handleMenu()
}

func handleMenu() {
	// create a new menu
	menu := []string{
		"1. Upload a file",
		"2. List the files",
		"3. Download a file",
		"4. Exit",
	}

	// print the menu
	for {
		// If connection is closed, then stop the client gracefully
		if conn == nil {
			fmt.Println("Connection is closed by the server")
			os.Exit(0)
		}

		for _, item := range menu {
			fmt.Println(item)
		}

		fmt.Println("")

		fmt.Print("Enter your choice: ")

		// read the user's choice
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')

		// remove the newline character
		choice = strings.TrimSpace(choice)

		// handle the choice
		switch choice {
		case UPLOAD:
			handleUpload()
		case LIST:
			handleList()
		case DOWNLOAD:
			handleDownload()
		case EXIT:
			conn.Write([]byte(EXIT + "\n"))
			conn.Close() // Close connection when exiting
			os.Exit(0)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func setupDir() {
	// create the uploads directory if not exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// create the downloads directory if not exists
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		os.Mkdir(downloadDir, 0755)
	}
}

func connectToServer() {
	// connect to the server
	fmt.Print("Enter the server address, e.g. 127.0.0.1:8080: ")
	reader := bufio.NewReader(os.Stdin)
	address, _ := reader.ReadString('\n')

	// remove the newline character
	address = strings.TrimSpace(address)

	fmt.Println("Connecting to the server", address, "...")

	// connect to the server
	var err error
	conn, err = net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to the server", err.Error())
		os.Exit(1)
	}

	fmt.Println("✅ Connected to the server", address)
	fmt.Println("You can now upload, list, and download files")
	fmt.Println("")
}

func handleUpload() {
	conn.Write([]byte(UPLOAD + "\n"))
	
	// create a reader
	reader := bufio.NewReader(os.Stdin)

	// read the file name
	fmt.Print("Enter the file name: ")
	fileName, _ := reader.ReadString('\n')

	// remove the newline character
	fileName = strings.TrimSpace(fileName)

	_, err := conn.Write([]byte(fileName + "\n"))
	if err != nil {
		fmt.Println("Error sending file name:", err)
		return
	}

	// read the file content
	fileContent, err := os.ReadFile(uploadDir + "/" + fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Send the file size first
	fileSize := strconv.Itoa(len(fileContent))
	_, err = conn.Write([]byte(fileSize + "\n"))
	if err != nil {
		fmt.Println("Error sending file size", err)
		return
	}

	// send the file content to the server
	_, err = conn.Write(fileContent)
	if err != nil {
		fmt.Println("Error sending file content:", err)
		return
	}

	// log the file is uploaded
	fmt.Println("File is uploaded to the server")
	fmt.Println("")
}

func handleList() {
	// send the option to the server
	conn.Write([]byte(LIST + "\n"))

	// create a reader
	reader := bufio.NewReader(conn)

	fmt.Println("Files in the server:")

	// read the file names in one go
	fileNames, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file names:", err)
		return
	}

	// Split the file names based on \t
	names := strings.Split(strings.TrimSpace(fileNames), "\t")
	for idx, name := range names {
		fmt.Println(idx+1, name)
	}

	// log the file names are received
	fmt.Println("File names are received from the server")
	fmt.Println("")
}

func handleDownload() {
	// create a reader
	reader := bufio.NewReader(os.Stdin)

	// read the file name
	fmt.Print("Enter the file name: ")
	fileName, _ := reader.ReadString('\n')

	// remove the newline character
	fileName = strings.TrimSpace(fileName)

	// send the option to the server
	conn.Write([]byte(DOWNLOAD + "\n"))

	// send the file name to the server
	conn.Write([]byte(fileName + "\n"))

	// read the file size from the server
	fileSizeStr, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}
	fileSizeStr = strings.TrimSpace(fileSizeStr)
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		fmt.Println("Error parsing file size:", err)
		return
	}

	// create a file to store the downloaded file
	file, err := os.Create(downloadDir + "/" + fileName)
	if err != nil {
		fmt.Println("Error creating the file ", err.Error())
		return
	}
	defer file.Close()

	// receive and write file content
	written, err := io.CopyN(file, conn, int64(fileSize))
	if err != nil {
		fmt.Println("Error receiving file content:", err)
		return
	}

	if written != int64(fileSize) {
		fmt.Println("Incomplete file transfer")
		return
	}

	// log the file is downloaded
	fmt.Println("File is downloaded from the server")
	fmt.Println("")
}
