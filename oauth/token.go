package oauth

import (
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
)

func EncodeToken(token *oauth2.Token) (string, error) {
	jsonToken, err := tokenToJSON(token)
	if err != nil {
		return "", err
	}
	return jsonToken, nil
}

func DecodeToken(tokenInterface interface{}) (*oauth2.Token, error) {
	tokenString, ok := tokenInterface.(string)
	if !ok {
		return nil, errors.New("token is not string")
	}

	token, err := tokenFromJSON(tokenString)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
