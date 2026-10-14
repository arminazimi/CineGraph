package narrative_test

import (
	"errors"
	"testing"

	"github.com/arminazimi/cinegraph/internal/narrative"
)

func TestClaim_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		claim   narrative.Claim
		wantErr bool
	}{
		{
			name: "valid claim",
			claim: narrative.Claim{
				ID:          "claim-1",
				Uncertainty: narrative.Certain,
				Evidence:    []narrative.Evidence{{ID: "evidence-1"}},
			},
		},
		{
			name: "nil evidence",
			claim: narrative.Claim{
				ID:          "claim-1",
				Uncertainty: narrative.Certain,
			},
			wantErr: true,
		},
		{
			name: "empty evidence",
			claim: narrative.Claim{
				ID:          "claim-1",
				Uncertainty: narrative.Certain,
				Evidence:    []narrative.Evidence{},
			},
			wantErr: true,
		},
		{
			name: "unsupported uncertainty",
			claim: narrative.Claim{
				ID:          "claim-1",
				Uncertainty: narrative.Uncertainty(255),
				Evidence:    []narrative.Evidence{{ID: "evidence-1"}},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.claim.Validate()
			if test.wantErr && !errors.Is(err, narrative.ErrInvalidClaim) {
				t.Fatalf("Validate() error = %v, want ErrInvalidClaim", err)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}
		})
	}
}
