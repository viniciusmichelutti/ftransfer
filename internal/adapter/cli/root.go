package cli

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/viniciusmichelutti/ftransfer/pkg/certs"
)

const version = "0.1.0"

func Execute() error {
	root := &cobra.Command{
		Use:           "ftransfer",
		Short:         "Cross-platform LAN file transfer over mDNS + TLS",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newSendCmd())
	root.AddCommand(newReceiveCmd())
	root.AddCommand(newPeersCmd())
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version)
			return nil
		},
	})
	return root.Execute()
}

func loadCert() (tls.Certificate, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return tls.Certificate{}, err
	}
	return certs.LoadOrCreate(filepath.Join(cfgDir, "ftransfer"))
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
	return ctx, cancel
}
