package handlers

import (
	"log"
	"net/http"
	"strings"
)



var doc_dir = "C:\\Users\\132968\\git\\paysim-web"

var RootHandler = func(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/paysim/" {
		resource_name := doc_dir + "\\html\\" + "index.html"
		http.ServeFile(w, r, resource_name)
		return
	}

	if strings.HasSuffix(r.RequestURI, ".html") || strings.HasSuffix(r.RequestURI, ".js") {
		tmp := strings.Split(r.RequestURI, "/")
		last_resource_name := tmp[len(tmp)-1:]
		sub_dir := "scripts"
		if strings.HasSuffix(r.RequestURI, ".html") {
			sub_dir = "html"
		}

		resource_name := doc_dir + "\\" + sub_dir + "\\" + last_resource_name[0]
		log.Println("serving static resource ", resource_name)
		http.ServeFile(w, r, resource_name)
	}

}
