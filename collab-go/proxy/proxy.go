package proxy

import (
	"io"
	"net/url"

	"github.com/gorilla/websocket"
)

type Proxy struct {
	coreURL string
}

func NewProxy(coreURL string) *Proxy {
	return &Proxy{
		coreURL: coreURL,
	}
}

func (p *Proxy) DialCore(docID string) (*websocket.Conn, error) {
	u, err := url.Parse(p.coreURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/doc-" + docID
	u.RawQuery = ""

	coreConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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
