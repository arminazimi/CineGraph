package engine_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/arminazimi/cinegraph/internal/engine"
	"github.com/arminazimi/cinegraph/internal/narrative"
)

type fakeClaimSource struct {
	calls           int
	lastQuery       engine.Query
	result          []narrative.Claim
	err             error
	receivedContext context.Context
}

func (f *fakeClaimSource) FindClaims(ctx context.Context, q engine.Query) ([]narrative.Claim, error) {
	f.calls++
	f.lastQuery = q
	f.receivedContext = ctx
	if f.err != nil {
		return nil, f.err
	}
	return f.result, nil
}

type valueClaimSource struct{}

func (valueClaimSource) FindClaims(context.Context, engine.Query) ([]narrative.Claim, error) {
	return []narrative.Claim{}, nil
}

func TestNew(t *testing.T) {
	t.Parallel()

	var typedNil *fakeClaimSource
	tests := []struct {
		name    string
		source  engine.ClaimSource
		wantErr bool
	}{
		{name: "nil interface", source: nil, wantErr: true},
		{name: "typed nil pointer", source: typedNil, wantErr: true},
		{name: "pointer source", source: &fakeClaimSource{}},
		{name: "value source", source: valueClaimSource{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := engine.New(test.source)
			if test.wantErr {
				if !errors.Is(err, engine.ErrNilSource) {
					t.Fatalf("New() error = %v, want ErrNilSource", err)
				}
				if got != nil {
					t.Fatalf("New() engine = %v, want nil", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("New() returned a nil Engine")
			}
		})
	}
}

func newEngine(t *testing.T, source engine.ClaimSource) *engine.Engine {
	t.Helper()
	eng, err := engine.New(source)
	if err != nil {
		t.Fatalf("creating engine: %v", err)
	}
	return eng
}

func TestEngine_ValidTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target engine.Target
	}{
		{name: "work target", target: engine.WorkTarget("w1")},
		{name: "entity target", target: engine.EntityTarget("ent1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			claims := []narrative.Claim{
				{
					ID:           "c1",
					Assertion:    "s",
					ContinuityID: "A",
					Uncertainty:  narrative.Certain,
					Evidence:     []narrative.Evidence{{ID: "e1"}},
				},
			}
			fake := &fakeClaimSource{result: claims}
			eng := newEngine(t, fake)
			query := engine.Query{Target: tt.target}
			res, err := eng.Query(context.Background(), query)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			if fake.calls != 1 {
				t.Fatalf("expected 1 call")
			}
			if fake.lastQuery != query {
				t.Fatalf("source query = %#v, want %#v", fake.lastQuery, query)
			}
			if len(res.Claims) != 1 {
				t.Fatalf("expected 1 claim")
			}
		})
	}
}

func TestEngine_MissingAndEmptyTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		query engine.Query
	}{
		{name: "missing target", query: engine.Query{}},
		{name: "empty work id", query: engine.Query{Target: engine.WorkTarget("")}},
		{name: "empty entity id", query: engine.Query{Target: engine.EntityTarget("")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fake := &fakeClaimSource{}
			eng := newEngine(t, fake)
			_, err := eng.Query(context.Background(), tt.query)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !errors.Is(err, engine.ErrInvalidQuery) {
				t.Fatalf("expected ErrInvalidQuery")
			}
			if fake.calls != 0 {
				t.Fatalf("source called %d times for invalid query", fake.calls)
			}
		})
	}
}

