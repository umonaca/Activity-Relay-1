package models

import (
	"strings"
	"testing"
)

func TestNewActivityPubActorFromSelfKey(t *testing.T) {
	relayConfig, err := NewRelayConfig()
	if err != nil {
		t.Fatal(err)
	}

	activityPubActor := NewActivityPubActorFromSelfKey(relayConfig)
	if activityPubActor.Context.([]string)[0] != "https://www.w3.org/ns/activitystreams" {
		t.Error("Failed parse: activityPubActor.Context[0]")
	}
	if activityPubActor.Context.([]string)[1] != "https://w3id.org/security/v1" {
		t.Error("Failed parse: activityPubActor.Context[1]")
	}

	validPublicKey := generatePublicKeyPEMString(&relayConfig.actorKey.PublicKey)
	if !strings.Contains(validPublicKey, "-----BEGIN RSA PUBLIC KEY-----") {
		t.Error("Failed in generatePublicKeyPEMString().")
	}
	if activityPubActor.PublicKey.PublicKeyPem != validPublicKey {
		t.Error("Failed parse: activityPubActor.PublicKey.PublicKeyPem")
	}
}
