package inflection

import "testing"

func TestParseRoot(t *testing.T) {
	valid := []string{
		"abc", "yap", "git", "k", "kat", "et", "ççç", "daB", "caC", "saD",
		"giD", "gök", "göK", "   uouou   ", "\t\r\nvvvrrruuuD\t\r\n", "cac",
		"bu(n)", "o(n)", "  ka(n)  ",
	}
	valid_out := []Root{
		Root("abc"), Root("yap"), Root("git"),
		Root("k"), Root("kat"), Root("et"),
		Root("ççç"), Root("daB"), Root("caC"),
		Root("saD"), Root("giD"), Root("gök"),
		Root("göK"), Root("uouou"), Root("vvvrrruuuD"),
		Root("cac"), Root("buN"), Root("oN"), Root("kaN"),
	}
	for i, s := range valid {
		r, ok := ParseRoot(s)
		if !ok || valid_out[i].String() != r.String() {
			t.Errorf("ParseRoot(%s) = (%v, %v), expected (%v, %v)", s, r, ok, valid_out[i], true)
		}
	}

	invalid := []string{
		"aBc", "aCb", "AcB", "öööI", "KaaaaK", "abc(d)", "ya(K)",
		"yaN", "KiD", "yaK(n)", "yaKn", "aCaK", "AcAK", "(I)a(n)",
		"(K)oK", "(K)o(n)", "(C)ac",
	}
	for _, s := range invalid {
		r, ok := ParseRoot(s)
		if ok {
			t.Errorf("ParseRoot(%s) = (%v, %v), expected (%v, %v)", s, r, ok, nil, false)
		}
	}
}
