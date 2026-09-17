package main

import (
	"errors"
	"testing"
)

// TestSplitUTXOsIfNeeded_CallsSubmitExactlyOnce is the regression test for finding 7: a
// historical copy-paste bug called SubmitCoinSplitRequest twice in a row with the identical
// computed amount, silently doubling the real split amount relative to what was logged. This
// test proves splitUTXOsIfNeeded (and therefore checkAndSplitUTXOs, which delegates to it)
// invokes the submit callback exactly once per invocation when a split is needed.
func TestSplitUTXOsIfNeeded_CallsSubmitExactlyOnce(t *testing.T) {
	var calls int
	var gotSplitAmt, gotNumSplits int

	count, err := splitUTXOsIfNeeded(10, 400, 8_000_000_000, func(splitAmt int, numSplits int) error {
		calls++
		gotSplitAmt = splitAmt
		gotNumSplits = numSplits
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected submit to be called exactly once, got %v calls", calls)
	}
	if count != 1 {
		t.Fatalf("expected splitUTXOsIfNeeded to report 1 call, got %v", count)
	}

	// balance / 10 / (minCount * 2) = 8_000_000_000 / 10 / 800 = 1_000_000
	wantSplitAmt := 1_000_000
	if gotSplitAmt != wantSplitAmt {
		t.Errorf("splitAmt = %v, want %v", gotSplitAmt, wantSplitAmt)
	}
	if gotNumSplits != 400 {
		t.Errorf("numSplits = %v, want 400", gotNumSplits)
	}
}

// TestSplitUTXOsIfNeeded_NoOpWhenAboveMinCount verifies that no split call is made at all when
// the existing UTXO count already meets the minimum.
func TestSplitUTXOsIfNeeded_NoOpWhenAboveMinCount(t *testing.T) {
	var calls int
	count, err := splitUTXOsIfNeeded(400, 400, 8_000_000_000, func(splitAmt int, numSplits int) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 0 || count != 0 {
		t.Fatalf("expected no submit calls when existingCount >= minCount, got calls=%v count=%v", calls, count)
	}
}

// TestSplitUTXOsIfNeeded_PropagatesSubmitError verifies that an error from the submit callback
// is propagated, and that the call count still reflects exactly what happened.
func TestSplitUTXOsIfNeeded_PropagatesSubmitError(t *testing.T) {
	wantErr := errors.New("boom")
	calls := 0
	count, err := splitUTXOsIfNeeded(0, 400, 8_000_000_000, func(splitAmt int, numSplits int) error {
		calls++
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
	if calls != 1 || count != 1 {
		t.Fatalf("expected exactly one attempted call, got calls=%v count=%v", calls, count)
	}
}

// TestSplitCallCount pins the regression-guard constant itself.
func TestSplitCallCount(t *testing.T) {
	if splitCallCount != 1 {
		t.Fatalf("splitCallCount must be 1 (one split submission per invocation), got %v", splitCallCount)
	}
}
