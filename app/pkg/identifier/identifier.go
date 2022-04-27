package identifier

import (
	"fmt"

	"github.com/rs/xid"
)

type ResourceType interface{ Prefix() string }

type ID[T ResourceType] xid.ID

func New[T ResourceType]() ID[T] {
	return ID[T](xid.New())
}

func (id ID[T]) String() string {
	var resourceType T
	return fmt.Sprintf("%s_%s", resourceType.Prefix(), xid.ID(id).String())
}

func (id ID[T]) FromString(idString string) ID[T] {
	value, err := xid.FromString(idString)
	if err != nil {
		fmt.Errorf("failed to parse id %v", err.Error())
	}

	return ID[T](value)
}
