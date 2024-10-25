package httphelper

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func GetSession(c *gin.Context) *sessions.Session {
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return nil
	}
	return session
}

func GetSessionValue(c *gin.Context, key string) any {
	session := GetSession(c)
	return session.Values[key]
}

func SetSessionValue(c *gin.Context, key string, value any) {
	session := GetSession(c)
	session.Values[key] = value
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CheckSession(c *gin.Context, key string, value any) bool {
	val := GetSessionValue(c, key)
	return val == value
}
