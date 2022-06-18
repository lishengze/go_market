package net

import "fmt"

type WSInfo struct {
	ID int64
}

func (s *WSInfo) String() string {
	return fmt.Sprintf("ID: %d ", s.ID)
}
