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
	uploadDir string = "assets"
	UPLOAD    string = "1"
	LIST      string = "2"
	DOWNLOAD  string = "3"
	EXIT      string = "4"
)

func main() {
	// create the directory if it does not exist
	setupDir()

	// create a listener
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	// close the listener when the server is stopped
	defer listener.Close()

	// log the server is running
	fmt.Printf("✅ Server is running at %s\n", listener.Addr())

	// accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			continue
		}
		// handle the connection concurrently
		go handleConnection(conn)
	}
}

func setupDir() {
	// create the directory if it does not exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("👌 Connected to the client with address", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading option", err.Error())
			return
		}

		option = strings.TrimSpace(option)

		switch option {
		case UPLOAD:
			if err := handleUpload(reader, conn); err != nil {
				fmt.Println("Error handling upload", err.Error())
				return
			}
		case LIST:
			if err := handleList(conn); err != nil {
				fmt.Println("Error handling list", err.Error())
				return
			}
		case DOWNLOAD:
			if err := handleDownload(conn, reader); err != nil {
				fmt.Println("Error handling download", err.Error())
				return
			}
		case EXIT:
			fmt.Println("Client requested to exit")
			return
		default:
			fmt.Println("Invalid option")
			return
		}
	}
}

// handleUpload function
func handleUpload(reader *bufio.Reader, conn net.Conn) error {
	fileName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file name", err.Error())
		return err
	}
	fileName = strings.TrimSpace(fileName)
	fmt.Println("Received file name", fileName)

	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return err
	}

	fileSizeStr = strings.TrimSpace(fileSizeStr)
	fileSize, err := strconv.Atoi(fileSizeStr)
	
	if err != nil {
		fmt.Println("Error parsing file size:", err)
		return err
	}


	// create a new file
	file, err := os.Create(uploadDir + "/" + fileName)
	if err != nil {
		fmt.Println("Error creating the file ", err.Error())
		return err
	}
	defer file.Close()

	written, err := io.CopyN(file, conn, int64(fileSize))
	if err != nil {
		fmt.Println("Error receiving file content:", err)
		return err
	}

	if written != int64(fileSize) {
		fmt.Println("Incomplete file transfer")
		return nil
	}

	// log the file is uploaded
	fmt.Println("File is added to the server")

	return nil
}

func handleDownload(conn net.Conn, reader *bufio.Reader) error {
	fileName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file name", err.Error())
		return err
	}
	fileName = strings.TrimSpace(fileName)
	fmt.Println("Received file name", fileName)

	fileContent, err := os.ReadFile(uploadDir + "/" + fileName)
	if err != nil {
		fmt.Println("Error reading file content", err.Error())
		return err
	}

	// Send the file size first
	fileSize := strconv.Itoa(len(fileContent))
	_, err = conn.Write([]byte(fileSize + "\n"))
	if err != nil {
		fmt.Println("Error sending file size", err.Error())
		return err
	}

	// Then send the file content
	_, err = conn.Write(fileContent)
	if err != nil {
		fmt.Println("Error sending file content", err.Error())
		return err
	}

	fmt.Println("File is sent to the client")
	return nil
}

func handleList(conn net.Conn) error {
	fileNames, err := os.ReadDir(uploadDir)
	if err != nil {
		fmt.Println("Error reading file names", err.Error())
		return err
	}

	var fileList string
	for _, file := range fileNames {
		fileList += file.Name() + "\t"
	}

	fileList += "\n"

	_, err = conn.Write([]byte(fileList))
	if err != nil {
		fmt.Println("Error sending file names", err.Error())
		return err
	}

	fmt.Println("File names are sent to the client")
	return nil
}
