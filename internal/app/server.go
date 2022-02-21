package app

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var hashTable = make(map[string]string)
var serverProtocol = "http"
var serverAddress = "localhost"
var serverPort = "8080"

func StartServer() {
	http.HandleFunc("/", Shortener)
	http.ListenAndServe(serverAddress+":"+serverPort, nil)
}

func Shortener(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		if r.URL.Path != "/" {
			log.Println("Bad Request: " + r.URL.Path)
			http.Error(w, "Bad Request. Only / allowed.", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error at io.ReadAll: " + err.Error())
			http.Error(w, "Error at io.ReadAll: "+err.Error(), http.StatusBadRequest)
			return
		}

		uri := string(body)
		uri = strings.TrimSpace(uri)
		if len(uri) == 0 {
			http.Error(w, "Empty body in request", http.StatusBadRequest)
			return
		}

		_, err = url.ParseRequestURI(uri)
		if err != nil {
			http.Error(w, "Illegal URL in request", http.StatusBadRequest)
			return
		}

		hash := GetHash(uri)

		// индекс есть в таблице
		if _, ok := hashTable[hash]; ok {
			// а строки разные - катастрофа
			if uri != hashTable[hash] {
				err := "Double hash " + hash + " for strings: <" + uri + "> & <" + hashTable[hash] + ">"
				log.Println(err)
				http.Error(w, err, http.StatusInternalServerError)
				return
			}
		} else {
			hashTable[hash] = uri
		}
		shortURL := serverProtocol + "://" + serverAddress + ":" + serverPort + "/" + hash
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))

	case http.MethodGet:
		if r.URL.Path == "/" {
			http.Error(w, "Bar Request", http.StatusBadRequest)
			return
		}
		hash := r.URL.Path[1:]
		if longURL, ok := hashTable[hash]; ok {
			w.Header().Set("Location", longURL)
			w.WriteHeader(http.StatusTemporaryRedirect)

			//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			//w.Write([]byte(longURL))
			return
		} else {
			http.Error(w, "Long URL for <"+hash+"> not found", http.StatusNotFound)
		}

	default:
		http.Error(w, "Bar Request", http.StatusBadRequest)

	}
}
