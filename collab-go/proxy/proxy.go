package proxy

import (
	"io"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Proxy struct {
	coreURL string
	dialer  *websocket.Dialer
}

func NewProxy(coreURL string) *Proxy {
	netDialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
		},
	}
	dialer := &websocket.Dialer{
		NetDialContext: netDialer.DialContext,
	}
	return &Proxy{
		coreURL: coreURL,
		dialer:  dialer,
	}
}

func (p *Proxy) DialCore(docID string) (*websocket.Conn, error) {
	u, err := url.Parse(p.coreURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/doc-" + docID
	u.RawQuery = ""

	coreConn, _, err := p.dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return coreConn, nil
}

func (p *Proxy) ProxyWS(src, dst *websocket.Conn, errCh chan<- error) {
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
