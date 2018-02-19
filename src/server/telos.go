package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"droprate"
)

var (
	certFile     = flag.String("cert_file", "", "Full path to cert.pem")
	keyFile      = flag.String("key_file", "", "Full path to privkey.pem")
	templateFile = flag.String("template_file", "", "Full path to template.html")
)

var t *template.Template

func main() {
	flag.Parse()
	if *certFile == "" || *keyFile == "" {
		log.Fatalf("Both --cert_file and --key_file must be set")
	}
	t = template.Must(template.ParseFiles(*templateFile))

	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":8080", *certFile, *keyFile, nil)
	log.Fatalf("Unable to listen: %v\n", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch req.Method {
	case "", "GET":
		if err := t.Execute(w, nil); err != nil {
			log.Printf("Unable to send form template: %v\n", err)
		}
	case "POST":
		req.ParseForm()
		handleForm(req.Form, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		output := "This server only supports GET and POST requests."
		w.Write([]byte(output))
	}
}

func handleForm(values url.Values, w http.ResponseWriter) {
	enrage, err := strconv.ParseInt(values.Get("enrage"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to parse enrage: %v", err)))
		return
	}
	streak, err := strconv.ParseInt(values.Get("streak"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to parse streak: %v", err)))
		return
	}
	lotd := values.Get("lotd") == "on"
	rate, err := droprate.DropRate(enrage, streak, lotd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to calculate drop rate: %v", err)))
		return
	}
	t.Execute(w, rate)
}
