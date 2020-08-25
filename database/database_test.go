package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmployeeIDReturnsErr(t *testing.T) {
	e := Employee{
		ID: 55,
	}
	err := Validate(e)
	assert.Len(t, err, 1)
}

func TestEmployeeEmailReturnsErr(t *testing.T) {
	e := Employee{
		Email: "abc.com",
	}
	err := Validate(e)
	assert.Len(t, err, 1)
}

func TestEmployeeDetailsReturnsErr(t *testing.T) {
	e := Employee{
		ID:    345,
		Name:  "Michale",
		Email: "abc@abc.com",
	}
	err := Validate(e)
	assert.Len(t, err, 1)
}
