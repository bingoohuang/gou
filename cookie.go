package go_utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

	cookie := http.Cookie{Name: cookieName, Value: cipher, Path: "/", MaxAge: 86400}
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

	cookie := http.Cookie{Domain: domain, Name: cookieName, Value: cipher, Path: "/", MaxAge: 86400}
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
