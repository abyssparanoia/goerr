package goerr_test

import (
	"testing"

	"github.com/abyssparanoia/goerr"
)

func newErrorWithBuilder() *goerr.Error {
	return goerr.NewBuilder().With("color", "orange").New("error")
}

func TestBuilderNew(t *testing.T) {
	err := newErrorWithBuilder()

	if err.Values()["color"] != "orange" {
		t.Errorf("Unexpected value: %v", err.Values())
	}
}

func TestBuilderWrap(t *testing.T) {
	cause := goerr.New("cause")
	err := goerr.NewBuilder().With("color", "blue").Wrap(cause, "error")

	if err.Values()["color"] != "blue" {
		t.Errorf("Unexpected value: %v", err.Values())
	}

	if err.Unwrap().Error() != "cause" {
		t.Errorf("Unexpected cause: %v", err.Unwrap().Error())
	}
}

func TestBuilderWithMetadata(t *testing.T) {
	t.Run("metadata is set on New()", func(t *testing.T) {
		err := goerr.NewBuilder().
			With("internal", "secret").
			WithMetadata("invitation_id", "inv_456").
			New("builder error")

		meta := err.Metadata()
		if meta["invitation_id"] != "inv_456" {
			t.Errorf("expected invitation_id=inv_456, got %q", meta["invitation_id"])
		}
		if _, ok := meta["internal"]; ok {
			t.Errorf("internal should not appear in Metadata")
		}
	})

	t.Run("metadata is set on Wrap()", func(t *testing.T) {
		cause := goerr.New("cause")
		err := goerr.NewBuilder().
			WithMetadata("domain", "example.com").
			Wrap(cause, "wrapped")

		meta := err.Metadata()
		if meta["domain"] != "example.com" {
			t.Errorf("expected domain=example.com, got %q", meta["domain"])
		}
	})

	t.Run("Builder.WithMetadata does not mutate original Builder", func(t *testing.T) {
		base := goerr.NewBuilder().WithMetadata("base_key", "base_val")
		derived := base.WithMetadata("extra_key", "extra_val")

		baseErr := base.New("base")
		derivedErr := derived.New("derived")

		if _, ok := baseErr.Metadata()["extra_key"]; ok {
			t.Error("base Builder should not have extra_key")
		}
		if derivedErr.Metadata()["extra_key"] != "extra_val" {
			t.Errorf("derived should have extra_key, got %q", derivedErr.Metadata()["extra_key"])
		}
	})
}
