// Package engine provides the primary Engine interface for querying narrative claims.
package engine

import (
	"errors"
	"fmt"

	"github.com/arminazimi/cinegraph/internal/narrative"
)

// ErrInvalidQuery is returned when a query is invalid.
var ErrInvalidQuery = errors.New("invalid query")

// TargetKind identifies what the query targets.
type TargetKind uint8

const (
	TargetKindWork TargetKind = iota + 1
	TargetKindEntity
)

// Target is a query target that selects exactly one of Work or Entity.
type Target struct {
	kind     TargetKind
	workID   narrative.WorkID
	entityID narrative.EntityID
}

// WorkTarget creates a target for a Work.
// It rejects empty IDs.
func WorkTarget(id narrative.WorkID) Target {
	if id == "" {
		return Target{kind: 0}
	}
	return Target{
		kind:   TargetKindWork,
		workID: id,
	}
}

// EntityTarget creates a target for an Entity.
// It rejects empty IDs.
func EntityTarget(id narrative.EntityID) Target {
	if id == "" {
		return Target{kind: 0}
	}
	return Target{
		kind:     TargetKindEntity,
		entityID: id,
	}
}

// Kind returns the target kind.
func (t Target) Kind() TargetKind {
	return t.kind
}

// WorkID returns the work ID if this is a Work target.
func (t Target) WorkID() narrative.WorkID {
	return t.workID
}

// EntityID returns the entity ID if this is an Entity target.
func (t Target) EntityID() narrative.EntityID {
	return t.entityID
}

// Query represents a structured query for claims.
type Query struct {
	Target           Target
	ContinuityFilter narrative.ContinuityID // empty means all continuities
}

// Validate checks that the query is well-formed.
func (q Query) Validate() error {
	if q.Target.kind == 0 {
		return fmt.Errorf("missing target: %w", ErrInvalidQuery)
	}
	switch q.Target.kind {
	case TargetKindWork:
		if q.Target.workID == "" {
			return fmt.Errorf("empty work id: %w", ErrInvalidQuery)
		}
	case TargetKindEntity:
		if q.Target.entityID == "" {
			return fmt.Errorf("empty entity id: %w", ErrInvalidQuery)
		}
	default:
		return fmt.Errorf("unknown target kind %d: %w", q.Target.kind, ErrInvalidQuery)
	}
	return nil
}
