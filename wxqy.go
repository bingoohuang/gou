package go_utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

//var (
//	accessToken            string
//	accessTokenExpiredTime time.Time
//	accessTokenMutex       sync.Mutex
//)

func GetAccessToken(corpId, corpSecret string) (string, error) {
	//accessTokenMutex.Lock()
	//defer accessTokenMutex.Unlock()
	//if accessToken != "" && accessTokenExpiredTime.After(time.Now()) {
	//	return accessToken, nil
	//}

	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpId + "&corpsecret=" + corpSecret
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		//accessToken = ""
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
		//accessToken = tokenResult.AccessToken
		//accessTokenExpiredTime = time.Now().Add(time.Duration(tokenResult.ExpiresInSeconds) * time.Second)
		//return accessToken, nil

		return tokenResult.AccessToken, nil
	}

	return "", errors.New(tokenResult.ErrMsg)
}

func CreateWxQyLoginUrl(cropId, agentId, redirectUri, csrfToken string) string {
	return "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=" +
		cropId + "&agentid=" + agentId + "&redirect_uri=" + url.QueryEscape(redirectUri) + "&state=" + csrfToken
}

func SendWxQyMsg(corpId, corpSecret, agentId, content string) (string, error) {
	// AccessToken是企业号的全局唯一票据，调用接口时需携带AccessToken。
	// AccessToken需要用CorpID和Secret来换取，不同的Secret会返回不同的AccessToken。
	// 正常情况下AccessToken有效期为7200秒，有效期内重复获取返回相同结果。access_token至少保留512字节的存储空间。
	// 企业号可能会出于运营需要，提前使accesstoken失效，企业开发者也应实现accesstoken失效时重试获取的逻辑。
	accessToken, err := GetAccessToken(corpId, corpSecret)
	if err != nil {
		return "", err
	}

	msg := map[string]interface{}{
		"touser": "@all", "toparty": "@all", "totag": "@all", "msgtype": "text", "agentid": agentId, "safe": 0,
		"text": map[string]string{
			"content": content,
		},
	}
	_, err = HttpPost("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+accessToken, msg)
	return accessToken, err
}
