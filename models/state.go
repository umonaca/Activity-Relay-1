package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

// RelayState is local cache of permanent state information.
// All information are stored redis.
type RelayState struct {
	redisClient    *redis.Client
	relayOptions   relayOptions   `json:"relayConfig"`
	subscriptions  []subscription `json:"subscriptions"`
	limitedDomains []string       `json:"limitedDomains"`
	blockedDomains []string       `json:"blockedDomains"`
}

// RelayOptionsIota is flags for relay options.
type RelayOptionsIota uint

const (
	BlockService     RelayOptionsIota = 1
	ManuallyAccept   RelayOptionsIota = 2
	CreateAsAnnounce RelayOptionsIota = 3
)

type relayOptions struct {
	blockService     bool `json:"blockService"`
	manuallyAccept   bool `json:"manuallyAccept"`
	createAsAnnounce bool `json:"createAsAnnounce"`
}

// DomainConfigTypeIota is flags for limited domain subscribe.
type DomainConfigTypeIota uint

const (
	Limited DomainConfigTypeIota = 1
	Blocked DomainConfigTypeIota = 2
)

type SubscriptionTypeIota uint

const (
	Subscribe SubscriptionTypeIota = 1
	Pending   SubscriptionTypeIota = 2
)

type subscription struct {
	Domain     string `json:"domain,omitempty"`
	InboxURL   string `json:"inbox_url,omitempty"`
	ActivityID string `json:"activity_id,omitempty"`
	ActorID    string `json:"actor_id,omitempty"`
	Available  bool   `json:"available,omitempty"`
}

func (relayState *RelayState) loadConfig() error {
	relayConfigKey := "relay:config"

	blockServiceInt, err := redisHGetOrCreateWithDefault(relayState.redisClient, relayConfigKey, "block_service", "0")
	if err != nil {
		return err
	}
	relayState.relayOptions.blockService = blockServiceInt == "1"

	manuallyAcceptInt, err := redisHGetOrCreateWithDefault(relayState.redisClient, relayConfigKey, "manually_accept", "0")
	if err != nil {
		return err
	}
	relayState.relayOptions.manuallyAccept = manuallyAcceptInt == "1"

	createAsAnnounceInt, err := redisHGetOrCreateWithDefault(relayState.redisClient, relayConfigKey, "create_as_announce", "0")
	if err != nil {
		return err
	}
	relayState.relayOptions.createAsAnnounce = createAsAnnounceInt == "1"

	relayState.limitedDomains, err = relayState.redisClient.HKeys("relay:config:limitedDomain").Result()
	if err != nil {
		return err
	}

	relayState.blockedDomains, err = relayState.redisClient.HKeys("relay:config:blockedDomain").Result()
	if err != nil {
		return err
	}

	relayState.subscriptions, err = relayState.loadSubscription()
	if err != nil {
		return err
	}

	return nil
}

func (relayState *RelayState) loadSubscription() ([]subscription, error) {
	relaySubscriptionKey := "relay:subscription"
	var subscriptions []subscription

	domains, err := relayState.redisClient.Keys(relaySubscriptionKey + ":*").Result()
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		domainName := strings.Replace(domain, relaySubscriptionKey+":", "", 1)
		data, err := relayState.redisClient.HMGet(domain, "inbox_url", "activity_id", "actor_id", "available").Result()
		if err != nil {
			return nil, err
		}
		if inboxURL, ok := data[0].(string); ok {
			activityID := ""
			actorID := ""
			available := true
			if _activityID, ok := data[1].(string); ok {
				activityID = _activityID
			}
			if _actorID, ok := data[2].(string); ok {
				actorID = _actorID
			}
			if _available, ok := data[3].(string); ok {
				if _available != "true" {
					available = false
				}
			}

			subscriptions = append(subscriptions, subscription{
				Domain:     domainName,
				InboxURL:   inboxURL,
				ActivityID: activityID,
				ActorID:    actorID,
				Available:  available,
			})
		}
	}
	return subscriptions, nil
}

func (relayState *RelayState) StateRefreshListener(refreshListener chan<- bool) error {
	messages := relayState.redisClient.Subscribe("relay_refresh").Channel()
	if messages == nil {
		return errors.New("state changes subscribe failed")
	}

	go func() {
		for range messages {
			fmt.Println("config refreshed from state changed")
			err := relayState.loadConfig()
			if err != nil {
				fmt.Println(err.Error())
				if refreshListener != nil {
					refreshListener <- false
				}
			} else {
				if refreshListener != nil {
					refreshListener <- true
				}
			}
		}
	}()

	return nil
}

func (relayState *RelayState) GetOptionValue(iota RelayOptionsIota) (bool, error) {
	switch iota {
	case BlockService:
		return relayState.relayOptions.blockService, nil
	case ManuallyAccept:
		return relayState.relayOptions.manuallyAccept, nil
	case CreateAsAnnounce:
		return relayState.relayOptions.createAsAnnounce, nil
	default:
		return false, errors.New("RelayOption not found")
	}
}

