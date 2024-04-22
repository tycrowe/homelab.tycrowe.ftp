package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Handle incoming FTP requests
	// Start by checking the request method
	if r.Method == "POST" {
		// Parse the request, this will allow us to get the file from the request
		// The 10 << 20 is the maximum file size we can accept, in this case 10MB.
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing request, too large probably!", http.StatusBadRequest)
			return
		}

		// Handle POST requests
		// Get the file from the request
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error getting file from request", http.StatusBadRequest)
			return
		}

		// Close the file when we are done
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				http.Error(w, "Error closing file", http.StatusInternalServerError)
			}
		}(file) // This is a closure, it will close the file when the function is done

		dst, err := os.Create(handler.Filename)

		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}

		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				http.Error(w, "Error closing file", http.StatusInternalServerError)
				return
			}
		}(dst)

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}

		// Send a response to the client
		_, err = w.Write([]byte("File uploaded successfully!"))
		if err != nil {
			return
		}
	}
}

// main
// This function starts the server and listens for incoming requests.
func main() {
	// Start an FTP server
	http.HandleFunc("/upload", handleRequest)
	httpError := http.ListenAndServe(":8080", nil)
	if httpError != nil {
		panic(httpError)
	}
}
