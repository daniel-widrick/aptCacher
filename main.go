package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {

	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("---------")
		log.Printf("Host %v", r.Host)
		log.Printf("Path %v", r.URL.Path)
		log.Printf("Method %v", r.Method)
		log.Printf("Request: %v", r.Header)

		localPath := path.Join("./cache", r.Host, r.URL.Path)
		localFile := path.Join(localPath, "file")
		log.Printf("Checking for local file: %v", localPath)

		//Check for local cache path and create if it doesn't exist
		_, err := os.Stat(localPath)
		if err != nil {
			err := os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				msg := fmt.Sprintf("Unable to create path: %s, %v", localPath, err)
				log.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		//Check that the local file exits and compare the requested modified time to file time
		log.Printf("Checking for file: %s", localFile)
		fileInfo, err := os.Stat(localFile)
		var doFetch = false
		if err != nil || fileInfo == nil {
			//File doesn't exist? Fetch
			doFetch = true
		} else {
			//Get modified time from request header
			log.Printf("Checking cache expiration: %s", localFile)
			modTimeRaw := r.Header.Get("If-Modified-Since")
			log.Printf("modified header: %s", modTimeRaw)
			if modTimeRaw != "" {
				modTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", modTimeRaw)
				log.Printf("Parsed date: %v", modTime)
				if err != nil {
					msg := fmt.Sprintf("Error parsing date: %s. %v", modTimeRaw, err)
					log.Printf(msg)
					doFetch = true
				} else {
					log.Printf("modtime: %v", modTime)
					log.Printf("fileinfo: %v", fileInfo)
					if modTime.After(fileInfo.ModTime()) {
						log.Printf("File: %s expired", localFile)
						doFetch = true
					}
				}
			}
		}

		log.Printf("Do fetch: %v", doFetch)
		if doFetch {
			log.Printf("Fetching...")
			fetchUrl := fmt.Sprintf("https://%s/%s", r.Host, r.URL.Path)
			resp, err := http.Get(fetchUrl)
			if err != nil {
				msg := fmt.Sprintf("Unable to fetch %s\n%v", fetchUrl, err)
				log.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg := fmt.Sprintf("Error reading response: %v", err)
				log.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			//Write file
			f, err := os.Create(localFile)
			if err != nil {
				log.Printf("Error creating file: %s. %v", localFile, err)
			} else {
				log.Printf("Writing file %s.", localFile)
				_, err := f.Write(body)
				if err != nil {
					log.Printf("Error writing file: %s. %v", localFile, err)
				}
			}
			_, err = w.Write(body)
			if err != nil {
				log.Printf("Error writing response: %s. %v", localFile, err)
				return
			}
		} else {
			log.Printf("Serving file from cache: %s", localFile)
			f, err := os.Open(localFile)
			if err != nil {
				msg := fmt.Sprintf("Error opening file: %s. %v", localFile, err)
				log.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			body, err := io.ReadAll(f)
			if err != nil {
				msg := fmt.Sprintf("Error reading file: %s. %v", localFile, err)
				log.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)
			if err != nil {
				log.Printf("Error sending response: %v", err)
			}
		}
	})

	log.Fatalf("%v", http.ListenAndServe("0.0.0.0:8700", s))
}
