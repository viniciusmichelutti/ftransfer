package usecase

import (
	"context"
	"fmt"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
	"github.com/viniciusmichelutti/ftransfer/pkg/archive"
	"github.com/viniciusmichelutti/ftransfer/pkg/fsutil"
)

type SendDeps struct {
	Transport Transport
	Progress  Progress
	Sender    string
}

func Send(ctx context.Context, deps SendDeps, target domain.Peer, sources []string) error {
	entries, total, err := fsutil.Walk(sources)
	if err != nil {
		return fmt.Errorf("walk sources: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("no files to send")
	}

	manifest := domain.Manifest{
		Sender:     deps.Sender,
		TotalBytes: total,
		Files:      stripAbs(entries),
	}

	conn, err := deps.Transport.Dial(ctx, target.Addr())
	if err != nil {
		return fmt.Errorf("dial %s: %w", target.Addr(), err)
	}
	defer conn.Close()

	if err := writeFrame(conn, manifest); err != nil {
		return fmt.Errorf("send manifest: %w", err)
	}

	var resp domain.AcceptResponse
	if err := readFrame(conn, &resp); err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if !resp.Accepted {
		reason := resp.Reason
		if reason == "" {
			reason = "rejected by receiver"
		}
		return fmt.Errorf("%s", reason)
	}

	deps.Progress.Start("Sending", total)
	defer deps.Progress.Finish()

	w := deps.Progress.WrapWriter(conn)
	if err := archive.StreamTo(w, entries); err != nil {
		return fmt.Errorf("stream: %w", err)
	}
	return nil
}

func stripAbs(entries []domain.FileEntry) []domain.FileEntry {
	out := make([]domain.FileEntry, len(entries))
	for i, e := range entries {
		e.AbsPath = ""
		out[i] = e
	}
	return out
}
