package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
  "sort"
  "strconv"
)

func main() {
  store := make(map[string]string)
  store["failureChance"] = "0"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// rnd := rand.New(rand.NewSource(0))

    failureChance := 0.0
    if s, err := strconv.ParseFloat(store["failureChance"], 64); err == nil {
      failureChance = s
    }
		if rand.Float64() > 1 - failureChance {
			http.Error(w, "Random error\n", http.StatusInternalServerError)
    }

		// Create return string
		var request []string

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

    // if accept, ok := r.Header["Accept"]; ok {
    //   if strings.Contains(accept[0], "html") {
    //     request = append(request, fmt.Sprintf("<h4>Update</h4>"))
    //   }
    // }

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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
