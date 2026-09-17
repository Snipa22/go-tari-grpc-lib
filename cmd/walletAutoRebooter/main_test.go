package main

import "testing"

func TestShouldCheckHeightBasedReboot(t *testing.T) {
	cases := []struct {
		name          string
		scannedHeight uint64
		bestHeight    uint64
		hasResponded  bool
		wantShouldRun bool
	}{
		{
			name:          "zero peers responded does not underflow into an unconditional reboot",
			scannedHeight: 1000,
			bestHeight:    0,
			hasResponded:  false,
			wantShouldRun: false,
		},
		{
			name:          "peer responded and wallet scan is genuinely behind",
			scannedHeight: 100,
			bestHeight:    200,
			hasResponded:  true,
			wantShouldRun: true,
		},
		{
			name:          "peer responded and wallet scan is caught up",
			scannedHeight: 200,
			bestHeight:    200,
			hasResponded:  true,
			wantShouldRun: false,
		},
		{
			name:          "peer responded, wallet scan is exactly at the lag threshold boundary",
			scannedHeight: 190,
			bestHeight:    200,
			hasResponded:  true,
			wantShouldRun: false,
		},
		{
			name:          "peer responded, wallet scan is one block past the lag threshold",
			scannedHeight: 189,
			bestHeight:    200,
			hasResponded:  true,
			wantShouldRun: true,
		},
		{
			name:          "zero peers responded even when scanned height is huge",
			scannedHeight: 0,
			bestHeight:    0,
			hasResponded:  false,
			wantShouldRun: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldCheckHeightBasedReboot(tc.scannedHeight, tc.bestHeight, tc.hasResponded)
			if got != tc.wantShouldRun {
				t.Errorf("shouldCheckHeightBasedReboot(%v, %v, %v) = %v, want %v",
					tc.scannedHeight, tc.bestHeight, tc.hasResponded, got, tc.wantShouldRun)
			}
		})
	}
}

// TestShouldCheckHeightBasedReboot_UnderflowRegression is the specific regression case from
// finding 6: previously, bestHeight stayed at 0 when every oracle peer was unreachable, and
// uint64(0) - heightLagThreshold underflowed to a value near math.MaxUint64, making the height
// comparison unconditionally true. This test would fail (or, without the hasResponded guard,
// return true) if that bug were reintroduced.
func TestShouldCheckHeightBasedReboot_UnderflowRegression(t *testing.T) {
	if shouldCheckHeightBasedReboot(0, 0, false) {
		t.Fatal("expected no reboot when zero oracle peers responded, regardless of heights")
	}
}
