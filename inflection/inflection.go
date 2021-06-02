package inflection

import (
	"regexp"
	"strings"
)

/* set of voiceless consonants for quick access, fıstıkçı şahap */
var voiceless = map[rune]bool{
	'f': true,
	's': true,
	't': true,
	'k': true,
	'ç': true,
	'ş': true,
	'h': true,
	'p': true,
}

/* set of vowels */
var vowel = map[rune]bool{
	'a': true,
	'e': true,
	'ı': true,
	'i': true,
	'o': true,
	'ö': true,
	'u': true,
	'ü': true,
	'A': true,
	'I': true,
}

type quality struct {
	front, round, high bool
}

/* map of vowels to vowel (harmony) qualities */
var vowel_to_quality = map[rune]quality{
	'a': quality{false, false, false},
	'e': quality{true, false, false},
	'ı': quality{false, false, true},
	'i': quality{true, false, true},
	'o': quality{false, true, false},
	'ö': quality{true, true, false},
	'u': quality{false, true, true},
	'ü': quality{true, true, true},
	'A': quality{false, false, false}, /* can be both front or back but always low */
	'I': quality{false, false, true},  /* can be both front or back, round or flat but always high */
}

/* map front, round, high qualities, in order, to vowels */
var quality_to_vowel = map[quality]rune{
	quality{false, false, false}: 'a',
	quality{true, false, false}:  'e',
	quality{false, false, true}:  'ı',
	quality{true, false, true}:   'i',
	quality{false, true, false}:  'o',
	quality{true, true, false}:   'ö',
	quality{false, true, true}:   'u',
	quality{true, true, true}:    'ü',
}

/*
A Root can contain only exact characters except the final character which can be B,C,D,K,(n).
The final letter is realized when a suffix is appended or it is converted to a word.
*/
type Root []rune

/*
A Stem can contain only exact characters except the final character which can be B,C,D,K,(n).
The final letter is realized when a suffix is appended or it is converted to a word.
*/
type Stem []rune

/* A Word only contains exact (fully resolved) characters. */
type Word []rune

/*
A Suffix has body that is a list of runes and an optional, single-character head and tail
A value of 0 for head and tail means no character. Only 'n' is a valid tail.
*/
type Suffix struct {
	Head, Tail rune
	Body       []rune
}

/*
takes in a vowel and front/round harmony it should conform to
If vowel is A/I, returns the new quality and adjusted form of vowel
If vowel is exact, returns the same vowel and its quality
*/
func resolve_vowel(vowel rune, front, round bool) (q quality, v rune) {
	switch vowel {
	case 'A':
		v = quality_to_vowel[quality{front, false, false}]
	case 'I':
		v = quality_to_vowel[quality{front, round, true}]
	default:
		v = vowel
	}
	q = vowel_to_quality[v]

	return q, v
}

/*
takes in B/C/D/K/N or an exact consonant and the previous and next consonant
returns correct consonant mutation form based on previous and next
*/
func resolve_cons(prev, c, next rune) rune {
	if vowel[next] && prev != 0 && !voiceless[prev] {
		switch c {
		case 'B':
			c = 'b'
		case 'C':
			c = 'c'
		case 'D':
			c = 'd'
		case 'K':
			c = 'g'
		}
	} else {
		switch c {
		case 'B':
			c = 'p'
		case 'C':
			c = 'ç'
		case 'D':
			c = 't'
		case 'K':
			c = 'k'
		}
	}
	if vowel[prev] && c == 'g' {
		c = 'ğ'
	}
	if c == 'N' {
		if next == 0 {
			c = 0
		} else {
			c = 'n'
		}
	}
	return c
}

