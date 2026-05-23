package cli

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/viniciusmichelutti/ftransfer/internal/adapter/mdns"
	"github.com/viniciusmichelutti/ftransfer/internal/adapter/transport"
	"github.com/viniciusmichelutti/ftransfer/internal/adapter/tui"
	"github.com/viniciusmichelutti/ftransfer/internal/domain"
	"github.com/viniciusmichelutti/ftransfer/internal/usecase"
	"github.com/viniciusmichelutti/ftransfer/pkg/netutil"
)

const defaultPort = 7777

func newSendCmd() *cobra.Command {
	var to string

	cmd := &cobra.Command{
		Use:   "send PATH [PATH...]",
		Short: "Send files or folders to a LAN peer",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signalContext()
			defer cancel()

			cert, err := loadCert()
			if err != nil {
				return fmt.Errorf("load cert: %w", err)
			}

			tr := transport.New(cert)
			pr := tui.NewProgress()
			prompter := tui.NewPrompter()

			target, err := resolveTarget(ctx, to, prompter)
			if err != nil {
				return err
			}

			deps := usecase.SendDeps{
				Transport: tr,
				Progress:  pr,
				Sender:    netutil.LocalHostname(),
			}
			return usecase.Send(ctx, deps, target, args)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "target peer name or IP[:port]")
	return cmd
}

func resolveTarget(ctx context.Context, to string, prompter *tui.Prompter) (domain.Peer, error) {
	if to != "" {
		return parseTarget(to)
	}
	fmt.Println("Discovering peers on LAN...")
	disc := mdns.New()
	peers, err := disc.Browse(ctx)
	if err != nil {
		return domain.Peer{}, fmt.Errorf("browse: %w", err)
	}
	return prompter.PickPeer(peers)
}

func parseTarget(s string) (domain.Peer, error) {
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			return domain.Peer{Name: s, Host: s, Port: defaultPort}, nil
		}
		return domain.Peer{}, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return domain.Peer{}, fmt.Errorf("invalid port: %w", err)
	}
	return domain.Peer{Name: host, Host: host, Port: port}, nil
}
