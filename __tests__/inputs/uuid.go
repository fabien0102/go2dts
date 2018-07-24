package client

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func nonEmptyUUIDRule(value interface{}) error {
	switch v := value.(type) {
	case uuid.UUID:
		if v == uuid.Nil {
			return errors.New("cannot be blank")
		}
		return nil
	default:
		return errors.New("unable to parse uuid")
	}
}
