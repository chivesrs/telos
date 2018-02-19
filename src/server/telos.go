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

// Parameters are the parameters used to render the HTML template.
type Parameters struct {
	Rate  int64
	Error error
}

func main() {
	flag.Parse()
	if *certFile == "" || *keyFile == "" {
		log.Fatal("Both --cert_file and --key_file must be set")
	}
	if *templateFile == "" {
		log.Fatal("--template_file must be set")
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
		params := Parameters{
			Error: fmt.Errorf("unable to parse enrage: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		t.Execute(w, params)
		return
	}

	streak, err := strconv.ParseInt(values.Get("streak"), 10, 64)
	if err != nil {
		params := Parameters{
			Error: fmt.Errorf("unable to parse streak: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		t.Execute(w, params)
		return
	}

	lotd := values.Get("lotd") == "on"
	rate, err := droprate.DropRate(enrage, streak, lotd)
	if err != nil {
		params := Parameters{
			Error: fmt.Errorf("unable to calculate drop rate: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		t.Execute(w, params)
		return
	}

	params := Parameters{
		Rate: rate,
	}
	w.WriteHeader(http.StatusOK)
	t.Execute(w, params)
}
