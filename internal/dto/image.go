package dto

import "io"

type Image struct {
	ID        int
	UserID    int
	Name      string
	Extension string
	Data      io.Reader
}
