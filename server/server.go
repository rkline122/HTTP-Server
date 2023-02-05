/*
Project 1: HTTP Server
By Ryan Kline
	---
CIS 457 - Data Communications
Winter 2023
*/
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8636"
	SERVER_TYPE = "tcp"
)

func main() {
	/*

	   Starts up server using the host, port, and
	   protocol defined above. Once a client is connected,
	   the processClient() function is ran as a goroutine (multithread)

	*/
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection)
	}
}

type Cookie struct {
	ID      string
	History map[string]string
}

func processClient(connection net.Conn) {
	/*

	   Processes the request sent from the client and sends
	   an appropriate response based on the validity of their request.

	*/
	var (
		responseCode  string
		header        string
		content       string
		newConnection bool
		clientCookie  Cookie
		cookies       = make(map[string]Cookie)
		buffer        = make([]byte, 1024)
	)
	rand.Seed(time.Now().UnixNano())

	// Opens JSON file containing cookie data
	jsonFile, err := os.Open("cookies.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Loads JSON data into cookies map
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &cookies)

	// Reads and deconstructs client message
	messageLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	bufferToString := string(buffer[:messageLen])
	headers := strings.Split(bufferToString, "\r\n")
	request := strings.Fields(headers[0])
	filePath := strings.TrimLeft(request[1], "/")
	fmt.Println("File Requested: ", filePath)

	// Checks for cookie header
	for _, header := range headers {
		if strings.HasPrefix(header, "Cookie: ") {
			data := header[len("Cookie: "):]
			id := (strings.Split(data, "="))[1]
			newConnection = false
			clientCookie.ID = id
			break
		} else {
			cookieID := randSeq(5)
			clientCookie.ID = cookieID
			newConnection = true
		}
	}

	if newConnection {
		clientCookie.History = make(map[string]string)
	} else {
		clientCookie.History = cookies[clientCookie.ID].History
	}

	// Check if requested filepath exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		filePath = "content/notFound.html"
		responseCode = "HTTP/1.1 404 Not Found\r\n"
		header += "Content-Type: text/html\r\n"
	} else {
		responseCode = "HTTP/1.1 200 OK\r\n"

		// Check file type
		if strings.HasSuffix(filePath, ".html") {
			header += "Content-Type: text/html\r\n"
		} else if strings.HasSuffix(filePath, ".jpg") {
			header += "Content-Type: image/jpeg\r\n"
		} else {
			filePath = "content/isDirectory.html"
			header += "Content-Type: text/html\r\n"
		}
		clientCookie.History[filePath] = time.Now().Format("Mon, 02 Jan 2006 15:04:05")
	}
	content = GetFileContent(filePath)

	if newConnection {
		header += fmt.Sprintf("Set-Cookie: id=%s\r\n", clientCookie.ID)
	}
	header += "Content-Length: " + strconv.Itoa(len(content)) + "\r\n"
	header += "\r\n"

	// Update cookies map and write to JSON file
	cookies[clientCookie.ID] = clientCookie
	writeCookieToJSONFile(cookies, "cookies.json")

	// Send response and close connection
	response := responseCode + header + content
	connection.Write([]byte(response))
	connection.Close()
}

func GetFileContent(fileURL string) string {
	/*

	   Helper function to retrieve contents of a
	   file in the form of a string if it exists.

	*/
	file, err := os.Open(fileURL)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
	}
	defer file.Close()

	fileBuffer := new(bytes.Buffer)
	fileBuffer.ReadFrom(file)
	content := fileBuffer.String()

	return content
}

func writeCookieToJSONFile(cookies map[string]Cookie, fileName string) error {
	// Write cookie data to JSON
	file, err := json.MarshalIndent(cookies, "", " ")
	err = ioutil.WriteFile(fileName, file, 0644)

	return err
}

func randSeq(n int) string {
	// Generate a random sequence of characters of length n
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
