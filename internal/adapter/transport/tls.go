package transport

import (
	"context"
	"crypto/tls"
	"net"
)

type TLS struct {
	cert tls.Certificate
}

func New(cert tls.Certificate) *TLS {
	return &TLS{cert: cert}
}

func (t *TLS) serverConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{t.cert},
		MinVersion:   tls.VersionTLS13,
		ClientAuth:   tls.NoClientCert,
	}
}

func (t *TLS) clientConfig() *tls.Config {
	return &tls.Config{
		Certificates:       []tls.Certificate{t.cert},
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true, // auth is the manual receiver prompt; v1 design
	}
}

func (t *TLS) Listen(addr string) (net.Listener, error) {
	return tls.Listen("tcp", addr, t.serverConfig())
}

func (t *TLS) Dial(ctx context.Context, addr string) (net.Conn, error) {
	d := tls.Dialer{Config: t.clientConfig()}
	return d.DialContext(ctx, "tcp", addr)
}
