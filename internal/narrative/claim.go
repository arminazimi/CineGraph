// Package narrative defines core domain types for narrative knowledge.
package narrative

import (
	"errors"
	"fmt"
)

// ErrInvalidClaim identifies a Claim that violates a domain invariant.
var ErrInvalidClaim = errors.New("invalid claim")

// ID types for narrative domain.
type WorkID string
type EntityID string
type ClaimID string
type EvidenceID string
type ContinuityID string

// Uncertainty represents the confidence/certainty of a claim.
type Uncertainty uint8

const (
	Uncertain Uncertainty = iota // 0
	Certain
	Probable
	Disputed
)

// String returns a string representation of Uncertainty.
func (u Uncertainty) String() string {
	switch u {
	case Certain:
		return "certain"
	case Probable:
		return "probable"
	case Uncertain:
		return "uncertain"
	case Disputed:
		return "disputed"
	default:
		return fmt.Sprintf("unknown_uncertainty_%d", u)
	}
}

// Validate ensures the uncertainty value is one of the supported values.
func (u Uncertainty) Validate() error {
	switch u {
	case Certain, Probable, Uncertain, Disputed:
		return nil
	default:
		return fmt.Errorf("invalid uncertainty: %d", u)
	}
}

// Evidence represents a traceable source supporting or challenging a claim.
type Evidence struct {
	ID      EvidenceID
	Source  string // source identifier or source URI
	Locator string // optional locator within the source
}

// Claim represents a sourced assertion about narrative fact.
type Claim struct {
	ID           ClaimID
	Assertion    string
	ContinuityID ContinuityID
	Uncertainty  Uncertainty
	Evidence     []Evidence
}

// Validate checks the invariants required for a Claim to enter the Engine.
func (c Claim) Validate() error {
	if len(c.Evidence) == 0 {
		return fmt.Errorf("claim %q has no evidence: %w", c.ID, ErrInvalidClaim)
	}
	if err := c.Uncertainty.Validate(); err != nil {
		return fmt.Errorf("claim %q has %v: %w", c.ID, err, ErrInvalidClaim)
	}
	return nil
}
