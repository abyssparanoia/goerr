package main

import (
	"github.com/abyssparanoia/goerr"
	"go.uber.org/zap"
)

var (
	testErr = goerr.New("test error").
		WithCategory("test_category").
		WithCode("TEST123").
		WithDetail("test detail").
		WithValue("key", "value")
)

func main() {
	l, _ := zap.NewDevelopment()
	if err := wrapDoSomething(); err != nil {
		l.Error("error", goerr.ZapError(err))
	}
}

func wrapDoSomething() error {
	if err := doSomething(); err != nil {
		return err
	}
	return nil
}

func doSomething() error {
	err := testErr.Errorf("failed to do something").
		WithValue("other_key", "other_value")
	return err
}
