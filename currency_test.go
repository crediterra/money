package money

import (
	"sort"
	"testing"
)

func TestIsKnownCurrency(t *testing.T) {
	if IsKnownCurrency(CurrencyCode("")) {
		t.Error("Empty currency should not be known")
	}
	for _, c := range currencies {
		if !IsKnownCurrency(c) {
			t.Errorf("%s should be known currency", c)
		}
	}
}

// tsSupportedCurrencies is the exact, alphabetically-sorted output of
// JavaScript's Intl.supportedValuesOf('currency') as captured from Node.js
// on 2026-09-02 (ICU's current ISO 4217 active-currency list, 162 codes).
// ActiveCurrencies must match this set exactly so that the Go and
// TypeScript sides of EventHappening agree on which currencies are active.
var tsSupportedCurrencies = []CurrencyCode{
	"AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN",
	"BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL",
	"BSD", "BTN", "BWP", "BYN", "BZD",
	"CAD", "CDF", "CHF", "CLP", "CNY", "COP", "CRC", "CUC", "CUP", "CVE", "CZK",
	"DJF", "DKK", "DOP", "DZD",
	"EGP", "ERN", "ETB", "EUR",
	"FJD", "FKP",
	"GBP", "GEL", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD",
	"HKD", "HNL", "HRK", "HTG", "HUF",
	"IDR", "ILS", "INR", "IQD", "IRR", "ISK",
	"JMD", "JOD", "JPY",
	"KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT",
	"LAK", "LBP", "LKR", "LRD", "LSL", "LYD",
	"MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN",
	"NAD", "NGN", "NIO", "NOK", "NPR", "NZD",
	"OMR",
	"PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG",
	"QAR",
	"RON", "RSD", "RUB", "RWF",
	"SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLE", "SLL", "SOS", "SRD", "SSP", "STN", "SVC", "SYP", "SZL",
	"THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS",
	"UAH", "UGX", "USD", "UYU", "UZS",
	"VES", "VND", "VUV",
	"WST",
	"XAF", "XCD", "XCG", "XDR", "XOF", "XPF", "XSU",
	"YER",
	"ZAR", "ZMW", "ZWG", "ZWL",
}

func TestActiveCurrencies_MatchesTypeScriptIntlSupportedValues(t *testing.T) {
	if len(tsSupportedCurrencies) != 162 {
		t.Fatalf("test fixture drifted: expected 162 TS currency codes, got %d", len(tsSupportedCurrencies))
	}

	active := ActiveCurrencies()

	got := make([]string, len(active))
	for i, c := range active {
		got[i] = string(c)
	}
	sort.Strings(got)

	want := make([]string, len(tsSupportedCurrencies))
	for i, c := range tsSupportedCurrencies {
		want[i] = string(c)
	}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("ActiveCurrencies() has %d codes, want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ActiveCurrencies() mismatch at index %d: got %q, want %q\nfull got:  %v\nfull want: %v", i, got[i], want[i], got, want)
		}
	}
}

func TestRetiredCurrencies_StillParseButAreNotActive(t *testing.T) {
	retired := []CurrencyCode{"BYR", "MRO", "STD", "VEF"}
	for _, c := range retired {
		if !IsKnownCurrency(c) {
			t.Errorf("%s should still be a known currency for historical data parsing", c)
		}
		if !c.IsRetired() {
			t.Errorf("%s should be reported as retired", c)
		}
		if c.IsActive() {
			t.Errorf("%s should not be active", c)
		}
	}

	for _, active := range ActiveCurrencies() {
		for _, r := range retired {
			if active == r {
				t.Errorf("retired currency %s should not appear in ActiveCurrencies()", r)
			}
		}
	}
}

func TestIsActive_UnknownAndMalformedCodesAreNotActive(t *testing.T) {
	for _, c := range []CurrencyCode{"", "XXX", "GGP[G]", "usd"} {
		if c.IsActive() {
			t.Errorf("%q should not be reported as active", c)
		}
	}
}

func TestCurrencyCode_IsMoney(t *testing.T) {
	if !CurrencyUSD.IsMoney() {
		t.Error("expected USD to be money")
	}
	if CurrencyCode("INVALID").IsMoney() {
		t.Error("expected INVALID not to be money")
	}
}

func TestHasCurrencyPrefix(t *testing.T) {
	if !HasCurrencyPrefix("$100") {
		t.Error("expected $100 to have currency prefix")
	}
	if HasCurrencyPrefix("100 USD") {
		t.Error("expected 100 USD not to have currency prefix")
	}
}

func TestCleanupCurrency(t *testing.T) {
	if got := CleanupCurrency(CurrencyUSD.SignAndCode()); got != CurrencyUSD {
		t.Errorf("expected USD from SignAndCode, got %v", got)
	}
	if got := CleanupCurrency("USD"); got != CurrencyUSD {
		t.Errorf("expected USD from string, got %v", got)
	}
	if got := CleanupCurrency("XYZ"); got != CurrencyCode("XYZ") {
		t.Errorf("expected fallback XYZ, got %v", got)
	}
}

func TestCurrencyCode_Sign_and_SignAndCode(t *testing.T) {
	if got := CurrencyUSD.Sign(); got != "$" {
		t.Errorf("expected $, got %v", got)
	}
	if got := CurrencyCode("XYZ").Sign(); got != "XYZ" {
		t.Errorf("expected XYZ, got %v", got)
	}

	if got := CurrencyUSD.SignAndCode(); got != "$ USD" {
		t.Errorf("expected '$ USD', got %v", got)
	}
	if got := CurrencyCode("XYZ").SignAndCode(); got != "XYZ" {
		t.Errorf("expected XYZ, got %v", got)
	}
}

