package engine_test

import (
	"errors"
	"testing"

	"github.com/arminazimi/cinegraph/internal/engine"
)

func TestQuery_TargetConstructors(t *testing.T) {
	t.Parallel()

	t.Run("work target validation", func(t *testing.T) {
		t.Parallel()
		q := engine.Query{Target: engine.WorkTarget("")}
		if err := q.Validate(); !errors.Is(err, engine.ErrInvalidQuery) {
			t.Fatalf("expected error for empty work id")
		}
	})

	t.Run("entity target validation", func(t *testing.T) {
		t.Parallel()
		q := engine.Query{Target: engine.EntityTarget("")}
		if err := q.Validate(); !errors.Is(err, engine.ErrInvalidQuery) {
			t.Fatalf("expected error for empty entity id")
		}
	})

	t.Run("missing target", func(t *testing.T) {
		t.Parallel()
		q := engine.Query{}
		if err := q.Validate(); !errors.Is(err, engine.ErrInvalidQuery) {
			t.Fatalf("expected error for missing target")
		}
	})
}
