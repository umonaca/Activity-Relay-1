/*
Models provide struct for config and type definition for ActivityPub, Nodeinfo, Webfinger.
*/
package models

// ActivityPubActor is compacted ActivityPub actor resource.
// reference: https://www.w3.org/TR/activitypub/#actor-objects
type ActivityPubActor struct {
	Context           interface{} `json:"@context,omitempty"`
	ID                string      `json:"id,omitempty"`
	Type              string      `json:"type,omitempty"`
	Name              string      `json:"name,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	Inbox             string      `json:"inbox,omitempty"`
	Endpoints         *struct {
		SharedInbox string `json:"sharedInbox,omitempty"`
	} `json:"endpoints,omitempty"`
	PublicKey struct {
		ID           string `json:"id,omitempty"`
		Owner        string `json:"owner,omitempty"`
		PublicKeyPem string `json:"publicKeyPem,omitempty"`
	} `json:"publicKey,omitempty"`
	Icon struct {
		URL string `json:"url,omitempty"`
	} `json:"icon,omitempty"`
	Image struct {
		URL string `json:"url,omitempty"`
	} `json:"image,omitempty"`
}

// NewActivityPubActorFromSelfKey create relay server's self ActivityPub actor object.
func NewActivityPubActorFromSelfKey(globalConfig *RelayConfig) ActivityPubActor {
	hostname := globalConfig.domain.String()
	publicKey := &globalConfig.actorKey.PublicKey
	publicKeyPemString := generatePublicKeyPEMString(publicKey)

	newActor := ActivityPubActor{
		Context:           []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
		ID:                hostname + "/actor",
		Type:              "Service",
		Name:              globalConfig.serviceName,
		PreferredUsername: "relay",
		Summary:           globalConfig.serviceSummary,
		Inbox:             hostname + "/inbox",
		PublicKey: struct {
			ID           string `json:"id,omitempty"`
			Owner        string `json:"owner,omitempty"`
			PublicKeyPem string `json:"publicKeyPem,omitempty"`
		}{
			ID:           hostname + "/actor#main-key",
			Owner:        hostname + "/actor",
			PublicKeyPem: publicKeyPemString,
		},
		Icon: struct {
			URL string `json:"url,omitempty"`
		}{
			URL: globalConfig.serviceIconURL.String(),
		},
		Image: struct {
			URL string `json:"url,omitempty"`
		}{
			URL: globalConfig.serviceImageURL.String(),
		},
	}

	return newActor
}

// Webfinger is webfinger response resource.
// reference: https://tools.ietf.org/html/rfc7033
type Webfinger struct {
	Subject string `json:"subject,omitempty"`
	Links   []struct {
		Rel  string `json:"rel,omitempty"`
		Type string `json:"type,omitempty"`
		Href string `json:"href,omitempty"`
	} `json:"links,omitempty"`
}

// Nodeinfo is server information about distributed social networks.
// reference: http://nodeinfo.diaspora.software
type Nodeinfo struct {
	Version  string `json:"version"`
	Software struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Repository string `json:"repository,omitempty"`
	} `json:"software"`
	Protocols []string `json:"protocols"`
	Services  struct {
		Inbound  []string `json:"inbound"`
		Outbound []string `json:"outbound"`
	} `json:"services"`
	OpenRegistrations bool `json:"openRegistrations"`
	Usage             struct {
		Users struct {
			Total          int `json:"total"`
			ActiveMonth    int `json:"activeMonth"`
			ActiveHalfyear int `json:"activeHalfyear"`
		} `json:"users"`
	} `json:"usage"`
	Metadata struct {
	} `json:"metadata"`
}
