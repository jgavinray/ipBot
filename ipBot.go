package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"sync"
)

func checkIpAddressChange(c chan<- string) {
	var ip string

	for {
		resp, err := http.Get("https://ifconfig.co")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		htmlData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result := string(htmlData)

		// Only return the result if there has been a change in IP Address
		if result != ip {
			ip = result
			c <- result
        }
        time.Sleep(1 * time.Minute)

	}
}

func slackIpAddressChange(c chan string) {
	
	// ws, _ := slackConnect(os.Args[1])
	ws, _ := slackConnect("xoxb-109879420662-lyrIyOq8nsnnBcmaL8IWqOrd")
	h, _ := os.Hostname()
	for {
		msg := <-c
		fmt.Println(msg)

		m := Message{Type: "message", Channel: "C3GVB1RB8", Text: h + " : " + msg}
		postMessage(ws, m)
	
	}
}

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Fprintf(os.Stderr, "usage: ipBot slack-bot-token\n")
	// 	os.Exit(1)
	// }

	fmt.Println("Starting up...")
	
	var c chan string = make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)

	go checkIpAddressChange(c)
	go slackIpAddressChange(c)

	wg.Wait()
}