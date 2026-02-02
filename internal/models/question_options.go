package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type QuestionOptions []Option

func (qo QuestionOptions) Value() (driver.Value, error) {
	return json.Marshal(qo)
}

func (qo *QuestionOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &qo)
}
