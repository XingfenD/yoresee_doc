package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func validateToken(tokenString, secret string) bool {
	if secret == "" {
		return true
	}
	if tokenString == "" {
		return false
	}
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return err == nil
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":1235"
	}
	secret := os.Getenv("JWT_SECRET")
	coreURL := os.Getenv("COLLAB_CORE_URL")
	if coreURL == "" {
		coreURL = "ws://collab:1234"
	}

	http.HandleFunc("/collab", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/collab") {
			http.NotFound(w, r)
			return
		}

		token := r.URL.Query().Get("token")
		if !validateToken(token, secret) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		u, err := url.Parse(coreURL)
		if err != nil {
			return
		}
		u.Path = r.URL.Path
		u.RawQuery = ""

		coreConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			return
		}
		defer coreConn.Close()

		errCh := make(chan error, 2)
		go proxyWS(conn, coreConn, errCh)
		go proxyWS(coreConn, conn, errCh)

		<-errCh
	})

	log.Printf("collab-gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func proxyWS(src, dst *websocket.Conn, errCh chan<- error) {
	for {
		msgType, r, err := src.NextReader()
		if err != nil {
			errCh <- err
			return
		}
		w, err := dst.NextWriter(msgType)
		if err != nil {
			errCh <- err
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			errCh <- err
			return
		}
		if err := w.Close(); err != nil {
			errCh <- err
			return
		}
	}
}
