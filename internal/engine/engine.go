// Package engine provides the primary Engine interface for querying narrative claims.
package engine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/arminazimi/cinegraph/internal/narrative"
)

// ErrNilSource is returned when a nil ClaimSource is provided.
var ErrNilSource = errors.New("nil claim source")

// ClaimSource is the dependency interface through which the Engine obtains claims.
type ClaimSource interface {
	FindClaims(
		ctx context.Context,
		query Query,
	) ([]narrative.Claim, error)
}

// Result contains the claims returned by a query.
type Result struct {
	Claims []narrative.Claim
}

// Engine is the primary interface used by callers and tests.
type Engine struct {
	source ClaimSource
}

// New creates a new Engine with the given source.
// It rejects nil sources.
func New(source ClaimSource) (*Engine, error) {
	if isNilClaimSource(source) {
		return nil, ErrNilSource
	}
	return &Engine{source: source}, nil
}

func isNilClaimSource(source ClaimSource) bool {
	if source == nil {
		return true
	}

	value := reflect.ValueOf(source)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// Query executes a structured query and returns normalized results.
func (e *Engine) Query(ctx context.Context, query Query) (Result, error) {
	if err := query.Validate(); err != nil {
		return Result{}, err
	}

	candidates, err := e.source.FindClaims(ctx, query)
	if err != nil {
		return Result{}, fmt.Errorf("querying claims: %w", err)
	}
	for _, candidate := range candidates {
		if err := candidate.Validate(); err != nil {
			return Result{}, fmt.Errorf("validating source claim: %w", err)
		}
	}

	// Make a copy so we don't mutate source-owned slices
	result := make([]narrative.Claim, len(candidates))
	copy(result, candidates)

	// Apply continuity filter if specified
	if query.ContinuityFilter != "" {
		filtered := make([]narrative.Claim, 0, len(result))
		for i := range result {
			if result[i].ContinuityID == query.ContinuityFilter {
				filtered = append(filtered, result[i])
			}
		}
		result = filtered
	}

	// Copy evidence before normalizing caller-visible order.
	for i := range result {
		ev := make([]narrative.Evidence, len(result[i].Evidence))
		copy(ev, result[i].Evidence)
		result[i].Evidence = ev
	}

	// Sort evidence stably by ID so equal IDs retain their source order.
	for i := range result {
		sort.SliceStable(result[i].Evidence, func(a, b int) bool {
			return result[i].Evidence[a].ID < result[i].Evidence[b].ID
		})
	}

	// Sort claims deterministically: continuity ID, then claim ID
	// Use stable sort to preserve original relative order on ties if any
	sort.SliceStable(result, func(a, b int) bool {
		if result[a].ContinuityID != result[b].ContinuityID {
			return result[a].ContinuityID < result[b].ContinuityID
		}
		return result[a].ID < result[b].ID
	})

	return Result{Claims: result}, nil
}
