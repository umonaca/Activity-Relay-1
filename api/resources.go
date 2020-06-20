package api

import (
	"Activity-Relay/models"
	"encoding/json"
	"net/http"
)

func handleResourceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		w.Write(nil)

		return
	}

	var body []byte
	var err error

	switch r.URL.Path {
	case "/.well-known/nodeinfo":
		data := models.NewWellKnownNodeinfo(selfHostname.Host)
		body, err = json.Marshal(&data)
	case "/.well-known/webfinger":
		res := r.URL.Query()["resource"]
		if len(res) == 0 {
			w.WriteHeader(400)
			w.Write(nil)

			return
		}

		data := models.NewWebfingerFromActor(selfActor, selfHostname.Host)
		if data.Subject != res[0] {
			w.WriteHeader(404)
			w.Write(nil)

			return
		}

		body, err = json.Marshal(&data)
	case "/actor":
		body, err = json.Marshal(&selfActor)
	case "/nodeinfo/2.1":
		// TODO: Implement nodeinfo 2.1
		body = []byte("nodeinfo 2.1")
	}

	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}
