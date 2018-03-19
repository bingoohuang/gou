package go_utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type WxLoginUserId struct {
	UserId  string `json:"UserId"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 企业微信：https://work.weixin.qq.com/

// refer https://work.weixin.qq.com/api/doc#10719
func GetLoginUserId(accessToken, code string) (string, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var wxLoginUserId WxLoginUserId
	err = json.Unmarshal(body, &wxLoginUserId)
	if err != nil {
		return "", err
	}
	if wxLoginUserId.UserId == "" {
		return "", errors.New(string(body))
	}

	return wxLoginUserId.UserId, nil
}

type WxUserInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	UserId string `json:"userid"`
}

func GetUserInfo(accessToken, userId string) (*WxUserInfo, error) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=" + accessToken + "&userid=" + userId
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wxUserInfo WxUserInfo
	err = json.Unmarshal(body, &wxUserInfo)
	if err != nil {
		return nil, err
	}

	return &wxUserInfo, nil
}

type TokenResult struct {
	ErrCode          int    `json:"errcode"`
	ErrMsg           string `json:"errmsg"`
	AccessToken      string `json:"access_token"`
	ExpiresInSeconds int    `json:"expires_in"`
}

var (
	accessToken            string
	accessTokenExpiredTime time.Time
	accessTokenMutex       sync.Mutex
)

func GetAccessToken(corpId, corpSecret string) (string, error) {
	accessTokenMutex.Lock()
	defer accessTokenMutex.Unlock()
	if accessToken != "" && accessTokenExpiredTime.After(time.Now()) {
		return accessToken, nil
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpId + "&corpsecret=" + corpSecret
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		accessToken = ""
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResult TokenResult
	json.Unmarshal(body, &tokenResult)
	if tokenResult.ErrCode == 0 {
		accessToken = tokenResult.AccessToken
		accessTokenExpiredTime = time.Now().Add(time.Duration(tokenResult.ExpiresInSeconds) * time.Second)
		return accessToken, nil
	}

	return "", errors.New(tokenResult.ErrMsg)
}

func CreateWxQyLoginUrl(cropId, agentId, redirectUri, csrfToken string) string {
	return "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=" +
		cropId + "&agentid=" + agentId + "&redirect_uri=" + redirectUri + "&state=" + csrfToken
}

//func redirectWxQyLogin(w http.ResponseWriter, r *http.Request, url string) {
//	http.Redirect(w, r, url, 302) // Temporarily Move
//}

type CookieValue interface {
	ExpiredTime() time.Time
}

func WriteUserInfoCookie(w http.ResponseWriter, encryptKey, cookieName string, cookieValue CookieValue) error {
	cookieVal, err := json.Marshal(cookieValue)
	if err != nil {
		return err
	}

	cipher, err := CBCEncrypt(encryptKey, string(cookieVal))
	if err != nil {
		return err
	}

	cookie := http.Cookie{Name: cookieName, Value: cipher, Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)

	return nil
}

func ReadLoginCookie(r *http.Request, encryptKey, cookieName string, cookieValue CookieValue) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	log.Println("cookie value:", cookie.Value)
	decrypted, err := CBCDecrypt(encryptKey, cookie.Value)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(decrypted), cookieValue)
	if err != nil {
		log.Println("unamrshal error:", err)
		return err
	}

	log.Println("cookie parsed:", cookieValue)

	return nil
}
