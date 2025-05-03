package service

import "context"

type StoragePinger interface {
	Ping(context.Context) error
}

type Ping struct {
	sping StoragePinger
}

func NewPing(sping StoragePinger) *Ping {
	return &Ping{
		sping: sping,
	}
}

func (s *Ping) PingStorage(ctx context.Context) error {
	err := s.sping.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
