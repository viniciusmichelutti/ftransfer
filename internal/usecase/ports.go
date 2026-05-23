package usecase

import (
	"context"
	"io"
	"net"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

// Discoverer advertises the local peer and browses for remote peers via mDNS.
type Discoverer interface {
	Advertise(ctx context.Context, name string, port int) error
	Browse(ctx context.Context) ([]domain.Peer, error)
}

// Transport opens authenticated, encrypted streams between peers.
type Transport interface {
	Listen(addr string) (net.Listener, error)
	Dial(ctx context.Context, addr string) (net.Conn, error)
}

// Prompter handles all user interaction (peer selection, accept/reject).
type Prompter interface {
	Confirm(message string) (bool, error)
	PickPeer(peers []domain.Peer) (domain.Peer, error)
}

// Progress reports byte-level transfer progress to the UI.
type Progress interface {
	Start(label string, total int64)
	Wrap(r io.Reader) io.Reader
	WrapWriter(w io.Writer) io.Writer
	Finish()
}
