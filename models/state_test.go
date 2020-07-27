package models

import (
	"strconv"
	"testing"
)

func TestNewRelayState(t *testing.T) {
	t.Run("success valid state with configuration", func(t *testing.T) {
		relayConfig, err := NewRelayConfig()
		if err != nil {
			t.Fatal(err)
		}

		_, err = NewRelayState(relayConfig)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func createRelayState(t *testing.T, relayConfig *RelayConfig) *RelayState {
	if relayConfig == nil {
		relayConfig = createRelayConfig(t)
	}

	relayState, err := NewRelayState(relayConfig)
	if err != nil {
		t.Fatal(err)
	}
	return relayState
}

func TestRelayState_GetOptionValue(t *testing.T) {
	t.Run("success get initial option values", func(t *testing.T) {
		var optionLabels = []RelayOptionsIota{
			BlockService,
			ManuallyAccept,
			CreateAsAnnounce,
		}
		relayState := createRelayState(t, nil)

		for _, optionLabel := range optionLabels {
			value, err := relayState.GetOptionValue(optionLabel)
			if err != nil {
				t.Error(err)
			}
			if value != false {
				t.Error("Failed get options value correctly:" + strconv.Itoa(int(optionLabel)))
			}
		}
	})

	t.Run("return error for unknown option label", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		_, err := relayState.GetOptionValue(12345)
		if err == nil {
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func TestRelayState_SetOptionValue(t *testing.T) {
	t.Run("success set and get option values", func(t *testing.T) {
		var optionLabels = []RelayOptionsIota{
			BlockService,
			ManuallyAccept,
			CreateAsAnnounce,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, optionLabel := range optionLabels {
			err := relayState.SetOptionValue(optionLabel, true)
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			value, err := relayState.GetOptionValue(optionLabel)
			if err != nil {
				t.Error(err)
			}

			if value != true {
				t.Error("Failed get options value correctly:" + strconv.Itoa(int(optionLabel)))
			}

			err = relayState.SetOptionValue(optionLabel, false)
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}
		}
	})

	t.Run("return error for unknown option label", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		err := relayState.SetOptionValue(12345, true)
		if err == nil {
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func createDemoSubscription() subscription {
	return subscription{
		Domain:     "example.yukimochi.dev",
		InboxURL:   "https://example.yukimochi.dev/inbox",
		ActivityID: "https://example.yukimochi.dev/be0802b1-8648-4598-b794-2ed19532100d",
		ActorID:    "https://example.yukimochi.dev/actor",
		Available:  true,
	}
}

func TestRelayState_GetDomainConfig(t *testing.T) {
	t.Run("return error for unknown option label", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		_, err = relayState.GetDomainConfig(12345)
		if err == nil {
			<-refreshListener
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func TestRelayState_SetDomainConfig(t *testing.T) {
	t.Run("success set and get limited and blacked domain", func(t *testing.T) {
		var labels = []DomainConfigTypeIota{
			Limited,
			Blocked,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, label := range labels {
			err := relayState.SetDomainConfig(label, "example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			value, err := relayState.GetDomainConfig(label)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, domain := range value {
				if domain == "example.yukimochi.dev" {
					contain = true
				}
			}

			if contain != true {
				t.Error("Failed set or get limited, blacked domain: " + strconv.Itoa(int(label)))
			}

			err = relayState.DelDomainConfig(label, "example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}
		}
	})

	t.Run("return error for unknown option label", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		err = relayState.SetDomainConfig(12345, "example.yukimochi.dev")
		if err == nil {
			<-refreshListener
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func TestRelayState_DelDomainConfig(t *testing.T) {
	t.Run("success del limited and blacked domain", func(t *testing.T) {
		var labels = []DomainConfigTypeIota{
			Limited,
			Blocked,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, label := range labels {
			err := relayState.SetDomainConfig(label, "example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			err = relayState.DelDomainConfig(label, "example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			value, err := relayState.GetDomainConfig(label)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, domain := range value {
				if domain == "example.yukimochi.dev" {
					contain = true
				}
			}

			if contain != false {
				t.Error("Failed del limited, blacked domain: " + strconv.Itoa(int(label)))
			}
		}
	})

	t.Run("return error for unknown domain", func(t *testing.T) {
		var labels = []DomainConfigTypeIota{
			Limited,
			Blocked,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, label := range labels {
			err = relayState.DelDomainConfig(label, "example.yukimochi.dev")
			if err == nil {
				<-refreshListener
				t.Error("Failed return error, given unknown domain")
			}
		}
	})

	t.Run("return error for unknown option label", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		err = relayState.DelDomainConfig(12345, "example.yukimochi.dev")
		if err == nil {
			<-refreshListener
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func TestRelayState_GetSubscriptions(t *testing.T) {
	t.Run("return error for unknown SubscriptionType", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		_, err = relayState.GetSubscriptions(12345)
		if err == nil {
			t.Error("Failed return error, given unknown SubscriptionType")
		}
	})
}

func TestRelayState_SetSubscription(t *testing.T) {
	t.Run("success set and get subscribe, pending", func(t *testing.T) {
		var labels = []SubscriptionTypeIota{
			Subscribe,
			Pending,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, label := range labels {
			srcSubscription := createDemoSubscription()
			if label == Pending {
				srcSubscription.Available = false
			}

			err := relayState.SetSubscription(srcSubscription)
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			value, err := relayState.GetSubscriptions(label)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, subscription := range value {
				if subscription.Domain == srcSubscription.Domain {
					contain = true
				}
			}

			if contain != true {
				t.Error("Failed set or get subscribe, pending: " + strconv.Itoa(int(label)))
			}

			err = relayState.DelSubscription("example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}
		}
	})
}

func TestRelayState_DelSubscription(t *testing.T) {
	t.Run("success del subscribe, pending", func(t *testing.T) {
		var labels = []SubscriptionTypeIota{
			Subscribe,
			Pending,
		}
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		for _, label := range labels {
			srcSubscription := createDemoSubscription()
			if label == Pending {
				srcSubscription.Available = false
			}

			err := relayState.SetSubscription(srcSubscription)
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			err = relayState.DelSubscription("example.yukimochi.dev")
			if err != nil {
				t.Error(err)
			} else {
				<-refreshListener
			}

			value, err := relayState.GetSubscriptions(label)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, subscription := range value {
				if subscription.Domain == srcSubscription.Domain {
					contain = true
				}
			}

			if contain != false {
				t.Error("Failed del subscribe, pending: " + strconv.Itoa(int(label)))
			}
		}
	})

	t.Run("return error for unknown domain", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		err = relayState.DelSubscription("example.yukimochi.dev")
		if err == nil {
			<-refreshListener
			t.Error("Failed return error, given unknown option label")
		}
	})
}

func TestRelayState_PromoteSubscription(t *testing.T) {
	t.Run("success promote pending to subscribe", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		srcSubscription := createDemoSubscription()
		srcSubscription.Available = false

		err = relayState.SetSubscription(srcSubscription)
		if err != nil {
			t.Error(err)
		} else {
			<-refreshListener
		}

		{
			value, err := relayState.GetSubscriptions(Pending)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, subscription := range value {
				if subscription.Domain == srcSubscription.Domain {
					contain = true
				}
			}

			if contain != true {
				t.Error("Failed set or get pending")
			}
		}

		err = relayState.PromoteSubscription("example.yukimochi.dev")
		if err != nil {
			t.Error(err)
		} else {
			<-refreshListener
		}

		{
			value, err := relayState.GetSubscriptions(Subscribe)
			if err != nil {
				t.Error(err)
			}

			contain := false
			for _, subscription := range value {
				if subscription.Domain == srcSubscription.Domain {
					contain = true
				}
			}

			if contain != true {
				t.Error("Failed promote pending to subscribe")
			}
		}
	})

	t.Run("return error for unknown domain", func(t *testing.T) {
		relayState := createRelayState(t, nil)

		refreshListener := make(chan bool)
		err := relayState.StateRefreshListener(refreshListener)
		if err != nil {
			t.Error(err)
		}

		err = relayState.PromoteSubscription("example.yukimochi.dev")
		if err == nil {
			<-refreshListener
			t.Error("Failed return error, given unknown domain")
		}
	})
}
