package uuidGenerator

import "github.com/google/uuid"

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (s UUIDGenerator) Generate() (uuid.UUID, error) {
	uuidNew, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuidNew, nil
}
