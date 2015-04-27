package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
)

var (
	c = flag.String("config", "/etc/gocas.yaml", "path to GoCAS configuration file")
)

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", util.Url("/login"))
	w.WriteHeader(http.StatusFound)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	config.Set(*c)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", redirect).Methods("GET")

	prefix := config.Get().UrlPrefix
	sr := r
	if prefix != "" {
		sr = r.PathPrefix(prefix).Subrouter()
		sr.HandleFunc("/", redirect)
	}
	sr.HandleFunc("/login", loginRequestorHandler).Methods("GET")
	sr.HandleFunc("/login", loginAcceptorHandler).Methods("POST")
	sr.HandleFunc("/validate", validateHandler).Methods("GET")
	sr.HandleFunc("/serviceValidate", serviceValidateHandler).Methods("GET")
	sr.HandleFunc("/logout", logoutHandler).Methods("GET")

	logrus.Infof("started gocas CAS server, %s", time.Now())
	http.ListenAndServe("0.0.0.0:8080", r)
}
