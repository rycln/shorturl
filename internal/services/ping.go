package services

import "context"

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type PingStorage interface {
	Ping(context.Context) error
}

// Ping provides storage health-check functionality.
//
// The service verifies connectivity and readiness of the underlying storage system.
type Ping struct {
	strg PingStorage
}

// NewPing creates a new Ping service instance.
func NewPing(strg PingStorage) *Ping {
	return &Ping{
		strg: strg,
	}
}

// Check verifies storage connectivity.
//
// Returns nil if storage is responsive and ready to handle requests.
func (s *Ping) PingStorage(ctx context.Context) error {
	err := s.strg.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
