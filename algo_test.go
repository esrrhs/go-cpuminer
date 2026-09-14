package main

import (
	"strings"
	"testing"

	"github.com/esrrhs/gohome/crypto"
)

func supportedAlgoNames() []string {
	seen := make(map[string]bool)
	var names []string
	for id := CN_0; id < MAX; id++ {
		name := (&Algorithm{id: id}).supportAlgoName()
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return names
}

func TestSupportedAlgosMatchCrypto(t *testing.T) {
	got := supportedAlgoNames()
	want := crypto.Algo()
	if len(got) != len(want) {
		t.Fatalf("supported algo count = %d, crypto.Algo() = %d\ngot  %v\nwant %v", len(got), len(want), got, want)
	}
	wantSet := make(map[string]bool, len(want))
	for _, name := range want {
		wantSet[name] = true
	}
	for _, name := range got {
		if !wantSet[name] {
			t.Errorf("supported algo %q is missing from crypto.Algo()", name)
		}
	}
}

func TestAlgorithmNamesParse(t *testing.T) {
	for _, item := range algorithm_names {
		t.Run(item.name, func(t *testing.T) {
			a := NewAlgorithm(item.name)
			if a.id != item.id {
				t.Fatalf("id = %d, want %d", a.id, item.id)
			}
		})
		if item.shortName == "" {
			continue
		}
		t.Run(item.shortName, func(t *testing.T) {
			a := NewAlgorithm(item.shortName)
			if a.id != item.id {
				t.Fatalf("id = %d, want %d", a.id, item.id)
			}
		})
		t.Run("upper/"+item.name, func(t *testing.T) {
			a := NewAlgorithm(strings.ToUpper(item.name))
			if a.id != item.id {
				t.Fatalf("id = %d, want %d", a.id, item.id)
			}
		})
	}
}

func TestAlgorithmInvalid(t *testing.T) {
	a := NewAlgorithm("not-a-real-algo")
	if a.id != INVALID {
		t.Fatalf("id = %d, want INVALID", a.id)
	}
	if a.supportAlgoName() != "" {
		t.Fatalf("unsupported algo returned %q", a.supportAlgoName())
	}
	if a.family() != UNKNOWN {
		t.Fatalf("family = %d, want UNKNOWN", a.family())
	}
	if a.name() != "invalid" || a.shortName() != "invalid" {
		t.Fatalf("invalid name = %q / %q", a.name(), a.shortName())
	}
}

func TestAlgorithmFamily(t *testing.T) {
	cases := []struct {
		id     int
		family int
	}{
		{CN_0, CN},
		{CN_1, CN},
		{CN_2, CN},
		{CN_R, CN},
		{CN_FAST, CN},
		{CN_HALF, CN},
		{CN_XAO, CN},
		{CN_RTO, CN},
		{CN_RWZ, CN},
		{CN_ZLS, CN},
		{CN_DOUBLE, CN},
		{CN_CCX, CN},
		{CN_LITE_0, CN_LITE},
		{CN_LITE_1, CN_LITE},
		{CN_HEAVY_0, CN_HEAVY},
		{CN_HEAVY_TUBE, CN_HEAVY},
		{CN_HEAVY_XHV, CN_HEAVY},
		{CN_PICO_0, CN_PICO},
		{CN_PICO_TLO, CN_PICO},
		{RX_0, RANDOM_X},
		{RX_WOW, RANDOM_X},
		{RX_ARQ, RANDOM_X},
		{RX_SFX, RANDOM_X},
		{RX_KEVA, RANDOM_X},
		{AR2_CHUKWA, ARGON2},
		{AR2_CHUKWA_V2, ARGON2},
		{AR2_WRKZ, ARGON2},
		{ASTROBWT_DERO, ASTROBWT},
		{KAWPOW_RVN, KAWPOW},
		{INVALID, UNKNOWN},
	}
	for _, tc := range cases {
		a := &Algorithm{id: tc.id}
		if got := a.family(); got != tc.family {
			t.Errorf("id %d family = %d, want %d", tc.id, got, tc.family)
		}
	}
}

func TestUnsupportedAlgosHaveNoHashName(t *testing.T) {
	unsupported := []int{CN_CCX, RX_0, RX_WOW, RX_ARQ, RX_SFX, RX_KEVA, AR2_CHUKWA, AR2_CHUKWA_V2, AR2_WRKZ, ASTROBWT_DERO, KAWPOW_RVN}
	for _, id := range unsupported {
		a := &Algorithm{id: id}
		if name := a.supportAlgoName(); name != "" {
			t.Errorf("id %d unexpectedly supported as %q", id, name)
		}
	}
}

func TestNewBenchmarkAllAlgos(t *testing.T) {
	for _, name := range crypto.Algo() {
		al := NewAlgorithm(name)
		if al.id == INVALID {
			t.Errorf("%s did not parse", name)
			continue
		}
		if al.supportAlgoName() == "" {
			t.Errorf("%s is not in supportAlgoName()", name)
		}
	}
}

func TestNewBenchmarkInvalid(t *testing.T) {
	if _, err := NewBenchmark("not-a-real-algo"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := NewBenchmark("rx/0"); err == nil {
		t.Fatal("expected unsupported algo error")
	}
}

func TestNewTester(t *testing.T) {
	if _, err := NewTester("not-a-real-algo"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := NewTester("rx/0"); err == nil {
		t.Fatal("expected unsupported algo error")
	}
}
