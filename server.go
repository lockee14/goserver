/*
	TO LEARN ABOUT GOLANG:
	interface
	channel
	context
*/

/*
	TODO:
	handle every route
	write a logger
	test
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Conf struct {
	port string
	path string
}

var port = flag.String("p", "port", "port to serve on")
var dir = flag.String("d", "dir", "the directory of static file to host")

func servFile(res http.ResponseWriter, req *http.Request) {
	f, e := os.Open(*dir + "\\index.html")
	defer f.Close()
	if e != nil {
		log.Fatal(e)
	}
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		// end of file
		if n == 0 {
			fmt.Println("print end of file: ", err)
			break
		}

		// write a chunk
		if _, err := res.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	return
}

func main() {
	flag.Parse()
	log.Printf("Serving %s on HTTP port: %s\n", *dir, *port)
	http.Handle("/", http.FileServer(http.Dir(*dir)))
	http.HandleFunc("/market/", marketHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
