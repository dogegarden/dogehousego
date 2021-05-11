package dogehousego

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthResponse struct {
	Username string `json:"username"`
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Error string `json:"error"`
}

type authPost struct {
	 ApiKey string `json:"apiKey"`
}


func Auth(apiKey string) (*AuthResponse, error) {
	postData := authPost{ApiKey: apiKey};
	dataString, err := json.Marshal(postData);

	if err != nil {
		return nil, errors.New("Failed to marshal apiKey. Error: " + err.Error());
	}

	resp, err := http.Post(HttpBaseUrl+ "/bot/auth", "application/json", bytes.NewBuffer(dataString));

	if err != nil {
		return nil, errors.New("Failed to send auth request. Error: " + err.Error());
	}

	var retVal *AuthResponse;

	err = json.NewDecoder(resp.Body).Decode(&retVal);
	if err != nil {
		return nil, errors.New("Failed to unmarshal response auth data. Error: " + err.Error());
	}

	if retVal.Error != "" {
		return nil, errors.New("Failed to get tokens from api-key. Reason: " + retVal.Error);
	}

	return retVal, nil;
}