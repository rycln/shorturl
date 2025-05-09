package service

import "context"

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type PingStorage interface {
	Ping(context.Context) error
}

type Ping struct {
	strg PingStorage
}

func NewPing(strg PingStorage) *Ping {
	return &Ping{
		strg: strg,
	}
}

func (s *Ping) PingStorage(ctx context.Context) error {
	err := s.strg.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
