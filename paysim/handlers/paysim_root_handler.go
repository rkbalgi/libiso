package handlers

import (
	"log"
	"net/http"
	"strings"
)

var docDir = "C:\\Users\\xyz\\git\\paysim-web"

var RootHandler = func(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/paysim/" {
		resourceName := docDir + "\\html\\" + "index.html"
		http.ServeFile(w, r, resourceName)
		return
	}

	if strings.HasSuffix(r.RequestURI, ".html") || strings.HasSuffix(r.RequestURI, ".js") {
		tmp := strings.Split(r.RequestURI, "/")
		lastResourceName := tmp[len(tmp)-1:]
		subDir := "scripts"
		if strings.HasSuffix(r.RequestURI, ".html") {
			subDir = "html"
		}

		resourceName := docDir + "\\" + subDir + "\\" + lastResourceName[0]
		log.Println("serving static resource ", resourceName)
		http.ServeFile(w, r, resourceName)
	}

}
