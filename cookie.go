package go_utils

import (
	"net/http"
	"time"
)

func ClearCookie(w http.ResponseWriter, cookieName string) {
	cookie := http.Cookie{Name: cookieName, Value: "", Path: "/", Expires: time.Now().AddDate(-1, 0, 0)}
	http.SetCookie(w, &cookie)
}