func TestEngine_ContinuityFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		filter           narrative.ContinuityID
		expectedCount    int
		expectedClaimIDs []narrative.ClaimID
	}{
		{name: "filters other continuities", filter: "A", expectedCount: 1, expectedClaimIDs: []narrative.ClaimID{"c1"}},
		{name: "empty allows all", filter: "", expectedCount: 2, expectedClaimIDs: []narrative.ClaimID{"c1", "c2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			claims := []narrative.Claim{
				{
					ID: "c1", Assertion: "s1", ContinuityID: "A", Uncertainty: narrative.Certain,
					Evidence: []narrative.Evidence{{ID: "e1"}},
				},
				{
					ID: "c2", Assertion: "s2", ContinuityID: "B", Uncertainty: narrative.Certain,
					Evidence: []narrative.Evidence{{ID: "e1"}},
				},
			}
			fake := &fakeClaimSource{result: claims}
			eng := newEngine(t, fake)
			res, err := eng.Query(
				context.Background(),
				engine.Query{
					Target:           engine.WorkTarget("w1"),
					ContinuityFilter: tt.filter,
				},
			)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			if len(res.Claims) != tt.expectedCount {
				t.Fatalf("expected %d claims, got %d", tt.expectedCount, len(res.Claims))
			}
			for i, wantID := range tt.expectedClaimIDs {
				if res.Claims[i].ID != wantID {
					t.Fatalf("claim %d ID = %q, want %q", i, res.Claims[i].ID, wantID)
				}
			}
		})
	}
}

func TestEngine_SupportedUncertainties(t *testing.T) {
	t.Parallel()

	uncs := []narrative.Uncertainty{narrative.Certain, narrative.Probable, narrative.Uncertain, narrative.Disputed}
	for _, u := range uncs {
		t.Run(u.String(), func(t *testing.T) {
			t.Parallel()
			claims := []narrative.Claim{
				{ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: u, Evidence: []narrative.Evidence{{ID: "e1"}}},
			}
			fake := &fakeClaimSource{result: claims}
			eng := newEngine(t, fake)
			res, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")})
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			if res.Claims[0].Uncertainty != u {
				t.Fatalf("uncertainty not retained")
			}
		})
	}
}

func TestEngine_DeterministicOrdering(t *testing.T) {
	t.Parallel()

	t.Run("by continuity then claim id", func(t *testing.T) {
		t.Parallel()
		claims := []narrative.Claim{
			{
				ID: "c3", Assertion: "s3", ContinuityID: "B", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{{ID: "e1"}},
			},
			{
				ID: "c1", Assertion: "s1", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{{ID: "e1"}},
			},
			{
				ID: "c2", Assertion: "s2", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{{ID: "e1"}},
			},
		}
		fake := &fakeClaimSource{result: claims}
		eng := newEngine(t, fake)
		res, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")})
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if res.Claims[0].ID != "c1" || res.Claims[1].ID != "c2" || res.Claims[2].ID != "c3" {
			t.Fatalf("wrong order")
		}
	})

	t.Run("evidence ordered by id", func(t *testing.T) {
		t.Parallel()
		claims := []narrative.Claim{
			{
				ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{{ID: "e3"}, {ID: "e1"}, {ID: "e2"}},
			},
		}
		fake := &fakeClaimSource{result: claims}
		eng := newEngine(t, fake)
		res, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")})
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		got := res.Claims[0].Evidence
		if got[0].ID != "e1" || got[1].ID != "e2" || got[2].ID != "e3" {
			t.Fatalf("wrong evidence order")
		}
	})

	t.Run("evidence stable with equal ids", func(t *testing.T) {
		t.Parallel()
		claims := []narrative.Claim{
			{
				ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{
					{ID: "e1", Source: "srcA", Locator: "loc1"},
					{ID: "e1", Source: "srcB", Locator: "loc2"},
				},
			},
		}
		fake := &fakeClaimSource{result: claims}
		eng := newEngine(t, fake)
		res, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")})
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if res.Claims[0].Evidence[0].Source != "srcA" || res.Claims[0].Evidence[1].Source != "srcB" {
			t.Fatalf("stable order not preserved")
		}
	})
}

func TestEngine_SourceFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sourceErr error
	}{
		{name: "sentinel error", sourceErr: errors.New("boom")},
		{name: "canceled context", sourceErr: context.Canceled},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fake := &fakeClaimSource{err: test.sourceErr}
			eng := newEngine(t, fake)
			_, err := eng.Query(
				context.Background(),
				engine.Query{Target: engine.WorkTarget("w1")},
			)
			if !errors.Is(err, test.sourceErr) {
				t.Fatalf("Query() error = %v, want wrapped %v", err, test.sourceErr)
			}
			if !strings.Contains(err.Error(), "querying claims") {
				t.Fatalf("Query() error = %q, want operation context", err)
			}
		})
	}
}

