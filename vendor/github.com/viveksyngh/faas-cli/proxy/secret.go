package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	types "github.com/openfaas/faas-provider/types"
)

const (
	secretEndpoint = "/system/secrets"
)

// GetSecretList get secrets list
func GetSecretList(gateway string, tlsInsecure bool, namespace string) ([]types.Secret, error) {
	return GetSecretListToken(gateway, tlsInsecure, "", namespace)
}

// GetSecretListToken get secrets lists with taken as auth
func GetSecretListToken(gateway string, tlsInsecure bool, token, namespace string) ([]types.Secret, error) {
	var results []types.Secret

	gateway = strings.TrimRight(gateway, "/")
	gatewayURL, err := url.Parse(gateway)
	if err != nil {
		return results, fmt.Errorf("invalid gateway URL: %s", gateway)
	}
	gatewayURL.Path = path.Join(gatewayURL.Path, secretEndpoint)
	if len(namespace) > 0 {
		q := gatewayURL.Query()
		q.Set("namespace", namespace)
		gatewayURL.RawQuery = q.Encode()
	}

	client := MakeHTTPClient(&defaultCommandTimeout, tlsInsecure)
	getRequest, err := http.NewRequest(http.MethodGet, gatewayURL.String(), nil)
	if len(token) > 0 {
		SetToken(getRequest, token)
	} else {
		SetAuth(getRequest, gateway)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	res, err := client.Do(getRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:

		bytesOut, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read result from OpenFaaS on URL: %s", gateway)
		}

		jsonErr := json.Unmarshal(bytesOut, &results)
		if jsonErr != nil {
			return nil, fmt.Errorf("cannot parse result from OpenFaaS on URL: %s\n%s", gateway, jsonErr.Error())
		}

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized access, run \"faas-cli login\" to setup authentication for this server")

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return results, nil
}

// UpdateSecret update a secret via the OpenFaaS API by name
func UpdateSecret(gateway string, secret types.Secret, tlsInsecure bool) (int, string) {
	return UpdateSecretToken(gateway, secret, tlsInsecure, "")
}

// UpdateSecretToken update a secret with token as auth
func UpdateSecretToken(gateway string, secret types.Secret, tlsInsecure bool, token string) (int, string) {
	var output string

	gateway = strings.TrimRight(gateway, "/")
	client := MakeHTTPClient(&defaultCommandTimeout, tlsInsecure)

	reqBytes, _ := json.Marshal(&secret)

	putRequest, err := http.NewRequest(http.MethodPut, gateway+"/system/secrets", bytes.NewBuffer(reqBytes))
	if len(token) > 0 {
		SetToken(putRequest, token)
	} else {
		SetAuth(putRequest, gateway)
	}

	if err != nil {
		output += fmt.Sprintf("cannot connect to OpenFaaS on URL: %s", gateway)
		return http.StatusInternalServerError, output
	}

	res, err := client.Do(putRequest)
	if err != nil {
		output += fmt.Sprintf("cannot connect to OpenFaaS on URL: %s", gateway)
		return http.StatusInternalServerError, output
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		output += fmt.Sprintf("Updated: %s\n", res.Status)
		break

	case http.StatusNotFound:
		output += fmt.Sprintf("unable to find secret: %s", secret.Name)

	case http.StatusUnauthorized:
		output += fmt.Sprintf("unauthorized access, run \"faas-cli login\" to setup authentication for this server")

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			output += fmt.Sprintf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return res.StatusCode, output
}

// RemoveSecret remove a secret via the OpenFaaS API by name
func RemoveSecret(gateway string, secret types.Secret, tlsInsecure bool) error {
	return RemoveSecretToken(gateway, secret, tlsInsecure, "")
}

// RemoveSecretToken remove a secret with token as auth
func RemoveSecretToken(gateway string, secret types.Secret, tlsInsecure bool, token string) error {

	gateway = strings.TrimRight(gateway, "/")
	client := MakeHTTPClient(&defaultCommandTimeout, tlsInsecure)

	body, _ := json.Marshal(secret)

	getRequest, err := http.NewRequest(http.MethodDelete, gateway+"/system/secrets", bytes.NewBuffer(body))

	if len(token) > 0 {
		SetToken(getRequest, token)
	} else {
		SetAuth(getRequest, gateway)
	}

	if err != nil {
		return fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	res, err := client.Do(getRequest)
	if err != nil {
		return fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		break
	case http.StatusNotFound:
		return fmt.Errorf("unable to find secret: %s", secret.Name)
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized access, run \"faas-cli login\" to setup authentication for this server")

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return nil
}

// CreateSecret create secret
func CreateSecret(gateway string, secret types.Secret, tlsInsecure bool) (int, string) {
	return CreateSecretToken(gateway, secret, tlsInsecure, "")
}

// CreateSecretToken create secret with token as auth
func CreateSecretToken(gateway string, secret types.Secret, tlsInsecure bool, token string) (int, string) {
	var output string

	gateway = strings.TrimRight(gateway, "/")

	reqBytes, _ := json.Marshal(&secret)
	reader := bytes.NewReader(reqBytes)

	client := MakeHTTPClient(&defaultCommandTimeout, tlsInsecure)
	request, err := http.NewRequest(http.MethodPost, gateway+"/system/secrets", reader)

	if len(token) > 0 {
		SetToken(request, token)
	} else {
		SetAuth(request, gateway)
	}

	if err != nil {
		output += fmt.Sprintf("cannot connect to OpenFaaS on URL: %s\n", gateway)
		return http.StatusInternalServerError, output
	}

	res, err := client.Do(request)
	if err != nil {
		output += fmt.Sprintf("cannot connect to OpenFaaS on URL: %s\n", gateway)
		return http.StatusInternalServerError, output
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		output += fmt.Sprintf("Created: %s\n", res.Status)

	case http.StatusUnauthorized:
		output += fmt.Sprintln("unauthorized access, run \"faas-cli login\" to setup authentication for this server")

	case http.StatusConflict:
		output += fmt.Sprintf("secret with the name %q already exists\n", secret.Name)

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			output += fmt.Sprintf("server returned unexpected status code: %d - %s\n", res.StatusCode, string(bytesOut))
		}
	}

	return res.StatusCode, output
}
