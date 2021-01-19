// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type LoginParameters struct {
	ClientID     string
	ClientSecret string
	Scopes       string
	Token        string
	URL          string
	RedirectURL  string
	AuthorizeURL string
}

type RefreshParameters struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	URL          string
}

type AuthorizationResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type UserAuthorizationQueryResponse struct {
	Code  string
	State string
}

type LoginResponse struct {
	Response  AuthorizationResponse
	ExpiresAt time.Time
}

const ClientCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=client_credentials"

const UserCredentialsURL = "https://id.twitch.tv/oauth2/token?grant_type=authorization_code"
const UserAuthorizeURL = "https://id.twitch.tv/oauth2/authorize?response_type=code"

const RefreshTokenURL = "https://id.twitch.tv/oauth2/token?grant_type=refresh_token"

const RevokeTokenURL = "https://id.twitch.tv/oauth2/revoke"

func ClientCredentialsLogin(p LoginParameters) (LoginResponse, error) {
	twitchClientCredentialsURL := fmt.Sprintf(`%s&client_id=%s&client_secret=%s`, p.URL, p.ClientID, p.ClientSecret)

	resp, err := loginRequest(http.MethodPost, twitchClientCredentialsURL, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	r, err := handleLoginResponse(resp.Body)
	if err != nil {
		log.Printf("Error handling login: %v", err)
		return LoginResponse{}, nil
	}

	log.Printf("App Access Token: %s", r.Response.AccessToken)
	return r, nil
}

func UserCredentialsLogin(p LoginParameters) (LoginResponse, error) {
	twitchAuthorizeURL := fmt.Sprintf(`%s&client_id=%s&redirect_uri=%s&force_verify=true`, p.AuthorizeURL, p.ClientID, p.RedirectURL)

	if p.Scopes != "" {
		twitchAuthorizeURL += "&scope=" + p.Scopes
	}

	state, err := generateState()
	if err != nil {
		log.Fatal(err.Error())
	}

	twitchAuthorizeURL += "&state=" + state

	openBrowser(twitchAuthorizeURL)

	ur, err := userAuthServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	if ur.State != state {
		log.Fatal("state mismatch")
	}

	twitchUserTokenURL := fmt.Sprintf(`%s&client_id=%s&client_secret=%s&redirect_uri=%s&code=%s`, p.URL, p.ClientID, p.ClientSecret, p.RedirectURL, ur.Code)
	resp, err := loginRequest(http.MethodPost, twitchUserTokenURL, nil)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	r, err := handleLoginResponse(resp.Body)
	if err != nil {
		log.Printf("Error handling login: %v", err)
		return LoginResponse{}, nil
	}

	log.Printf("User Access Token: %s\nRefresh Token: %s\nExpires At: %s\nScopes: %s", r.Response.AccessToken, r.Response.RefreshToken, r.ExpiresAt, r.Response.Scope)
	return r, nil
}

func CredentialsLogout(p LoginParameters) (LoginResponse, error) {
	twitchClientCredentialsURL := fmt.Sprintf(`%s?client_id=%s&token=%s`, p.URL, p.ClientID, p.Token)

	resp, err := loginRequest(http.MethodPost, twitchClientCredentialsURL, nil)
	if err != nil {
		log.Printf(err.Error())
		return LoginResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API responded with an error while revoking token: %v", resp.Body)
		return LoginResponse{}, errors.New("API responded with an error while revoking token")
	}

	log.Printf("Token %s has been successfully revoked.", p.Token)
	return LoginResponse{}, nil
}

func RefreshUserToken(p RefreshParameters) (LoginResponse, error) {
	twitchRefreshTokenURL := fmt.Sprintf(`%s&client_id=%s&client_secret=%s&redirect_uri=&refresh_token=%s`, p.URL, p.ClientID, p.ClientSecret, p.RefreshToken)
	resp, err := loginRequest(http.MethodPost, twitchRefreshTokenURL, nil)
	if err != nil {
		return LoginResponse{}, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return LoginResponse{}, errors.New("Error with client while refreshing. Please rerun twitch configure")
	}

	r, err := handleLoginResponse(resp.Body)
	if err != nil {
		log.Printf("Error handling login: %v", err)
		return LoginResponse{}, err
	}

	return r, nil
}

func handleLoginResponse(body []byte) (LoginResponse, error) {
	var r AuthorizationResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return LoginResponse{}, err
	}
	expiresAt := time.Now().Add(time.Duration(int64(time.Second) * int64(r.ExpiresIn)))
	storeInConfig(r.AccessToken, r.RefreshToken, r.Scope, expiresAt)

	return LoginResponse{
		Response:  r,
		ExpiresAt: expiresAt,
	}, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func userAuthServer() (UserAuthorizationQueryResponse, error) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":3000", Handler: m}
	userAuth := make(chan UserAuthorizationQueryResponse)
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Feel free to close this browser window."))

		var u = UserAuthorizationQueryResponse{
			Code:  r.URL.Query().Get("code"),
			State: r.URL.Query().Get("state"),
		}
		userAuth <- u
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
	}()

	userAuthResponse := <-userAuth

	s.Shutdown(context.Background())
	return userAuthResponse, nil
}

func storeInConfig(token string, refresh string, scopes []string, expiresAt time.Time) {
	viper.Set("accessToken", token)
	viper.Set("refreshToken", refresh)
	viper.Set("tokenScopes", scopes)
	viper.Set("tokenExpiration", expiresAt.Format(time.RFC3339))

	err := viper.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = viper.SafeWriteConfig()
	}

	if err != nil {
		log.Fatalf("Error writing configuration: %s", err)
	}
}