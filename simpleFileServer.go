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
	filename := r.RequestURI[1:]
	filePath := rootPath + string(os.PathSeparator) + filename
	if contents, err := os.ReadFile(filePath); err == nil {
		extension := strings.ToLower(filepath.Ext(filename))
		if extension != "" {
			if mimeType, ok := mimeTypes[extension[1:]]; ok {
				w.Header().Add("Content-Type", mimeType)
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "*://localhost:*/*")
		w.Write(contents)
	} else {
		fmt.Printf("Error reading file: %s\n", err)
	}
}
