// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mqttsnippets

// [START iot_mqtt_jwt]

import (
	"errors"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// createJWT creates a Cloud IoT Core JWT for the given project id.
// algorithm can be one of ["RSA256", "ES256"].
func createJWT(projectID string, privateKeyPath string, algorithm string, expiration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(expiration).Unix(),
	}

	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), claims)

	switch algorithm {
	case "RS256":
		privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		return token.SignedString(privKey)
	case "ES256":
		privKey, _ := jwt.ParseECPrivateKeyFromPEM(keyBytes)
		return token.SignedString(privKey)
	}

	return "", errors.New("Cannot find JWT algorithm. Specify 'ES256' or 'RS256'")
}

// [END iot_mqtt_jwt]
