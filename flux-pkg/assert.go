package fluxpkg

import "fmt"

const (
	assertMessagePrefix = "FXSERVER:CRITICAL:ASSERT:"
)

func AssertT(assert func() bool, message string) {
	if !assert() {
		panic(assertMessagePrefix + message)
	}
}

func Assert(true bool, message string) {
	if !true {
		panic(assertMessagePrefix + message)
	}
}

func AssertNil(v interface{}, message string, args ...interface{}) {
	if nil != v {
		panic(assertMessagePrefix + fmt.Sprintf(message, args...))
	}
}

func AssertNotNil(v interface{}, message string, args ...interface{}) {
	if nil == v {
		panic(assertMessagePrefix + fmt.Sprintf(message, args...))
	}
}
