package consts

import "database/sql/driver"

type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
)

func (ct *Gender) Scan(value interface{}) error {
	*ct = Gender(value.([]byte))
	return nil
}

func (ct Gender) Value() (driver.Value, error) {
	return string(ct), nil
}
