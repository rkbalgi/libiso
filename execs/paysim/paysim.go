package main

import (
	"github.com/rkbalgi/go/paysim/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/paysim/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test: Paysim v1.00"))
	})

	http.HandleFunc("/", handlers.RootHandler)
	http.Handle("/paysim/get_layout", new(handlers.PaysimDefaultHandler))
	http.Handle("/paysim/parse_trace", new(handlers.ParseTraceHandlerHandler))
	http.Handle("/paysim/send_message", new(handlers.SendMessageHandlerHandler))
	

	//start and listen at 8080
	http.ListenAndServe(":8080", nil)

}
