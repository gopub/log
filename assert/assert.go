package assert

import (
	"github.com/gopub/log"
	"github.com/stretchr/testify/require"
)

var Logger = log.Default()

type testingT struct {
}

var tt require.TestingT = testingT{}

func (t testingT) Errorf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

func (t testingT) FailNow() {

}

func NoError(err error) {
	require.NoError(tt, err)
}

func False(val bool, msgAndArgs ...interface{}) {
	require.False(tt, val, msgAndArgs...)
}

func True(val bool, msgAndArgs ...interface{}) {
	require.True(tt, val, msgAndArgs...)
}

func NotEmpty(val interface{}, msgAndArgs ...interface{}) {
	require.NotEmpty(tt, val, msgAndArgs...)
}

func Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	require.Equal(tt, expected, actual, msgAndArgs...)
}
