package main

/*
Inflectional morphology of turkish words. Specifically, functions to perform agglutination
(root + suffix + suffix ...) while respecting phonotactics (vowel harmony, consonant mutation) and exceptions.
*/

//import "fmt"

/*
Each suffix has a body that is always included and an optional head and tail character that are included for
phonotactics. A consonant head is realized iff the stem's final character is a vowel and a vowel head is realized
iff the stem's final character is a consonant.
The only character that appears as an optional tail is 'n', which appears in the roots bu, şu, o
(bu + i -> bunu, o + a -> ona, etc.) and the 3rd person possessive suffix -(s)I(n-) (as well as its plural form
-lar + -(s)I(n-) -> -larI(n-)). Therefore, the optional tail, when part of the stem, is represented by:
 - N = 'n' or nothing
N is realized as 'n' if and only if a suffix is added after it. (The actual rules are more complex and not fully
implemented, this is a simplifying assumption.)

The addition of a suffix can never cause consecutive vowels (and generally does not cause consecutive consonants
in the same syllable, an exception is -t: yaptır + -t -> yaptırt) so if the stem ends in a vowel, there is no
optional consonant head, and the body begins with a vowel, the stem's final vowel is dropped. The only instance
of this is the -iyor suffix: for -iyor: başla + iyor -> başlıyor

For the stem and the suffix head and body, the characters used for encoding are:
- lowercase characters are not subject to vowel harmony or consonant mutation and match themselves exactly.
- uppercase characters are subject to both and the encodings are:
  B = b/p, C = c,ç, D = d,t K = k/g/ğ	A = a/e (low), I = ı/i/u/ü (high)
  K becomes a ğ when a vowel follows it except when preceded by a consonant in which case it becomes g
  consonant voicing generally happens in multisyllabic stems but there are exceptions (git, gidiyorum) which are
  encoded in the root; "git" would be encoded as "giD"
The tail can only contain N which can also be the stem's final character. N is realized as 'n' if followed by
a vowel and is removed otherwise.

as an example, the suffix -iyor would be written as "Iyor" and the root 'bu' is represented as the stem 'buN'

Summary:
- Suffix heads are optional letters and are realized iff the stem's final phoneme is a vowel.
- Suffix/stem tail is an optional 'N' which is realized iff a suffix follows it.
- Suffix bodies are mandatory lists of phonemes.
	- If the first phoneme is a vowel and the stem ends in a vowel, the stem's vowel is droped.
	- Vowel harmony: A = a/e (low), I = ı/i/u/ü (high) based on front-back and roundness of previous vowel.
	- Consonant mutation:  B = b/p, C = c,ç, D = d,t K = k/g/ğ
		- voiced if preceded by a voiced consonant
		- voiceless if preceded by a voiceless consonant
		- voiced if in between two vowels
		- g instead of ğ iff preceded by a consonant

An example input can be of the form
- yap + Iyor + (y)sA + (I)m
which should produce the output
- yapıyorsam

*/

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
	quality{false, false, false}: 'A',
	quality{false, false, true}:  'I',
}

/*
The Stem representation can contain the unrealized forms B,C,D,K,N only at the end, all other
letters must be fully realized. The final letter is realized when a suffix is appended or the
word is converted to a string
*/
type Stem []rune
type Root Stem

/*
head and tail are optional (single character for both) and a value of 0 implies no head/tail
body must be non-empty
*/
type Suffix struct {
	head, tail rune
	body       []rune
}

var suffixes = map[string]Suffix{
	/* tense/aspect (does not include -makta, which can be encoded as as -mak + -ta */
	"known past": Suffix{head: 0, tail: 0, body: []rune("DA")},
	"infer past": Suffix{head: 0, tail: 0, body: []rune("mIş")},
	"aorist a":   Suffix{head: 'A', tail: 0, body: []rune("r")},
	"aorist i":   Suffix{head: 'I', tail: 0, body: []rune("r")},
	"aorist neg": Suffix{head: 0, tail: 0, body: []rune("mAz")},
	"pres cont":  Suffix{head: 0, tail: 0, body: []rune("Iyor")},
	"fut":        Suffix{head: 'y', tail: 0, body: []rune("AcAK")},
	/* verb negation */
	"neg": Suffix{head: 0, tail: 0, body: []rune("mA")},
	/* infinitive */
	"inf": Suffix{head: 0, tail: 0, body: []rune("mAk")},
}

/*
takes in a vowel and front/round harmony it should conform to
If vowel is A/I, returns the new quality and adjusted form of vowel
If vowel is exact, returns the same vowel and its quality
*/
func resolve_vowel(vowel rune, front, round bool) (q quality, v rune) {
	switch v {
	case 'A':
		v = quality_to_vowel[quality{front, round, false}]
	case 'E':
		v = quality_to_vowel[quality{front, round, true}]
	default:
		v = vowel
	}
	q = vowel_to_quality[v]

	return q, v
}

/*
V: vowel, X: any consonant, H: voiceless, S: voices, 0: empty (start/end of word)
V - S - V // k->ğ
S - S - V // derdi    t softens to d when followed by vowel     k->g
X - H - 0
0 - H - X // choose voiceless as default if first letter of word
H - H - V // kastırmak     üst üstün     test -> testin
H - H - H // shouldn't happen?   üsttürmek if üst was a verb, anyways if should be hard regardless
S - H - X //   dertli   t is hard,  dertsiz still hard
V - H - X // katlı   katsız  gerekli     gereksiz     stays hard
H - H - S // üstlü     sarkmak

might be more useful to analyze this as consonants are hard by default and soften in specific cases
*/

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
func (stem Stem) add(suffix Suffix) Stem {
	s := Stem(make([]rune, len(stem)))
	copy(s, stem)

	/* add optional suffix head if it is the opposite type (vowel/consonant) of the stem's final word */
	if suffix.head != 0 && vowel[s[len(s)-1]] != vowel[suffix.head] {
		s = append(s, suffix.head)
	}

	/* drop stem-final vowel if suffix begins with a vowel (-Iyor) */
	if vowel[s[len(s)-1]] && vowel[suffix.body[0]] {
		s = s[:len(s)-1]
	}

	s = append(s, suffix.body...)
	if suffix.tail != 0 {
		s = append(s, suffix.tail)
	}

	/* get quality of latest exact vowel in stem */
	front, round := false, false // quality of latest vowel
	for i := len(s) - 1; i >= 0; i-- {
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

	return s
}

/*
stringer interface, this function returns the fully resolved stem
the start of the stem must not contain A/E/B/C/D/K/N
*/
func (stem Stem) String() string {
	s := ""

	return s
}

func main() {

}
