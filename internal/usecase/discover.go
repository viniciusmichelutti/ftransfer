package usecase

import (
	"context"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

func Discover(ctx context.Context, d Discoverer) ([]domain.Peer, error) {
	return d.Browse(ctx)
}
