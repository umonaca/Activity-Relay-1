package models

import (
	"strings"
	"testing"
)

func TestNewActivityPubActorFromSelfKey(t *testing.T) {
	relayConfig := createRelayConfig(t)

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

func TestNewWebfingerFromActor(t *testing.T) {
	relayConfig := createRelayConfig(t)
	activityPubActor := NewActivityPubActorFromSelfKey(relayConfig)

	webfinger := NewWebfingerFromActor(activityPubActor, relayConfig.domain.Host)
	if webfinger.Subject != "acct:relay@relay.toot.yukimochi.jp" {
		t.Error("Failed parse: webfinger.Subject")
	}
	if len(webfinger.Links) != 1 {
		t.Error("Failed parse: webfinger.Links")
	} else {
		if webfinger.Links[0].Rel != "self" {
			t.Error("Failed parse: webfinger.Links[0].Rel")
		}
		if webfinger.Links[0].Type != "application/activity+json" {
			t.Error("Failed parse: webfinger.Links[0].Type")
		}
		if webfinger.Links[0].Href != activityPubActor.ID {
			t.Error("Failed parse: webfinger.Links[0].Href")
		}
	}
}

func TestNewWellKnownNodeinfo(t *testing.T) {
	relayConfig := createRelayConfig(t)

	wellknownNodeinfo := NewWellKnownNodeinfo(relayConfig.domain.Host)
	if len(wellknownNodeinfo.Links) != 1 {
		t.Error("Failed parse: wellknownNodeinfo.Links")
	} else {
		if wellknownNodeinfo.Links[0].Rel != "http://nodeinfo.diaspora.software/ns/schema/2.1" {
			t.Error("Failed parse: wellknownNodeinfo.Links[0].Rel")
		}
		if wellknownNodeinfo.Links[0].Href != "https://relay.toot.yukimochi.jp/nodeinfo/2.1" {
			t.Error("Failed parse: wellknownNodeinfo.Links[0].Href")
		}
	}
}
