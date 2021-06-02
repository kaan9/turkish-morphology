package inflection

import (
	"reflect"
	"testing"
)

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
		if !ok || !reflect.DeepEqual(valid_out[i], r) {
			t.Errorf("ParseRoot(%s) = (%#v, %v), expected (%#v, %v)", s, r, ok, valid_out[i], true)
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
			t.Errorf("ParseRoot(%s) = (%#v, %v), expected (%#v, %v)", s, r, ok, nil, false)
		}
	}
}

func TestParseSuffix(t *testing.T) {
	valid := []string{
		"abc", "AbC", "(I)n", "(n)In", "(I)m", "   (K)AAAAIIII(n)\t\r\n   ",
		"(y)AcAK", "sIz", "(s)I(n)", "lIK", "CI", "KI", "lI", "DAş",
		"(A)KDDCCBB(n)", "(I)AI(n)",
	}
	valid_out := []Suffix{
		Suffix{Head: 0, Tail: 0, Body: []rune("abc")},
		Suffix{Head: 0, Tail: 0, Body: []rune("AbC")},
		Suffix{Head: 'I', Tail: 0, Body: []rune("n")},
		Suffix{Head: 'n', Tail: 0, Body: []rune("In")},
		Suffix{Head: 'I', Tail: 0, Body: []rune("m")},
		Suffix{Head: 'K', Tail: 'n', Body: []rune("AAAAIIII")},
		Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},
		Suffix{Head: 0, Tail: 0, Body: []rune("sIz")},
		Suffix{Head: 's', Tail: 'n', Body: []rune("I")},
		Suffix{Head: 0, Tail: 0, Body: []rune("lIK")},
		Suffix{Head: 0, Tail: 0, Body: []rune("CI")},
		Suffix{Head: 0, Tail: 0, Body: []rune("KI")},
		Suffix{Head: 0, Tail: 0, Body: []rune("lI")},
		Suffix{Head: 0, Tail: 0, Body: []rune("DAş")},
		Suffix{Head: 'A', Tail: 'n', Body: []rune("KDDCCBB")},
		Suffix{Head: 'I', Tail: 'n', Body: []rune("AI")},
	}

	for i, s := range valid {
		suf, ok := ParseSuffix(s)
		if !ok || !reflect.DeepEqual(valid_out[i], suf) {
			t.Errorf("ParseSuffix(%s) = (%#v, %v), expected (%#v, %v)", s, suf, ok, valid_out[i], true)
		}
	}

	invalid := []string{
		"(I)", "(I)(n)", "buN", "bu(s)", "(KC)bcd", "(A)bee(I)", "()a()",
		"(N)A",	"(IA)I(n)", "n(A)", "(nn)nn(n)",
	}
	for _, s := range invalid {
		suf, ok := ParseSuffix(s)
		if ok {
			t.Errorf("ParseSuffix(%s) = (%#v, %v), expected (%#v, %v)", s, suf, ok, nil, false)
		}
	}

}


func TestParseRootSuffixes(t *testing.T) {
	valid := []string{
		"tanı (I)ş DIr (I)l (y)AmA  \t\r\n   (y)Abil (y)AcAK lAr DAn (y)mIş çA (s)I(n)  (y)A",
		"bu(n) lAr  (n)In  ki   lAr    DAn mI  (y)mIş",
	}

	valid_roots := []Root{
		Root("tanı"), Root("buN"),
	}

	valid_sufs := [][]Suffix{
		[]Suffix{
			Suffix{Head: 'I',  Tail: 0, Body: []rune("ş")},
			Suffix{Head: 0, Tail: 0, Body: []rune("DIr")},
			Suffix{Head: 'I', Tail: 0, Body: []rune("l")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("AmA")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("Abil")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},
			Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},
			Suffix{Head: 0, Tail: 0, Body: []rune("DAn")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("mIş")},
			Suffix{Head: 0, Tail: 0, Body: []rune("çA")},
			Suffix{Head: 's',  Tail: 'n', Body: []rune("I")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("A")},
		},
		[]Suffix{
			Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},
			Suffix{Head: 'n', Tail: 0, Body: []rune("In")},
			Suffix{Head: 0, Tail: 0, Body: []rune("ki")},
			Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},
			Suffix{Head: 0, Tail: 0, Body: []rune("DAn")},
			Suffix{Head: 0, Tail: 0, Body: []rune("mI")},
			Suffix{Head: 'y', Tail: 0, Body: []rune("mIş")},
		},

	}

	for i, s := range valid {
		root, sufs, ok := ParseRootSuffixes(s)
		if !ok || !reflect.DeepEqual(valid_roots[i], root) || !reflect.DeepEqual(valid_sufs[i], sufs) {
			t.Errorf("ParseRootSuffix(%s) = (%#v, %#v, %v), expected (%#v, %#v, %v)",
			s, root, sufs, ok, valid_roots[i], valid_sufs[i], true)
		}
	}


	invalid := []string{
		"",
	}
	for _, s := range invalid {
		root, sufs, ok := ParseRootSuffixes(s)
		if ok {
			t.Errorf("ParseRootSuffix(%s) = (%#v, %#v, %v), expected (%#v, %#v, %v)",
			s, root, sufs, ok, nil, nil, false)
		}
	}

}


func TestAppend(t *testing.T) {
	valid := []string{
		"tanı (I)ş DIr (I)l (y)AmA  \t\r\n   (y)Abil (y)AcAK lAr DAn (y)mIş çA (s)I(n)  (y)A",
		"bu(n) lAr  (n)In  ki   lAr    DAn mI  (y)mIş",
	}

	valid_stems := [][]Stem{
		[]Stem{
			Stem("tanı"),
			Stem("tanış"),
			Stem("tanıştır"),
			Stem("tanıştırıl"),
			Stem("tanıştırılama"),
			Stem("tanıştırılamayabil"),
			Stem("tanıştırılamayabileceK"),
			Stem("tanıştırılamayabilecekler"),
			Stem("tanıştırılamayabileceklerden"),
			Stem("tanıştırılamayabileceklerdenmiş"),
			Stem("tanıştırılamayabileceklerdenmişçe"),
			Stem("tanıştırılamayabileceklerdenmişçesiN"),
		},
		[]Stem{
			Stem("buN"),
			Stem("bunlar"),
			Stem("bunların"),
			Stem("bunlarınki"),
			Stem("bunlarınkiler"),
			Stem("bunlarınkilerden"),
			Stem("bunlarınkilerdenmi"),
			Stem("bunlarınkilerdenmiymiş"),
		},
	}

	valid_words := []Word{
		Word("tanıştırılamayabileceklerdenmişçesine"),
		Word("bunlarınkilerdenmiymiş"),
	}

	for i, s := range valid {
		root, sufs, _ := ParseRootSuffixes(s)
		stem := Stem(root)
		for j, suf := range sufs {
			if !reflect.DeepEqual(stem, valid_stems[i][j]) {
				t.Errorf("Append(%v) = %v, expected %v", suf, stem, valid_stems[i][j])
			}
			stem = stem.Append(suf)
		}
		word := Word(stem)
		if !reflect.DeepEqual(word, valid_words[i]) {
			t.Errorf("Word(%v) = %v, expected %v", stem, word, valid_words[i])
		}
	}
}
