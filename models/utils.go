package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/go-redis/redis"
	"io/ioutil"
)

func redisHGetOrCreateWithDefault(redisClient *redis.Client, key string, field string, defaultValue string) (string, error) {
	keyExist, err := redisClient.HExists(key, field).Result()
	if err != nil {
		return "", err
	}
	if keyExist {
		value, err := redisClient.HGet(key, field).Result()
		if err != nil {
			return "", err
		}
		return value, nil
	} else {
		_, err := redisClient.HSet(key, field, defaultValue).Result()
		if err != nil {
			return "", err
		}
		return defaultValue, nil
	}
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
