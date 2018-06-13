package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
)

const (
	crlf       = "\r\n"
	colonspace = ": "
)

func ChecksumMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// record responses for middleware
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, r)

		// read and sort headers from response
		var response = []string{}
		var headers = recorder.Header()
		for key, _ := range headers {
			response = append(response, key)
		}
		sort.Strings(response)

		// add hash, append, and join headers
		hash := sha1.New()
		io.WriteString(hash, fmt.Sprintf("%v", recorder.Code)+crlf)
		for _, header := range response {
			io.WriteString(hash, header+colonspace+strings.Join(headers[header], "")+crlf)
		}

		// write response and format headers in prep for checsum
		var body, _ = ioutil.ReadAll(recorder.Body)
		io.WriteString(hash, "X-Checksum-Headers: "+strings.Join(response, ";")+crlf+crlf+string(body))

		// add final checksum header to cannonical responses
		hashChecksum := hash.Sum(nil)
		w.Header().Set("X-Checksum", hex.EncodeToString(hashChecksum))

		h.ServeHTTP(w, r)
	})
}

// Do not change this function.
func main() {
	var listenAddr = flag.String("http", "localhost:8080", "HTTP Listener")
	flag.Parse()

	http.Handle("/", ChecksumMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Foo", "bar")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Date", "Wednesday, 13 June 2018 14:04:53 GMT")
		msg := "Testy Test.\n"
		w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
		fmt.Fprintf(w, msg)
	})))

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
