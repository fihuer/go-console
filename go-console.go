package main

import (
	"log"
	"os"
	"flag"
	"io"
	"io/ioutil"
	"github.com/fihuer/go-console/server"
)

var (
	Trace *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)



func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}


func main() {
	//Init Loggers
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//Init flags
	certPtr := flag.String("cert", "cert/cert.crt", "SSL Certificate")
	keyPtr := flag.String("key", "cert/private_key", "SSL Private Key")
	flag.Parse()

	//Creates Pages
	p1 := &server.Page{Title: "index.html", Body: []byte("This is the first page.")}
	p1.Save()
	Info.Println("Created index.html")
        p2 := &server.Page{Title: "index2.html", Body: []byte("This is the second page.")}
        p2.Save()
	Info.Println("Created index2.html")

	server.Start(*certPtr, *keyPtr)
}