/*
Combines the suffix with the stem but does not resolve final N/B/C/D/K after appending
Only resolves the consonant and vowel harmonies of the suffix and the final consonant
of the original stem if it exists. Does not modify inputted stem
*/
func (stem Stem) Append(suffix Suffix) Stem {
	s := Stem(make([]rune, len(stem)))
	copy(s, stem)

	/* add optional suffix head if it is the opposite type (vowel/consonant) of the stem's final word */
	if suffix.Head != 0 && vowel[s[len(s)-1]] != vowel[suffix.Head] {
		s = append(s, suffix.Head)
	}

	/* drop stem-final vowel if suffix begins with a vowel (-Iyor) */
	if vowel[s[len(s)-1]] && len(suffix.Body) != 0 && vowel[suffix.Body[0]] {
		s = s[:len(s)-1]
	}

	s = append(s, suffix.Body...)
	if suffix.Tail != 0 {
		s = append(s, 'N') /* (n) is the only valid suffix */
	}

	/* get quality of latest exact vowel in stem */
	front, round := false, false // quality of latest vowel
	for i := len(stem) - 1; i >= 0; i-- {
		if vowel[s[i]] && s[i] != 'A' && s[i] != 'I' {
			q := vowel_to_quality[s[i]]
			front, round = q.front, q.round
			break
		}
	}

	for i := len(stem) - 1; i < len(s)-1; i++ {
		if vowel[s[i]] {
			var q quality
			q, s[i] = resolve_vowel(s[i], front, round)
			front, round = q.front, q.round
		} else {
			var prev rune
			if i == 0 {
				prev = 0
			} else {
				prev = s[i-1]
			}
			s[i] = resolve_cons(prev, s[i], s[i+1])
		}
	}

	if vowel[s[len(s)-1]] {
		_, s[len(s)-1] = resolve_vowel(s[len(s)-1], front, round)
	}

	return s
}

/* fully resolves the stem (resolves final consonant) and returns as Word */
func (stem Stem) Word() Word {
	w := Word(make([]rune, len(stem)))
	copy(w, stem)
	if !vowel[w[len(w)-1]] {
		/* value of prev is irrelevant; next == 0 implies a voiceless */
		w[len(w)-1] = resolve_cons(0, w[len(w)-1], 0)
	}
	return w
}

func (suffix Suffix) String() string {
	head, tail := "", ""
	if suffix.Head != 0 {
		head = "(" + string(suffix.Head) + ")"
	}
	if suffix.Tail != 0 {
		tail = "(" + string(suffix.Tail) + ")"
	}
	return head + string(suffix.Body) + tail
}

func (root Root) String() string {
	return string(root)
}

func (stem Stem) String() string {
	return string(stem)
}

func (word Word) String() string {
	return string(word)
}

/*
The root of a word is a list of exact characters. The final character can be one of
B/C/D/K or (n). n must be parenthesized if it is used as an optional final character
*/
func ParseRoot(s string) (r Root, ok bool) {
	re := regexp.MustCompile(`^\s*([a-zçğıöşü]*)(?:([a-zçğıöşüBCDK])|(?:\((n)\)))\s*$`)
	if matches := re.FindStringSubmatch(s); len(matches) == 4 {
		if matches[3] != "" {
			return Root(matches[1] + "N"), true
		}
		return Root(matches[1] + matches[2]), true
	}
	return Root(""), false
}

/*
The suffix is a sequence of exact characters or A/I/B/C/D/K consisting of a body and
an optional head and tail character marked by parenthesis. The tail can only be (n).
*/
func ParseSuffix(s string) (suf Suffix, ok bool) {
	re := regexp.MustCompile(
		`^\s*(?:\(([a-zçğıöşüBCDKAI])\))?([a-zçğıöşüBCDKAI]+)(?:\((n)\))?\s*$`)
	if matches := re.FindStringSubmatch(s); len(matches) == 4 {
		var h, t rune = 0, 0
		if matches[1] != "" {
			h = ([]rune(matches[1]))[0]
		}
		if matches[3] != "" {
			t = ([]rune(matches[3]))[0]
		}
		return Suffix{Head: h, Tail: t, Body: []rune(matches[2])}, true
	}
	return Suffix{Head: 0, Tail: 0, Body: nil}, false

}

/*
Parses a root followed by a sequence of suffixes. The input must be of the form
ROOT SUFFIX SUFFIX ...
That is, the root followed by a sequence of suffixes with whitespace as delimiter
returns the Root and a slice of Suffixes and true on success
returns the nil values and false on error
*/
func ParseRootSuffixes(s string) (root Root, sufs []Suffix, ok bool) {
	words := strings.Fields(s)
	if len(words) == 0 {
		return Root(nil), []Suffix(nil), false
	}
	root, ok = ParseRoot(words[0])
	if !ok {
		return Root(nil), []Suffix(nil), false
	}
	for i := 1; i < len(words); i++ {
		suf, ok := ParseSuffix(words[i])
		if !ok {
			return Root(nil), []Suffix(nil), false
		}
		sufs = append(sufs, suf)
	}
	return root, sufs, true

}
