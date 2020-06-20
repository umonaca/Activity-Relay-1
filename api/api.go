package api

import (
	"Activity-Relay/models"
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"net/http"
	"net/url"
)

var (
	machineryServer *machinery.Server
	selfActor       models.ActivityPubActor
	selfHostname    *url.URL
)

func Entrypoint(globalConfig *models.RelayConfig) error {
	var err error

	machineryServer, err = models.NewMachineryServer(globalConfig)
	if err != nil {
		return err
	}
	selfActor = models.NewActivityPubActorFromSelfKey(globalConfig)
	selfHostname = globalConfig.ServerHostname()

	registResourceHandlers()

	fmt.Println("Staring API Server at", globalConfig.ServerBind())
	err = http.ListenAndServe(globalConfig.ServerBind(), nil)
	if err != nil {
		return err
	}

	return nil
}

func registResourceHandlers() {
	http.HandleFunc("/.well-known/nodeinfo", handleResourceRequest)
	http.HandleFunc("/.well-known/webfinger", handleResourceRequest)
	http.HandleFunc("/actor", handleResourceRequest)
	http.HandleFunc("/nodeinfo/2.1", handleResourceRequest)
}
