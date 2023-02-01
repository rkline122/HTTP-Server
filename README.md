# HTTP-Server

## Introduction

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

You may wish to review Section 2.2 of your textbook for details of the HTTP protocol.

<b>The server should be multithreaded -- a new thread should be spawned to deal with each client connection.
Served files</b>

Your server should work for various file types, not just text. At a minimum, test it with text and PNG images.

Your server must not be hard-coded to work with only certain file names. It needs to parse the HTTP GET request for an arbitrary file name, search for that file, and either serve it or send a 404 as appropriate.

You can make the files more complicated or interesting if you would like once you have the basic functionality of the server working.

Your server will most likely not have a domain name, so you will need to refer to it directly by its IP address. You will also need to use a nonstandard port, as described below. To have your browser request a web page from an IP address and nonstandard port, simply enter ip:port. For example, you could type the following into your web browser: http://192.168.1.70:12345/test1.html.
Note about ports

Typically, the server would listen on port 80. To ensure that we do not have many people trying to bind processes to the same port at once, you will hard-code the control port number for your server as follows:

The port number should be the last 4 digits of your G#. If the last 4 digits of your G# are less than 1024, add 10000 to the last 4 digits. For example, if the last 4 digits of your G# are 0123, your port number should be 10123.

## Extra credit

For up to 10% extra credit, add some functionality involving cookies to your server. This can be as simple as keeping a list of every page that has been visited by a particular user. Or, you can try to have more fun with it by amending the web page that you send back to them in some way based on what you know about them.

This will involve some sort of storage of user information on your server. You can keep the data in whatever format you like. It does not need to be in a real database -- a simple JSON, plaintext, or CSV file will be fine for this assignment.

Note that depending on how the browser you use to test is configured, it might choose to ignore your cookie. If so, you will need to use a different browser or change your configuration to demonstrate to me that you have successfully implemented cookies.

As an additional extra credit option, for up to 10% extra credit, you may implement a 304 Not Modified message in your server. To get credit for this, your server must

    detect when an HTTP client is making a conditional GET request (hint: look for the relevant header -- we talked about this in lab)
    check the file modification time
    either send back the file if it has been recently modified or send back a correctly formatted 304 message

You may do both the cookies and 304 extra credit if you wish.
Existing code

You may not simply call an HTTP library to perform the work for you. You should write your own code to parse the incoming message, formulate the response text, send the messages through a TCP socket, etc. However, you are welcome to use demonstration code from this course or code that you wrote for other labs or projects in this course as part of your solution. If you are not sure whether something is allowed, please ask.
What to turn in?

Turn in a copy of your project code with a written report. The report should include the following:

    screen captures demonstrating
        a successful request/response interaction
        a 404 error
        results of using cookies, if implemented
    the basic logic of the program, how you implemented each program, what problems you have encountered, and how you solved them.

The report and source code should be put into a zip file and uploaded below.

## Grading notes:

    No late work will be accepted. Be sure to start early and ask questions. If you do not have time to finish, turn in whatever work you have at the deadline.
    A server that is not multithreaded, or for which multithreading is implemented incorrectly, can earn no higher than an 80%.
