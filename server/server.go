/*

Project 1: HTTP Server
By Ryan Kline
   ---
CIS 457 - Data Communications
Winter 2023

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
        "encoding/json"
	"time"
        "io/ioutil"
)
const (
        SERVER_HOST = "localhost"
        SERVER_PORT = "8000"
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
        ID string
        History []string
}

var cookies = make(map[string]Cookie)

func processClient(connection net.Conn) {
        /*

        Processes the request sent from the client and sends 
        an appropriate response based on the validity of their request.

        */
        var clientCookie Cookie
        hasCookie := false
        buffer := make([]byte, 1024)
        mLen, err := connection.Read(buffer)

        if err != nil {
                fmt.Println("Error reading:", err.Error())
                return
        }

        bufferToString := string(buffer[:mLen])
        headers := strings.Split(bufferToString, "\r\n")
        request := strings.Fields(headers[0])
        filePath := strings.TrimLeft(request[1], "/")
        
        fmt.Println("File Requested: ", filePath)
        
        var data string
        // Checks for cookie header
        for _, header := range headers {
                if strings.HasPrefix(header, "Cookie: ") {
                        data = header[len("Cookie: "):]
                        id:= (strings.Split(data, "="))[1]
                        fmt.Println(id)
                        hasCookie = true
                        clientCookie.ID = id
                        break
                }
        }

        var (
                header string
                content string
        )

        var historyStr string
        if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
                content = GetFileContent("content/notFound.html")
                header += "HTTP/1.1 404 Not Found\r\n"
                header += "Content-Type: text/html\r\n"

                historyStr = fmt.Sprintf("Visited 'notFound.html' at %s", time.Now().Format(time.RFC3339))
        }else{
                content = GetFileContent(filePath)
                header += "HTTP/1.1 200 OK\r\n"
                historyStr = fmt.Sprintf("Visited '%s' at %s", filePath, time.Now().Format(time.RFC3339))
                

                if(strings.HasSuffix(filePath, ".html")){
                        header += "Content-Type: text/html\r\n"
                }else if(strings.HasSuffix(filePath, ".jpg")){
                        header += "Content-Type: image/jpeg\r\n"
                }else{
                        content = GetFileContent("content/isDirectory.html")
                        header += "Content-Type: text/html\r\n"           
                }
                
        }

        fmt.Println(cookies[clientCookie.ID].History)
        clientCookie.History = append(cookies[clientCookie.ID].History, historyStr)
        fmt.Println(clientCookie.History)

        cookies[clientCookie.ID] = clientCookie

        if(!hasCookie){
                cookieID := "a1b2c3"           //TODO: Make this random
                header += fmt.Sprintf("Set-Cookie: id=%s\r\n", cookieID)
        }
        header += "Content-Length: " + strconv.Itoa(len(content)) + "\r\n" 
        header += "\r\n"


        file, _ := json.MarshalIndent(cookies, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)

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


func writeCookieToJSONFile(cookie Cookie) error {
        file, err := os.Create("cookies.json")
        if err != nil {
            return err
        }
        defer file.Close()
    
        encoder := json.NewEncoder(file)
        return encoder.Encode(cookie)
}
    

func readCookiesFromJSONFile() ([]Cookie, error) {
        file, err := os.Open("cookies.json")
        if err != nil {
            return nil, err
        }
        defer file.Close()
    
        decoder := json.NewDecoder(file)
        var cookies []Cookie
        err = decoder.Decode(&cookies)
        if err != nil {
            return nil, err
        }
        return cookies, nil
    }
    
