package htt

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bingoohuang/gou/enc"
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

	cipher, err := enc.CBCEncrypt(encryptKey, string(cookieVal))
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

	cipher, err := enc.CBCEncrypt(encryptKey, string(cookieVal))
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
	decrypted, err := enc.CBCDecrypt(encryptKey, cookie.Value)
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
	UserID    string
	Name      string
	Avatar    string
	CsrfToken string
	Expired   time.Time
}

func (t *CookieValueImpl) ExpiredTime() time.Time {
	return t.Expired
}

type MustAuthParam struct {
	EncryptKey  string
	CookieName  string
	RedirectURI string
	LocalURL    string
	ForceLogin  bool
}

func PrepareMustAuthFlag(param *MustAuthParam) {
	flag.StringVar(&param.EncryptKey, "key", "", "key to encryption or decryption")
	flag.StringVar(&param.CookieName, "cookieName", "i-raiyee-cn-auth", "cookieName")
	flag.StringVar(&param.RedirectURI, "redirectUri", "", "redirectUri")
	flag.StringVar(&param.LocalURL, "localUrl", "", "localUrl")
	flag.BoolVar(&param.ForceLogin, "forceLogin", false, "forceLogin required")
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
func MustAuth(fn http.HandlerFunc, param MustAuthParam, cookieContextKey interface{}) http.HandlerFunc { // nolint
	if !param.ForceLogin {
		return fn
	}

	return func(w http.ResponseWriter, r *http.Request) {
		cookie := CookieValueImpl{}
		err := ReadCookie(r, param.EncryptKey, param.CookieName, &cookie)
		if err == nil && cookie.Name != "" {
			ctx := context.WithValue(r.Context(), cookieContextKey, &cookie)
			fn.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		urlx := param.RedirectURI
		if strings.Contains(param.RedirectURI, "?") {
			urlx += "&"
		} else {
			urlx += "?"
		}

		urlx += "cookie=" + param.CookieName + "&"
		urlx += "redirect=" + url.QueryEscape(param.LocalURL+r.RequestURI)
		http.Redirect(w, r, urlx, http.StatusFound)
	}
}
