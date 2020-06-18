package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/url"
)

type RelayConfig struct {
	actorKey        *rsa.PrivateKey
	domain          *url.URL
	redisClient     *redis.Client
	serviceName     string
	serviceSummary  string
	serviceIconURL  *url.URL
	serviceImageURL *url.URL
}

func NewRelayConfig() (*RelayConfig, error) {
	domain, err := url.ParseRequestURI("https://" + viper.GetString("RELAY_DOMAIN"))
	if err != nil {
		return nil, errors.New("RELAY_DOMAIN: " + err.Error())
	}

	iconURL, err := url.ParseRequestURI(viper.GetString("RELAY_ICON"))
	if err != nil {
		return nil, errors.New("RELAY_ICON: " + err.Error())
	}

	imageURL, err := url.ParseRequestURI(viper.GetString("RELAY_IMAGE"))
	if err != nil {
		return nil, errors.New("RELAY_IMAGE: " + err.Error())
	}

	privateKey, err := readPrivateKeyRSA(viper.GetString("ACTOR_PEM"))
	if err != nil {
		return nil, errors.New("ACTOR_PEM: " + err.Error())
	}

	redisOption, err := redis.ParseURL(viper.GetString("REDIS_URL"))
	if err != nil {
		return nil, errors.New("REDIS_URL: " + err.Error())
	}
	redisClient := redis.NewClient(redisOption)
	err = redisClient.Ping().Err()
	if err != nil {
		return nil, errors.New("Redis Connection Test: " + err.Error())
	}

	return &RelayConfig{
		privateKey,
		domain,
		redisClient,
		viper.GetString("RELAY_SERVICENAME"),
		viper.GetString("RELAY_SUMMARY"),
		iconURL,
		imageURL,
	}, nil
}

func readPrivateKeyRSA(keyPath string) (*rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	decoded, _ := pem.Decode(file)
	privateKey, err := x509.ParsePKCS1PrivateKey(decoded.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
