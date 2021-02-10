package util

import (
	"github.com/google/uuid"
	"strconv"
)

type UUIDGenerator struct {
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) Generate() string {
	return uuid.New().String()
}

type SequentialIDGenerator struct {
	nextID int
}

func NewSequentialIDGenerator() *SequentialIDGenerator {
	return &SequentialIDGenerator{}
}

func (s *SequentialIDGenerator) Generate() string {
	id := strconv.Itoa(s.nextID)
	s.nextID++
	return id
}
