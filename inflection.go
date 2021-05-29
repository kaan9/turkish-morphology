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
}

/*
The Stem representation contains the unrealized forms B,C,D,K,N,A,E. These are converted to actual
letter when the Stem is converted to a string
*/
type Stem []rune
type Root Stem

/* head and tail are optional (single character for both) and a value of 0 implies no head/tail */
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
this function appends the suffix to the stem but does not resolve the vowel and consonant mutations
it handles the optional head and tail and the vowel drop
*/
func (stem Stem) append(suffix Suffix) Stem {
//	if vowel[stem[len(stem)-1]] != vowel[suffix.head[len(suffix.head)-1]] {
//		stem += suffix.head
//	}
//	stem += suffix.body
	return stem
}

/*
stringer interface, this function returns the fully resolved stem
the start of the stem must not contain A/E/B/C/D/K/N
*/
func (stem Stem) String() string {
	s := ""

	s += string(stem[0])

	return s

}


func main() {

}
