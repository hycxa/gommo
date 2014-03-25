package god

import (
	"fmt"
	"hash"
)

type UUID struct {
	hash.Hash
}

func (id UUID) String() string {
	return fmt.Sprintf("%x", id.Sum(nil))
}
