package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/url"
)

// RelayConfig contains valid configuration.
type RelayConfig struct {
	actorKey        *rsa.PrivateKey
	domain          *url.URL
	redisClient     *redis.Client
	redisURL        string
	serverBind      string
	serviceName     string
	serviceSummary  string
	serviceIconURL  *url.URL
	serviceImageURL *url.URL
}

// NewRelayConfig create valid RelayConfig from viper configuration. If invalid configuration detected, return error.
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

	redisURL := viper.GetString("REDIS_URL")
	redisOption, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.New("REDIS_URL: " + err.Error())
	}
	redisClient := redis.NewClient(redisOption)
	err = redisClient.Ping().Err()
	if err != nil {
		return nil, errors.New("Redis Connection Test: " + err.Error())
	}

	serverBind := viper.GetString("RELAY_BIND")

	return &RelayConfig{
		actorKey:        privateKey,
		domain:          domain,
		redisClient:     redisClient,
		redisURL:        redisURL,
		serverBind:      serverBind,
		serviceName:     viper.GetString("RELAY_SERVICENAME"),
		serviceSummary:  viper.GetString("RELAY_SUMMARY"),
		serviceIconURL:  iconURL,
		serviceImageURL: imageURL,
	}, nil
}

// ServerBind: API Server's bind interface definition.
func (relayConfig RelayConfig) ServerBind() string {
	return relayConfig.serverBind
}

// ServerBind: API Server's bind interface definition.
func (relayConfig RelayConfig) ServerHostname() *url.URL {
	return relayConfig.domain
}

func (relayConfig RelayConfig) DumpWelcomeMessage(moduleName string) string {
	return fmt.Sprintf(`Welcome to YUKIMOCHI Activity-Relay [Project-Improve] - %s
 - Configuration
RELAY NAME   : %s
RELAY DOMAIN : %s
REDIS URL    : %s
BIND ADDRESS : %s
`, moduleName, relayConfig.serviceName, relayConfig.domain.Host, relayConfig.redisURL, relayConfig.serverBind)
}

// NewMachineryServer create Redis backed Machinery Server from RelayConfig.
func NewMachineryServer(globalConfig *RelayConfig) (*machinery.Server, error) {
	cnf := &config.Config{
		Broker:          globalConfig.redisURL,
		DefaultQueue:    "relay",
		ResultBackend:   globalConfig.redisURL,
		ResultsExpireIn: 1,
	}
	newServer, err := machinery.NewServer(cnf)

	return newServer, err
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

func generatePublicKeyPEMString(publicKey *rsa.PublicKey) string {
	publicKeyByte := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyByte,
		},
	)
	return string(publicKeyPem)
}
