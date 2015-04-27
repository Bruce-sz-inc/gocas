package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/protocol/cas"
	"github.com/apognu/gocas/protocol/oauth"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
)

var (
	c = flag.String("config", "/etc/gocas.yaml", "path to GoCAS configuration file")
)

type Protocol func(*mux.Router)

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", util.Url("/login"))
	w.WriteHeader(http.StatusFound)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	config.Set(*c)
	var protocols = map[string]Protocol{
		"cas":   cas.New,
		"oauth": oauth.New,
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", redirect)
	sr := r
	if config.Get().UrlPrefix != "" {
		sr = r.PathPrefix(config.Get().UrlPrefix).Subrouter()
		sr.HandleFunc("/", redirect)
	}
	sr.HandleFunc("/validate", validateHandler).Methods("GET")
	sr.HandleFunc("/serviceValidate", serviceValidateHandler).Methods("GET")
	sr.HandleFunc("/logout", logoutHandler).Methods("GET")

	protocols[config.Get().Protocol](sr)

	logrus.Infof("started gocas CAS server, %s", time.Now())
	http.ListenAndServe("0.0.0.0:8080", r)
}
