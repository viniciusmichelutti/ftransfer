package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/viniciusmichelutti/ftransfer/internal/adapter/mdns"
	"github.com/viniciusmichelutti/ftransfer/internal/adapter/transport"
	"github.com/viniciusmichelutti/ftransfer/internal/adapter/tui"
	"github.com/viniciusmichelutti/ftransfer/internal/usecase"
	"github.com/viniciusmichelutti/ftransfer/pkg/netutil"
)

func newReceiveCmd() *cobra.Command {
	var (
		out  string
		port int
		name string
		yes  bool
	)
	cmd := &cobra.Command{
		Use:   "receive",
		Short: "Listen for incoming transfers from LAN peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signalContext()
			defer cancel()

			cert, err := loadCert()
			if err != nil {
				return fmt.Errorf("load cert: %w", err)
			}

			outDir, err := resolveOutDir(out)
			if err != nil {
				return err
			}
			if name == "" {
				name = netutil.LocalHostname()
			}

			deps := usecase.ReceiveDeps{
				Transport:  transport.New(cert),
				Discoverer: mdns.New(),
				Prompter:   tui.NewPrompter(),
				Progress:   tui.NewProgress(),
				Name:       name,
				Port:       port,
				AutoAccept: yes,
			}

			ips := netutil.LocalIPv4s()
			fmt.Printf("Listening as %q on port %d (TLS). LAN IPs: %v\n", name, port, ips)
			fmt.Printf("Output dir: %s\n", outDir)
			fmt.Println("Waiting for transfers... (Ctrl+C to stop)")

			return usecase.Receive(ctx, deps, outDir)
		},
	}
	cmd.Flags().StringVar(&out, "out", "", "output directory (default: ~/Downloads/ftransfer)")
	cmd.Flags().IntVar(&port, "port", defaultPort, "TCP port to listen on")
	cmd.Flags().StringVar(&name, "name", "", "advertised peer name (default: hostname)")
	cmd.Flags().BoolVar(&yes, "yes", false, "auto-accept all incoming transfers")
	return cmd
}

func resolveOutDir(out string) (string, error) {
	if out != "" {
		return out, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Downloads", "ftransfer"), nil
}
