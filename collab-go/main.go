package main

import (
	"fmt"
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

func validateToken(tokenString, secret string) error {
	if secret == "" {
		return nil
	}
	if tokenString == "" {
		return http.ErrNoCookie
	}
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
	return err
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":1234"
	}
	secret := os.Getenv("JWT_SECRET")
	coreURL := os.Getenv("COLLAB_CORE_URL")
	if coreURL == "" {
		coreURL = "ws://collab-core:1234"
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/ws/doc/") {
			http.NotFound(w, r)
			return
		}
		docID := strings.TrimPrefix(r.URL.Path, "/ws/doc/")
		if docID == "" || strings.Contains(docID, "/") {
			http.Error(w, "invalid doc id", http.StatusBadRequest)
			return
		}

		token := r.URL.Query().Get("token")
		if err := validateToken(token, secret); err != nil {
			log.Printf("collab-gateway unauthorized path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("collab-gateway upgrade failed path=%s remote=%s err=%v", r.URL.Path, r.RemoteAddr, err)
			return
		}
		defer conn.Close()

		u, err := url.Parse(coreURL)
		if err != nil {
			log.Printf("collab-gateway invalid core url=%s err=%v", coreURL, err)
			return
		}
		u.Path = fmt.Sprintf("/doc-%s", docID)
		u.RawQuery = ""

		coreConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("collab-gateway dial core failed url=%s err=%v", u.String(), err)
			return
		}
		defer coreConn.Close()

		errCh := make(chan error, 2)
		go proxyWS(conn, coreConn, errCh)
		go proxyWS(coreConn, conn, errCh)

		<-errCh
	}

	http.HandleFunc("/", handler)

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
