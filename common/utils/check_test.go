package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyEmailFormat(t *testing.T) {

	assert.Equal(t, true, func() bool {
		return VerifyEmailFormat("wanghyu@163.com")
	}())

	assert.Equal(t, false, func() bool {
		return VerifyEmailFormat("wanghyu.com")
	}())
}

func TestVerifyMobileFormat(t *testing.T) {
	assert.Equal(t, true, func() bool {
		return VerifyMobileFormat("1572314")
	}())

	assert.Equal(t, false, func() bool {
		return VerifyMobileFormat("")
	}())
}
