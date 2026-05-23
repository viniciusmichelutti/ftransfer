package mdns

import (
	"context"
	"time"

	"github.com/grandcat/zeroconf"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

const service = "_ftransfer._tcp"
const domainSuffix = "local."

type Discoverer struct {
	server *zeroconf.Server
}

func New() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Advertise(ctx context.Context, name string, port int) error {
	srv, err := zeroconf.Register(name, service, domainSuffix, port, []string{"ftransfer=1"}, nil)
	if err != nil {
		return err
	}
	d.server = srv
	go func() {
		<-ctx.Done()
		srv.Shutdown()
	}()
	return nil
}

func (d *Discoverer) Browse(ctx context.Context) ([]domain.Peer, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	browseCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancel()

	if err := resolver.Browse(browseCtx, service, domainSuffix, entries); err != nil {
		return nil, err
	}

	var peers []domain.Peer
	for e := range entries {
		host := ""
		if len(e.AddrIPv4) > 0 {
			host = e.AddrIPv4[0].String()
		} else if len(e.AddrIPv6) > 0 {
			host = "[" + e.AddrIPv6[0].String() + "]"
		}
		if host == "" {
			continue
		}
		peers = append(peers, domain.Peer{
			Name: e.Instance,
			Host: host,
			Port: e.Port,
		})
	}
	return peers, nil
}
