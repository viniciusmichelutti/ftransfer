package usecase

import (
	"context"
	"fmt"
	"net"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
	"github.com/viniciusmichelutti/ftransfer/pkg/archive"
)

type ReceiveDeps struct {
	Transport  Transport
	Discoverer Discoverer
	Prompter   Prompter
	Progress   Progress

	Name       string // advertised mDNS name
	Port       int
	AutoAccept bool
}

func Receive(ctx context.Context, deps ReceiveDeps, outDir string) error {
	addr := fmt.Sprintf(":%d", deps.Port)
	ln, err := deps.Transport.Listen(addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	defer ln.Close()

	if err := deps.Discoverer.Advertise(ctx, deps.Name, deps.Port); err != nil {
		return fmt.Errorf("advertise: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept: %w", err)
		}
		// Handle one transfer at a time. KISS — receiver is a foreground command.
		if err := handleConn(deps, conn, outDir); err != nil {
			fmt.Printf("transfer failed: %v\n", err)
		}
	}
}

func handleConn(deps ReceiveDeps, conn net.Conn, outDir string) error {
	defer conn.Close()

	var manifest domain.Manifest
	if err := readFrame(conn, &manifest); err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	accept := deps.AutoAccept
	if !accept {
		msg := fmt.Sprintf("Incoming from %s: %d item(s), %s — accept?",
			manifest.Sender, len(manifest.Files), humanBytes(manifest.TotalBytes))
		ok, err := deps.Prompter.Confirm(msg)
		if err != nil {
			return err
		}
		accept = ok
	}

	resp := domain.AcceptResponse{Accepted: accept}
	if !accept {
		resp.Reason = "declined by user"
	}
	if err := writeFrame(conn, resp); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	if !accept {
		return nil
	}

	deps.Progress.Start("Receiving", manifest.TotalBytes)
	defer deps.Progress.Finish()

	r := deps.Progress.Wrap(conn)
	if err := archive.ExtractTo(r, outDir); err != nil {
		return fmt.Errorf("extract: %w", err)
	}
	fmt.Printf("Saved to %s\n", outDir)
	return nil
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}
