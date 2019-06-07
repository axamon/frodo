// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var userid = flag.String("user", "pippo", "username")
var password = flag.String("pass", "pippo", "password")
var remoteAddr = flag.String("r", "http://127.0.0.1:8080", "default http://127.0.0.1:8080")
var file = flag.String("file", "", "no default")

type info struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func main() {

	flag.Parse()

	// Crea il contesto.
	ctx := context.Background()

	var remoteURL string
	fmt.Println(*remoteAddr)
	remoteURL = *remoteAddr + "/upload"

	filedainviare := *file

	err := upload(ctx, remoteURL, filedainviare)
	if err != nil {
		log.Println(err.Error())
	}
}

func upload(ctx context.Context, url string, filedainviare string) (err error) {
	file, err := os.Open(filedainviare)
	if err != nil {
		log.Printf("impossible aprire file: %s errore: %s\n", filedainviare, err.Error())
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	if err != nil {
		log.Println(err.Error())
	}

	kvPairs, err := json.Marshal(info{Name: filedainviare, Data: encoded})

	//fmt.Printf("Sending JSON string '%s'\n", string(kvPairs))

	// Send request to OP's web server
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(kvPairs))
	if err != nil {
		log.Printf(err.Error())
	}

	req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")

	//Aggiunge sicurezza
	req.SetBasicAuth(*userid, *password)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	//body, err := ioutil.ReadAll(resp.Body)

	//fmt.Println("Response: ", string(body))
	return
}
