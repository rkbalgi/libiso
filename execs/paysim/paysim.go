package main

import (
	"github.com/rkbalgi/go/paysim/handlers"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/paysim/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test: Paysim v1.00"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/paysim/" {
			resource_name := "c:\\users\\132968\\Desktop\\ang-js\\index.html"
			http.ServeFile(w, r, resource_name)
			return;
		}

		if strings.HasSuffix(r.RequestURI, ".html") || strings.HasSuffix(r.RequestURI, ".js") {
			tmp := strings.Split(r.RequestURI, "/")
			last_resource_name := tmp[len(tmp)-1:]
			resource_name := "c:\\users\\132968\\Desktop\\ang-js\\" + last_resource_name[0]
			log.Println("serving static resource ", resource_name)
			http.ServeFile(w, r, resource_name)
		}

	})

	http.Handle("/paysim/get_layout", new(handlers.PaysimDefaultHandler))
	http.ListenAndServe(":8080", nil)

}
