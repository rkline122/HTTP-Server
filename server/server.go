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

import(
    "fmt"
    "net"
    "os"
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
                go processClient(connection)    // Starts goroutine for new connection
        }
}


func processClient(connection net.Conn) {
        buffer := make([]byte, 1024)
        mLen, err := connection.Read(buffer)
        if err != nil {
                fmt.Println("Error reading:", err.Error())
        }
        fmt.Println("Received: ", string(buffer[:mLen]))
        _, _ = connection.Write([]byte("Thanks! Got your message: [" + string(buffer[:mLen]) + "]"))
        connection.Close()
}