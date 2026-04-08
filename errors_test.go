package goerr_test

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"testing"

	"github.com/abyssparanoia/goerr"
)

func oops() *goerr.Error {
	return goerr.New("omg")
}

func normalError() error {
	return fmt.Errorf("red")
}

func wrapError() *goerr.Error {
	err := normalError()
	return goerr.Wrap(err, "orange")
}

func TestNew(t *testing.T) {
	err := oops()
	v := fmt.Sprintf("%+v", err)
	if !strings.Contains(v, "goerr_test.oops") {
		t.Error("Stack trace 'goerr_test.oops' is not found")
	}
	if !strings.Contains(err.Error(), "omg") {
		t.Error("Error message is not correct")
	}
}

func TestWrapError(t *testing.T) {
	err := wrapError()
	st := fmt.Sprintf("%+v", err)
	if !strings.Contains(st, "github.com/abyssparanoia/goerr_test.wrapError") {
		t.Error("Stack trace 'wrapError' is not found")
	}
	if !strings.Contains(st, "github.com/abyssparanoia/goerr_test.TestWrapError") {
		t.Error("Stack trace 'TestWrapError' is not found")
	}
	if strings.Contains(st, "github.com/abyssparanoia/goerr_test.normalError") {
		t.Error("Stack trace 'normalError' is found")
	}
	if !strings.Contains(err.Error(), "orange: red") {
		t.Error("Error message is not correct")
	}
}

func TestStackTrace(t *testing.T) {
	err := oops()
	st := err.Stacks()
	if len(st) != 4 {
		t.Errorf("Expected stack length of 4, got %d", len(st))
	}
	if st[0].Func != "github.com/abyssparanoia/goerr_test.oops" {
		t.Error("Stack trace 'github.com/abyssparanoia/goerr_test.oops' is not found")
	}
	if !regexp.MustCompile(`/goerr/errors_test\.go$`).MatchString(st[0].File) {
		t.Error("Stack trace file is not correct")
	}
	if st[0].Line != 16 {
		t.Errorf("Expected line number 13, got %d", st[0].Line)
	}
}

func TestMultiWrap(t *testing.T) {
	err1 := oops()
	err2 := goerr.Wrap(err1)
	if err1 == err2 {
		t.Error("Expected err1 and err2 to be different")
	}

	err3 := goerr.Wrap(err1, "some message")
	if err1 == err3 {
		t.Error("Expected err1 and err3 to be different")
	}
}

func TestErrorCode(t *testing.T) {
	rootErr := goerr.New("something bad")
	baseErr1 := goerr.New("oops").ID("code1")
	baseErr2 := goerr.New("oops").ID("code2")

	newErr := baseErr1.Wrap(rootErr).WithValue("v", 1)

	if !errors.Is(newErr, baseErr1) {
		t.Error("Expected newErr to be based on baseErr1")
	}
	if newErr == baseErr1 {
		t.Error("Expected newErr and baseErr1 to be different")
	}
	if newErr.Values()["v"] == nil {
		t.Error("Expected newErr to have a non-nil value for 'v'")
	}
	if baseErr1.Values()["v"] != nil {
		t.Error("Expected baseErr1 to have a nil value for 'v'")
	}
	if errors.Is(newErr, baseErr2) {
		t.Error("Expected newErr to not be based on baseErr2")
	}
}

func TestPrintable(t *testing.T) {
	cause := errors.New("test")
	err := goerr.Wrap(cause, "oops").ID("E001").WithValue("blue", "five")

	p := err.Printable()
	if p.Message != "oops" {
		t.Errorf("Expected message to be 'oops', got '%s'", p.Message)
	}
	if p.ID != "E001" {
		t.Errorf("Expected ID to be 'E001', got '%s'", p.ID)
	}
	if p.Cause != cause {
		t.Errorf("Expected cause to be '%v', got '%v'", cause, p.Cause)
	}
	if p.Values["blue"] != "five" {
		t.Errorf("Expected value for 'blue' to be 'five', got '%v'", p.Values["blue"])
	}
}

func TestUnwrap(t *testing.T) {
	err1 := goerr.New("omg").WithValue("color", "five")
	err2 := fmt.Errorf("oops: %w", err1)

	err := goerr.Unwrap(err2)
	if err == nil {
		t.Error("Expected unwrapped error to be non-nil")
	}
	values := err.Values()
	if values["color"] != "five" {
		t.Errorf("Expected value for 'color' to be 'five', got '%v'", values["color"])
	}
}

func TestFormat(t *testing.T) {
	err := goerr.New("test: %s", "blue")
	if err.Error() != "test: blue" {
		t.Errorf("Expected error message to be 'test: blue', got '%s'", err.Error())
	}
}

func TestErrorString(t *testing.T) {
	err := goerr.Wrap(goerr.Wrap(goerr.New("blue"), "orange"), "red")
	if err.Error() != "red: orange: blue" {
		t.Errorf("Expected error message to be 'red: orange: blue', got '%s'", err.Error())
	}
}

