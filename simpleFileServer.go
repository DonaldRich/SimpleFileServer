// Copyright (c) 2023 Donald Rich

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var rootPath = ""

var allowedClients = []string{"localhost"}

var mimeTypes = map[string]string{
	"bin":  "application/octet-stream",
	"bmp":  "image/bmp",
	"css":  "text/css",
	"csv":  "text/csv",
	"gif":  "image/gif",
	"htm":  "text/html",
	"html": "text/html",
	"ico":  "image/x-icon",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"js":   "text/javascript",
	"json": "application/json",
	"mjs":  "text/javascript",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"txt":  "text/plain",
	"webp": "image/webp",
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Usage: sfs <port> <path>")
		os.Exit(1)
	}
	portNumber, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		fmt.Printf("Invalid port")
		os.Exit(2)
	}
	rootPath = args[1]
	pathInfo, err := os.Stat(rootPath)
	if os.IsNotExist(err) || !pathInfo.IsDir() {
		fmt.Printf("Invalid path: %s\n", rootPath)
		os.Exit(3)
	}

	var moreMimeTypes map[string]string

	if jsonData, err := os.ReadFile("mime-types.json"); err == nil {
		if json.Unmarshal([]byte(jsonData), &moreMimeTypes) == nil {
			for extension, mimeType := range moreMimeTypes {
				mimeTypes[extension] = mimeType
			}
			fmt.Printf("MimeTypes added\n")
		} else {
			fmt.Printf("Failed to add MimeTypes: %s\n", err)
		}
	}

	var additionalClients []string

	if jsonData, err := os.ReadFile("clients.json"); err == nil {
		if json.Unmarshal([]byte(jsonData), &additionalClients) == nil {
			if len(additionalClients) == 1 && additionalClients[0] == "*" {
				allowedClients[0] = "*"
				fmt.Printf("All clients allowed\n")
			} else {
				allowedClients = append(allowedClients, additionalClients...)
				fmt.Printf("Clients added\n")
			}
		} else {
			fmt.Printf("Failed to add Clients: %s\n", err)
		}
	}

	http.HandleFunc("/", handleRequest)

	fmt.Printf("Starting server\n")
	err = http.ListenAndServe(":"+args[0], nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else {
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(4)
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if len(allowedClients) != 1 && allowedClients[0] != "*" {
		client := r.Header.Get("X-FORWARDED-FOR")
		if client == "" {
			client = r.RemoteAddr
		}
		if colonIndex := strings.Index(client, ":"); colonIndex != -1 {
			client = client[0:colonIndex]
		}

		clientAllowed := false
		for _, allowedClient := range allowedClients {
			if client == allowedClient {
				clientAllowed = true
				break
			}
		}
		if !clientAllowed {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	filename := r.RequestURI[1:]
	filePath := rootPath + string(os.PathSeparator) + filename
	if contents, err := os.ReadFile(filePath); err == nil {
		extension := strings.ToLower(filepath.Ext(filename))
		if extension != "" {
			if mimeType, ok := mimeTypes[extension[1:]]; ok {
				w.Header().Add("Content-Type", mimeType)
			}
		}
		w.Write(contents)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("Error reading file: %s\n", err)
	}
}