func (relayState *RelayState) SetOptionValue(iota RelayOptionsIota, value bool) error {
	relayConfigKey := "relay:config"
	var modifyField string

	switch iota {
	case BlockService:
		modifyField = "block_service"
	case ManuallyAccept:
		modifyField = "manually_accept"
	case CreateAsAnnounce:
		modifyField = "create_as_announce"
	default:
		return errors.New("RelayOption not found")
	}

	redisValue := "0"
	if value {
		redisValue = "1"
	}
	_, err := relayState.redisClient.HSet(relayConfigKey, modifyField, redisValue).Result()
	if err != nil {
		return err
	}
	_, err = relayState.redisClient.Publish("relay_refresh", nil).Result()
	if err != nil {
		return err
	}
	return nil
}

func (relayState *RelayState) GetDomainConfig(iota DomainConfigTypeIota) ([]string, error) {
	var subscriptions []string

	switch iota {
	case Limited:
		subscriptions = relayState.limitedDomains
	case Blocked:
		subscriptions = relayState.blockedDomains
	default:
		return nil, errors.New("DomainConfigType not found")
	}

	return subscriptions, nil
}

func (relayState *RelayState) SetDomainConfig(iota DomainConfigTypeIota, domain string) error {
	var modifyKey string

	switch iota {
	case Limited:
		modifyKey = "relay:config:limitedDomain"
	case Blocked:
		modifyKey = "relay:config:blockedDomain"
	default:
		return errors.New("DomainConfigType not found")
	}

	_, err := relayState.redisClient.HSet(modifyKey, domain, 1).Result()
	if err != nil {
		return err
	}
	_, err = relayState.redisClient.Publish("relay_refresh", nil).Result()
	if err != nil {
		return err
	}
	return nil
}

func (relayState *RelayState) DelDomainConfig(iota DomainConfigTypeIota, domain string) error {
	var modifyKey string

	switch iota {
	case Limited:
		modifyKey = "relay:config:limitedDomain"
	case Blocked:
		modifyKey = "relay:config:blockedDomain"
	default:
		return errors.New("DomainConfigType not found")
	}

	exist, err := relayState.redisClient.HExists(modifyKey, domain).Result()
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("subscription is not exist")
	}

	_, err = relayState.redisClient.HDel(modifyKey, domain).Result()
	if err != nil {
		return err
	}
	_, err = relayState.redisClient.Publish("relay_refresh", nil).Result()
	if err != nil {
		return err
	}
	return nil
}

func (relayState *RelayState) GetSubscriptions(filter SubscriptionTypeIota) ([]subscription, error) {
	var subscriptions []subscription
	var requestedIsAvailable bool

	switch filter {
	case Subscribe:
		requestedIsAvailable = true
	case Pending:
		requestedIsAvailable = false
	default:
		return nil, errors.New("SubscriptionType not found")
	}

	subscriptionsCount := len(relayState.subscriptions)
	skipCount := 0
	subscriptions = make([]subscription, subscriptionsCount)
	for index, subscription := range relayState.subscriptions {
		if subscription.Available == requestedIsAvailable {
			subscriptions[index-skipCount] = subscription
		} else {
			skipCount++
		}
	}
	subscriptions = subscriptions[:subscriptionsCount-skipCount]

	return subscriptions, nil
}

func (relayState *RelayState) SetSubscription(value subscription) error {
	var available string
	if value.Available {
		available = "true"
	} else {
		available = "false"
	}

	_, err := relayState.redisClient.HMSet("relay:subscription:"+value.Domain, map[string]interface{}{
		"inbox_url":   value.InboxURL,
		"activity_id": value.ActivityID,
		"actor_id":    value.ActorID,
		"available":   available}).Result()
	if err != nil {
		return err
	}
	_, err = relayState.redisClient.Publish("relay_refresh", nil).Result()
	if err != nil {
		return err
	}
	return nil
}

func (relayState *RelayState) DelSubscription(domain string) error {
	exist, err := relayState.redisClient.Exists("relay:subscription:" + domain).Result()
	if err != nil {
		return err
	}
	if exist != 1 {
		return errors.New("subscription is not exist")
	}

	_, err = relayState.redisClient.Del("relay:subscription:" + domain).Result()
	if err != nil {
		return err
	}
	_, err = relayState.redisClient.Publish("relay_refresh", nil).Result()
	if err != nil {
		return err
	}
	return nil
}

func (relayState *RelayState) PromoteSubscription(domain string) error {
	pendings, err := relayState.GetSubscriptions(Pending)
	if err != nil {
		return err
	}

	for _, pending := range pendings {
		if domain == pending.Domain {
			pending.Available = true

			err := relayState.SetSubscription(pending)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("subscription is not exist")
}

func NewRelayState(relayConfig *RelayConfig) (*RelayState, error) {
	var relayState RelayState
	relayState.redisClient = relayConfig.redisClient
	err := relayState.loadConfig()
	if err != nil {
		return nil, err
	}

	return &relayState, nil
}