func TestLoggingNestedError(t *testing.T) {
	err1 := goerr.New("e1").WithValue("color", "orange")
	err2 := goerr.Wrap(err1, "e2").WithValue("number", "five")
	out := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(out, nil))
	logger.Error("fail", slog.Any("error", err2))
	if !strings.Contains(out.String(), `"number":"five"`) {
		t.Errorf("Expected log output to contain '\"number\":\"five\"', got '%s'", out.String())
	}
	if !strings.Contains(out.String(), `"color":"orange"`) {
		t.Errorf("Expected log output to contain '\"color\":\"orange\"', got '%s'", out.String())
	}
}

func TestLoggerWithNil(t *testing.T) {
	out := &bytes.Buffer{}
	var err *goerr.Error
	logger := slog.New(slog.NewJSONHandler(out, nil))
	logger.Error("fail", slog.Any("error", err))
	if !strings.Contains(out.String(), `"error":null`) {
		t.Errorf("Expected log output to contain '\"error\":null', got '%s'", out.String())
	}
}

func TestUnstack(t *testing.T) {
	t.Run("original stack", func(t *testing.T) {
		err := oops()
		st := err.Stacks()
		if st == nil {
			t.Error("Expected stack trace to be nil")
		}
		if len(st) == 0 {
			t.Error("Expected stack trace length to be 0")
		}
		if st[0].Func != "github.com/abyssparanoia/goerr_test.oops" {
			t.Errorf("Not expected stack trace func name (github.com/abyssparanoia/goerr_test.oops): %s", st[0].Func)
		}
	})

	t.Run("unstacked", func(t *testing.T) {
		err := oops().Unstack()
		st1 := err.Stacks()
		if st1 == nil {
			t.Error("Expected stack trace to be non-nil")
		}
		if len(st1) == 0 {
			t.Error("Expected stack trace length to be non-zero")
		}
		if st1[0].Func != "github.com/abyssparanoia/goerr_test.TestUnstack.func2" {
			t.Errorf("Not expected stack trace func name (github.com/abyssparanoia/goerr_test.TestUnstack.func2): %s", st1[0].Func)
		}
	})

	t.Run("unstackN with 2", func(t *testing.T) {
		err := oops().UnstackN(2)
		st2 := err.Stacks()
		if st2 == nil {
			t.Error("Expected stack trace to be non-nil")
		}
		if len(st2) == 0 {
			t.Error("Expected stack trace length to be non-zero")
		}
		if st2[0].Func != "testing.tRunner" {
			t.Errorf("Not expected stack trace func name (testing.tRunner): %s", st2[0].Func)
		}
	})
}

func TestWithMetadata(t *testing.T) {
	t.Run("stores metadata separately from values", func(t *testing.T) {
		err := goerr.New("test error").
			WithValue("internal_id", "secret").
			WithMetadata("invitation_id", "inv_123").
			WithMetadata("email_domain", "example.com")

		meta := err.Metadata()
		if meta["invitation_id"] != "inv_123" {
			t.Errorf("expected invitation_id=inv_123, got %q", meta["invitation_id"])
		}
		if meta["email_domain"] != "example.com" {
			t.Errorf("expected email_domain=example.com, got %q", meta["email_domain"])
		}

		// WithValue must NOT bleed into Metadata
		if _, ok := meta["internal_id"]; ok {
			t.Errorf("internal_id should not appear in Metadata, but it does")
		}

		// Metadata must NOT bleed into Values
		vals := err.Values()
		if _, ok := vals["invitation_id"]; ok {
			t.Errorf("invitation_id should not appear in Values, but it does")
		}
	})

	t.Run("merges metadata across Wrap chain (outer wins)", func(t *testing.T) {
		inner := goerr.New("inner").WithMetadata("key", "inner_val").WithMetadata("only_inner", "yes")
		outer := goerr.New("outer").WithMetadata("key", "outer_val").Wrap(inner)

		meta := outer.Metadata()
		if meta["key"] != "outer_val" {
			t.Errorf("expected outer to win: got %q", meta["key"])
		}
		if meta["only_inner"] != "yes" {
			t.Errorf("expected only_inner from wrapped chain: got %q", meta["only_inner"])
		}
	})

	t.Run("New() on template copies metadata", func(t *testing.T) {
		template := goerr.New("template").WithMetadata("reason", "template_reason")
		instance := template.New()

		meta := instance.Metadata()
		if meta["reason"] != "template_reason" {
			t.Errorf("expected reason=template_reason, got %q", meta["reason"])
		}
	})

	t.Run("empty metadata returns empty map (not nil)", func(t *testing.T) {
		err := goerr.New("no metadata")
		meta := err.Metadata()
		if meta == nil {
			t.Error("expected non-nil map")
		}
		if len(meta) != 0 {
			t.Errorf("expected empty map, got %v", meta)
		}
	})
}
