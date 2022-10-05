package main

/* Usage:
 * PORT=10001 TARGET=http://proxy.target:1234 go run proxy.go
 */

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var port string
var target string

func init() {
	port = os.Getenv("PORT")
	target = os.Getenv("TARGET")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[INFO] Request to %s\n", r.URL)

	url := fmt.Sprintf("%s%s", target, r.URL.Path)

	var b bytes.Buffer
	b.ReadFrom(r.Body)
	r.Body.Close()
	body := ioutil.NopCloser(&b)

	req, err := http.NewRequest(r.Method, url, body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("[INFO] Proxy...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	respHeader := resp.Header
	wHeader := w.Header()

	for k := range respHeader {
		wHeader.Set(k, respHeader.Get(k))
	}
	wHeader.Set("Access-Control-Allow-Origin", "*")
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func main() {
	fmt.Println("[INFO] Add Handler")
	http.HandleFunc("/", handler)

	fmt.Println("[INFO] Listening...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	fmt.Println("[INFO] Terminated!")
}
