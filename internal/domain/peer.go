package domain

import "fmt"

type Peer struct {
	Name string
	Host string
	Port int
}

func (p Peer) Addr() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}
