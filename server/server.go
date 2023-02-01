/*
Develop a simple, multithreaded web server in the language of your choice. Your server will be "simple" in that it

  - can ignore most HTTP headers,
  - needs only to be able to return responses of type 200 OK or 404 Not Found (with any corresponding data), and
  - need not operate with persistent connections.

For the most part, however, it will be a fully functional web server.

The general process of your web server should go as follows:

  - listen on a specified port
  - create a connection socket when a browser contacts you
  - receive an HTTP request through the socket
  - parse the request to determine which file is being requested
  - if the file is not available, respond with a 404 error. Otherwise, proceed to the next steps.
  - create an HTTP response -- the contents of the file preceded by header lines
  - send the HTTP response over the connection
*/
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
        "errors"
        "strconv"
        "bytes"
)
const (
        SERVER_HOST = "localhost"
        SERVER_PORT = "9988"
        SERVER_TYPE = "tcp"
)
func main() {
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
                go processClient(connection)    // Multi-threading with goroutines
        }
}


func processClient(connection net.Conn) {
        /*

        Processes the request sent from the client
        and sends an appropriate response based on
        the validity of their request.

        */
        buffer := make([]byte, 1024)
        mLen, err := connection.Read(buffer)

        if err != nil {
                fmt.Println("Error reading:", err.Error())
                return
        }

        // Deconstruct the request to eventually retrieve the desired filepath
        bufferToString := string(buffer[:mLen])
        headers := strings.Split(bufferToString, "\r\n")
        request := strings.Fields(headers[0])
        filePath := strings.TrimLeft(request[1], "/")
        
        fmt.Println("File Requested: ", filePath)

        var (
                header string
                content string
        )

        if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
                content = GetFileContent("content/notFound.html")
                header += "HTTP/1.1 404 Not Found\r\n"
                header += "Content-Type: text/html\r\n"
                header += "Content-Length: " + strconv.Itoa(len(content)) + "\r\n"
                header += "\r\n"
        }else{
                content = GetFileContent(filePath)
                header += "HTTP/1.1 200 OK\r\n"
                header += "Content-Type: text/html\r\n"
                header += "Content-Length: " + strconv.Itoa(len(content)) + "\r\n"
                header += "\r\n"
        }

        response := header + content
        connection.Write([]byte(response))
        connection.Close()
}


func GetFileContent(fileURL string) string{
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