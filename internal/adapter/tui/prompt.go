package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

type Prompter struct {
	in *bufio.Reader
}

func NewPrompter() *Prompter {
	return &Prompter{in: bufio.NewReader(os.Stdin)}
}

func (p *Prompter) Confirm(message string) (bool, error) {
	fmt.Printf("%s [Y/n]: ", message)
	line, err := p.in.ReadString('\n')
	if err != nil {
		return false, err
	}
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" || line == "y" || line == "yes" {
		return true, nil
	}
	return false, nil
}

func (p *Prompter) PickPeer(peers []domain.Peer) (domain.Peer, error) {
	if len(peers) == 0 {
		return domain.Peer{}, fmt.Errorf("no peers found on LAN")
	}
	if len(peers) == 1 {
		fmt.Printf("Found peer: %s (%s)\n", peers[0].Name, peers[0].Addr())
		return peers[0], nil
	}
	fmt.Println("Discovered peers:")
	for i, peer := range peers {
		fmt.Printf("  %d) %-20s %s\n", i+1, peer.Name, peer.Addr())
	}
	for {
		fmt.Print("Select peer: ")
		line, err := p.in.ReadString('\n')
		if err != nil {
			return domain.Peer{}, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || n < 1 || n > len(peers) {
			fmt.Println("invalid selection")
			continue
		}
		return peers[n-1], nil
	}
}
