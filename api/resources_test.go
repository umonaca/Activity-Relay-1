package api

import (
	"Activity-Relay/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func createSelfObjects(t *testing.T) (models.ActivityPubActor, *url.URL) {
	relayConfig, err := models.NewRelayConfig()
	if err != nil {
		t.Fatal(err)
	}
	selfActor := models.NewActivityPubActorFromSelfKey(relayConfig)

	return selfActor, relayConfig.ServerHostname()
}

func Test_handleResourceRequest(t *testing.T) {
	selfActor, selfHostname = createSelfObjects(t)

	s := httptest.NewServer(http.HandlerFunc(handleResourceRequest))
	defer s.Close()

	client := new(http.Client)

	t.Run("return Bad Gateway", func(t *testing.T) {
		invalidMethods := []string{
			"POST", "PUT", "PATCH", "OPTION",
		}
		requestPath := []string{
			"/.well-known/nodeinfo", "/.well-known/webfinger", "/actor", "/nodeinfo/2.1",
		}

		for _, invalidMethod := range invalidMethods {
			for _, path := range requestPath {
				req, _ := http.NewRequest(invalidMethod, s.URL+path, nil)

				r, err := client.Do(req)
				if err != nil {
					t.Error("Request error:", err.Error())
				}
				if r.StatusCode != 400 {
					t.Error("Response status not Bad Gateway")
				}
			}
		}
	})

	t.Run("success /.well-known/nodeinfo", func(t *testing.T) {
		methods := []string{"GET", "HEAD"}

		for _, method := range methods {
			req, _ := http.NewRequest(method, s.URL+"/.well-known/nodeinfo", nil)

			r, err := client.Do(req)
			if err != nil {
				t.Error("Request error:", err.Error())
			}
			defer r.Body.Close()

			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("Response header missing: Content-Type")
			}
			if r.StatusCode != 200 {
				t.Error("Response status not OK")
			}

			if method == "GET" {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
				var wellKnownNodeinfo models.WellKnownNodeinfo
				err = json.Unmarshal(data, &wellKnownNodeinfo)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
			}
		}
	})

	t.Run("success /.well-known/webfinger", func(t *testing.T) {
		methods := []string{"GET", "HEAD"}

		for _, method := range methods {
			req, _ := http.NewRequest(method, s.URL+"/.well-known/webfinger", nil)
			q := req.URL.Query()
			q.Add("resource", "acct:relay@"+selfHostname.Host)
			req.URL.RawQuery = q.Encode()

			r, err := client.Do(req)
			if err != nil {
				t.Error("Request error:", err.Error())
			}
			defer r.Body.Close()

			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("Response header missing: Content-Type")
			}
			if r.StatusCode != 200 {
				t.Error("Response status not OK")
			}

			if method == "GET" {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
				var webfinger models.Webfinger
				err = json.Unmarshal(data, &webfinger)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
			}
		}
	})

	t.Run("return Bad Gateway /.well-known/webfinger", func(t *testing.T) {
		methods := []string{"GET", "HEAD"}

		for _, method := range methods {
			req, _ := http.NewRequest(method, s.URL+"/.well-known/webfinger", nil)

			r, err := client.Do(req)
			if err != nil {
				t.Error("Request error:", err.Error())
			}
			defer r.Body.Close()

			if r.StatusCode != 400 {
				t.Error("Response status not Bad Gateway")
			}
		}
	})

	t.Run("return Not Found /.well-known/webfinger", func(t *testing.T) {
		methods := []string{"GET", "HEAD"}

		for _, method := range methods {
			req, _ := http.NewRequest(method, s.URL+"/.well-known/webfinger", nil)
			q := req.URL.Query()
			q.Add("resource", "acct:nil@"+selfHostname.Host)
			req.URL.RawQuery = q.Encode()

			r, err := client.Do(req)
			if err != nil {
				t.Error("Request error:", err.Error())
			}
			defer r.Body.Close()

			if r.StatusCode != 404 {
				t.Error("Response status not Not Found")
			}
		}
	})

	t.Run("success /actor", func(t *testing.T) {
		methods := []string{"GET", "HEAD"}

		for _, method := range methods {
			req, _ := http.NewRequest(method, s.URL+"/actor", nil)

			r, err := client.Do(req)
			if err != nil {
				t.Error("Request error:", err.Error())
			}
			defer r.Body.Close()

			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("Response header missing: Content-Type")
			}
			if r.StatusCode != 200 {
				t.Error("Response status not OK")
			}

			if method == "GET" {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
				var activityPubActor models.ActivityPubActor
				err = json.Unmarshal(data, &activityPubActor)
				if err != nil {
					t.Error("Response body error:" + err.Error())
				}
			}
		}
	})
}
