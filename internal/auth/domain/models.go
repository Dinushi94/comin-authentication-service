// internal/auth/domain/models.go
package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func (s *Settings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &s)
}

// Implement driver.Valuer interface
func (s Settings) Value() (driver.Value, error) {
	return json.Marshal(s)
}