func TestEngine_NoMatchEmptySlice(t *testing.T) {
	t.Parallel()

	fake := &fakeClaimSource{result: []narrative.Claim{}}
	eng := newEngine(t, fake)
	res, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if res.Claims == nil {
		t.Fatalf("expected non-nil empty slice")
	}
}

func TestEngine_DoesNotMutateSource(t *testing.T) {
	t.Parallel()

	claims := []narrative.Claim{
		{
			ID: "c2", Assertion: "s2", ContinuityID: "B", Uncertainty: narrative.Certain,
			Evidence: []narrative.Evidence{{ID: "e1"}},
		},
		{
			ID: "c1", Assertion: "s1", ContinuityID: "A", Uncertainty: narrative.Certain,
			Evidence: []narrative.Evidence{{ID: "e2"}, {ID: "e1"}},
		},
	}
	want := cloneClaims(claims)
	fake := &fakeClaimSource{result: claims}
	eng := newEngine(t, fake)
	if _, err := eng.Query(context.Background(), engine.Query{Target: engine.WorkTarget("w1")}); err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if !reflect.DeepEqual(claims, want) {
		t.Fatalf("source claims mutated:\n got: %#v\nwant: %#v", claims, want)
	}
}

func TestEngine_ContextPropagation(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}
	key := ctxKey{}
	ctx := context.WithValue(context.Background(), key, "v")
	claims := []narrative.Claim{{
		ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
		Evidence: []narrative.Evidence{{ID: "e1"}},
	}}
	fake := &fakeClaimSource{result: claims}
	eng := newEngine(t, fake)
	if _, err := eng.Query(ctx, engine.Query{Target: engine.WorkTarget("w1")}); err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if fake.receivedContext != ctx {
		t.Fatalf("context not propagated")
	}
}

func TestEngine_EvidenceInvariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		claims  []narrative.Claim
		wantErr bool
	}{
		{
			name: "nil evidence rejected",
			claims: []narrative.Claim{{
				ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
			}},
			wantErr: true,
		},
		{
			name: "empty evidence rejected",
			claims: []narrative.Claim{{
				ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{},
			}},
			wantErr: true,
		},
		{
			name: "evidence retained",
			claims: []narrative.Claim{{
				ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
				Evidence: []narrative.Evidence{{ID: "e1"}},
			}},
		},
		{
			name: "invalid filtered claim rejected",
			claims: []narrative.Claim{
				{
					ID: "c1", Assertion: "s", ContinuityID: "A", Uncertainty: narrative.Certain,
					Evidence: []narrative.Evidence{{ID: "e1"}},
				},
				{ID: "c2", Assertion: "s", ContinuityID: "B", Uncertainty: narrative.Certain},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fake := &fakeClaimSource{result: test.claims}
			eng := newEngine(t, fake)
			result, err := eng.Query(
				context.Background(),
				engine.Query{
					Target:           engine.WorkTarget("w1"),
					ContinuityFilter: "A",
				},
			)
			if test.wantErr {
				if !errors.Is(err, narrative.ErrInvalidClaim) {
					t.Fatalf("Query() error = %v, want ErrInvalidClaim", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Query() unexpected error: %v", err)
			}
			if len(result.Claims) != 1 || len(result.Claims[0].Evidence) != 1 {
				t.Fatalf("Query() result = %#v, want one evidence-backed claim", result)
			}
		})
	}
}

func cloneClaims(claims []narrative.Claim) []narrative.Claim {
	cloned := make([]narrative.Claim, len(claims))
	copy(cloned, claims)
	for i := range cloned {
		cloned[i].Evidence = make([]narrative.Evidence, len(claims[i].Evidence))
		copy(cloned[i].Evidence, claims[i].Evidence)
	}
	return cloned
}
