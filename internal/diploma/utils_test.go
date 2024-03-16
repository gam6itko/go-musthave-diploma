package diploma

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LuhnValidate(t *testing.T) {
	//assert.True(t, LuhnValidate(2022))
	assert.True(t, LuhnValidate(4561_2612_1234_5467))
	assert.True(t, LuhnValidate(5580_4733_7202_4733))

	assert.False(t, LuhnValidate(5580_4733_7202_4732))
}
