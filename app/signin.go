package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  os.Getenv("DOMAIN_NAME") + "/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

type GoogleUserId struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthState", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthState")
	if r.FormValue("state") != oauthState.Value {
		log.Printf("invalid google oauth state cookie:%s\n state:%s\n", oauthState.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		errMsg := fmt.Sprintf("invalid google oauth state cookie:%s\n state:%s\n", oauthState.Value, r.FormValue("state"))
		log.Println(err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// store id info into session cookie
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
	}

	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
	}
	session.Values["id"] = userInfo.ID
	err = sessions.Save(r, w)
	if err != nil {
		log.Println(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//fmt.Fprint(w, string(data))
}

var oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	res, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get userInfo %s\n", err.Error())
	}
	return ioutil.ReadAll(res.Body)
}
