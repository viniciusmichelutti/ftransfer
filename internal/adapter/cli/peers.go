package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/viniciusmichelutti/ftransfer/internal/adapter/mdns"
	"github.com/viniciusmichelutti/ftransfer/internal/usecase"
)

func newPeersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "peers",
		Short: "List LAN peers currently advertising ftransfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signalContext()
			defer cancel()

			peers, err := usecase.Discover(ctx, mdns.New())
			if err != nil {
				return err
			}
			if len(peers) == 0 {
				fmt.Println("No peers found.")
				return nil
			}
			for _, p := range peers {
				fmt.Printf("%-24s %s\n", p.Name, p.Addr())
			}
			return nil
		},
	}
}
