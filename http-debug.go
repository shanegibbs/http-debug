package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Println("Starting")

	store := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], "HTTP_") {
			store[strings.ToLower(pair[0][5:])] = pair[1]
		}
	}

	log.Println("Initial store:", store)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Create return string
		var request []string

		{
			value, present := store["failure_chance"]
			if present {
				if s, err := strconv.ParseFloat(value, 64); err == nil {
					failureChance := s
					if rand.Float64() > 1-failureChance {
						http.Error(w, "Random error\n", http.StatusInternalServerError)
					}
				}
			}
		}

		{
			value, present := store["sleep_mills"]
			if present {
				if s, err := strconv.ParseInt(value, 10, 64); err == nil {
					sleepTime := time.Duration(s) * time.Millisecond
					time.Sleep(sleepTime)
				}
			}
		}

		{
			value, present := store["request_urls"]
			if present {
				response, err := http.Get(value)
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					contents, err := ioutil.ReadAll(response.Body)
					if err != nil {
						http.Error(w, fmt.Sprintf("Failed to get url: %s\n", value), http.StatusInternalServerError)
					}
					request = append(request, fmt.Sprintf("From %s:\n%s\n", value, string(contents)))
				}
			}
		}

		request = append(request, fmt.Sprintf("Remote: %s\n", r.RemoteAddr))

		// Add the request string
		url := fmt.Sprintf("%v %v %v\n", r.Method, r.URL, r.Proto)
		request = append(request, url)

		// Add the host
		request = append(request, fmt.Sprintf("Headers\n-------\nHost: %v", r.Host))

		// Loop through headers
		var headerNames []string
		for name, _ := range r.Header {
			headerNames = append(headerNames, name)
		}
		sort.Strings(headerNames)
		// for name, headers := range r.Header {
		for _, name := range headerNames {
			for _, h := range r.Header[name] {
				request = append(request, fmt.Sprintf("%v: %v", name, h))
			}
		}

		// If this is a POST, add post data
		if r.Method == "POST" {
			r.ParseForm()
			request = append(request, fmt.Sprintf("\nForm\n----"))
			request = append(request, r.Form.Encode())

			key, keyOk := r.Form["key"]
			value, valueOk := r.Form["value"]

			if keyOk && valueOk {
				store[key[0]] = value[0]
			}
		}

		request = append(request, fmt.Sprintf("\nStore\n-----\n"))
		var keyNames []string
		for name, _ := range store {
			keyNames = append(keyNames, name)
		}
		for _, name := range keyNames {
			request = append(request, fmt.Sprintf("[%s] %s", name, store[name]))
		}

		fmt.Fprintf(w, "%s\n", strings.Join(request, "\n"))
	})

	// listen on port 8080
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
