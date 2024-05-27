package token

import "time"

// for making tokens using diff algorithms
type Maker interface {
	// creates a new token
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// validates a token
	VerifyToken(token string) (*Payload, error)
}
