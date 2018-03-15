package ddpportal

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
)

type OauthProvider interface {
	CheckToken(token string) bool
}

type OauthAdeoProvider struct {
	baseUrl      string
	clientID     string
	clientSecret string
}

func (o OauthAdeoProvider) CheckToken(token string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s/check_token?token=%s", o.baseUrl, token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Error connecting to Oauth server: %s", o.baseUrl))
		return false
	}

	clientIDSecret := fmt.Sprintf("%s:%s", o.clientID, o.clientSecret)
	encodedClientID := base64.StdEncoding.EncodeToString([]byte(clientIDSecret))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedClientID))

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Error connecting to Oauth server: %s", o.baseUrl))
		return false
	}

	if resp.StatusCode != 200 {
		log.Error(fmt.Sprintf("Authentication failed: %s", resp.Status))
		return false
	}

	return true
}

func NewOauthAdeoProvider(baseUrl string, clientID string, clientSecret string) OauthAdeoProvider {
	return OauthAdeoProvider{baseUrl: baseUrl, clientID: clientID, clientSecret: clientSecret}
}

func NewOauthAuthorizerMiddleware(oauth OauthProvider) func(http.Handler) http.Handler {
	middleWare := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			token := strings.Split(bearer, " ")[1]
			if oauth.CheckToken(token) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Not Authorized", 401)
			}
		}
		return http.HandlerFunc(fn)
	}
	return middleWare
}
