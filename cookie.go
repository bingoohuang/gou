package go_utils

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"
)

func ClearCookie(w http.ResponseWriter, cookieName string) {
	cookie := http.Cookie{Name: cookieName, Value: "", Path: "/", Expires: time.Now().AddDate(-1, 0, 0)}
	http.SetCookie(w, &cookie)
}

type CookieValue interface {
	ExpiredTime() time.Time
}

func WriteCookie(w http.ResponseWriter, encryptKey, cookieName string, cookieValue CookieValue) error {
	cookieVal, err := json.Marshal(cookieValue)
	if err != nil {
		return err
	}

	cipher, err := CBCEncrypt(encryptKey, string(cookieVal))
	if err != nil {
		return err
	}

	maxAge := cookieValue.ExpiredTime().Unix() - time.Now().Unix()
	cookie := http.Cookie{Name: cookieName, Value: cipher, Path: "/", MaxAge: int(maxAge)}
	http.SetCookie(w, &cookie)

	return nil
}

func WriteDomainCookie(w http.ResponseWriter, domain, encryptKey, cookieName string, cookieValue CookieValue) error {
	cookieVal, err := json.Marshal(cookieValue)
	if err != nil {
		return err
	}

	cipher, err := CBCEncrypt(encryptKey, string(cookieVal))
	if err != nil {
		return err
	}

	maxAge := cookieValue.ExpiredTime().Unix() - time.Now().Unix()
	cookie := http.Cookie{Domain: domain, Name: cookieName, Value: cipher, Path: "/", MaxAge: int(maxAge)}
	http.SetCookie(w, &cookie)

	return nil
}

func ReadCookie(r *http.Request, encryptKey, cookieName string, cookieValue CookieValue) error {
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

	if cookieValue.ExpiredTime().Before(time.Now()) {
		return errors.New("cookie expired")
	}

	return nil
}

type CookieValueImpl struct {
	UserId    string
	Name      string
	Avatar    string
	CsrfToken string
	Expired   time.Time
}

func (t *CookieValueImpl) ExpiredTime() time.Time {
	return t.Expired
}

type MustAuthParam struct {
	EncryptKey  *string
	CookieName  *string
	RedirectUri *string
	LocalUrl    *string
	ForceLogin  *bool
}

func PrepareMustAuthFlag(param *MustAuthParam) {
	param.EncryptKey = flag.String("key", "", "key to encryption or decryption")
	param.CookieName = flag.String("cookieName", "i-raiyee-cn-auth", "cookieName")
	param.RedirectUri = flag.String("redirectUri", "", "redirectUri")
	param.LocalUrl = flag.String("localUrl", "", "localUrl")
	param.ForceLogin = flag.Bool("forceLogin", false, "forceLogin required")
}

/*
	fmt.Println(r.Proto)
	// output:HTTP/1.1
	fmt.Println(r.TLS)
	// output: <nil>
	fmt.Println(r.Host)
	// output: localhost:9090
	fmt.Println(r.RequestURI)
	// output: /index?id=1
*/
func MustAuth(fn http.HandlerFunc, param MustAuthParam) http.HandlerFunc {
	if !*param.ForceLogin {
		return fn
	}

	return func(w http.ResponseWriter, r *http.Request) {
		cookie := CookieValueImpl{}
		err := ReadCookie(r, *param.EncryptKey, *param.CookieName, &cookie)
		if err == nil && cookie.Name != "" {
			ctx := context.WithValue(r.Context(), "CookieValue", &cookie)
			fn.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		urlx := *param.RedirectUri + "?redirect=" + url.QueryEscape(*param.LocalUrl+r.RequestURI)
		http.Redirect(w, r, urlx, 302)
	}
}
