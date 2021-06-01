# Turkish Morphology
Algorithms implementing features of Turkish morphology

## Phonology

### Consonants
|	|	| Labial | Dental/Alveolar | Postalveolar | Palatal | Velar | Glottal |
|-------|-------|:------:|:---------------:|:------------:|:-------:|:-----:|:-------:|
| Nasal | | m | n | | | | |
| Plosive | voiceless | p | t | ç | (k) | k | |
|  | voiced | b | d | c | (g) | g | |
| Fricative | voiceless | f | s | ş | | | h |
|  | voiced | v | z | j | | | |
| Approximant | | (v) | (l) | l | y | (ğ) | |
| Flap | | | | r | | | |

### Vowels

|	| Front	|	| Back	|	|
|-------|:-----:|:-----:|:-----:|:-----:|
|	| flat	| round	| flat	| round	|
|Close (High)	| i	| ü	| ı	| u	|
|Open (Low)	| e	| ö	| a	| o	|


## Package `inflection`
Inflectional morphology of turkish words.
Functions to perform agglutination (root + suffix + suffix ...) while respecting phonotactics (vowel harmony, consonant mutation) and their exceptions.

The following representations are in terms of **exact** characters `a,b,c,ç,...` (the lowercase Turkish characters) which only match themselves and the **varying** characters `A,I,B,C,D,K` whose matching rules are described below.

A word consists of a `Root` and `Suffix`es attached to it. The `Root` and every form created by suffixation is a `Stem`. The `Stem` is an intermediate form that is finalized when it is converted to `Word`.

* A `Root` or `Stem` consists of only exact characters except for the final character which may be one of `B,C,D,K` and `(n)`.  For example `yap`, `giD`, or `bu(n)`.

* A `Suffix` has a body of at least 1 character and a parenthesized optional head and tail character. The head and body can contain all exact and varying characters. The tail can only be `(n)`.  A consonant head is realized if and only if the stem's final character is a vowel and vice-versa. The tail is realized if and only if a suffix is added.

* A `Word` only contains resolved (exact) characters.

---

The only character that appears as an optional tail is 'n', which appears in the roots bu, şu, o:
```
 bu(n) + (y)I -> bunu, o(n) + (y)A -> ona, etc.
```
and the 3rd person possessive suffix `-(s)I(n-)` as well as its plural form
```
lar + -(s)I(n) -> larI(n)
```
Therefore, the optional tail, when part of the stem, is internally represented in the `Root` or `Stem` as:
```
N = 'n' or nothing
```
`N` is realized as 'n' if and only if a suffix is added after it. (The actual rules are concerning the optional `n` are more complex and not implemented at this level.)

---

The addition of a suffix can never cause consecutive vowels and generally does not cause consecutive consonants
in the same syllable. An exception is `-t`: 
```
yaptır + -t -> yaptırt
```

Therefore, if the stem ends in a vowel, the body begins with a vowel, and there is no optional consonant head, the stem's final vowel is dropped. The only instance
of this is the -Iyor suffix:
```
başla + Iyor -> başlıyor
```

---

For the stem and the suffix head and body, the **varying** characters used for encoding are:

* `B = b/p, C = c/ç, D = d/t K = k/g/ğ`
	* voiced if preceded by a voiced consonant and followed by a vowel
	* voiced if in between two vowels
	* voiceless otherwise
	* g instead of ğ if preceded by a consonant

* `A = a/e, I = ı/i/u/ü`
	* `A` and `I` represent low and high vowels and are resolved using vowel front-back and rounding harmony
	
Consonant voicing generally happens in multisyllabic stems but there are exceptions (git, gidiyorum) which are
encoded in the root: `git` is encoded as `giD`.

---

The package defines:

* The types `Root`, `Suffix`, `Stem`, `Word` and their `Stringer` interface implementations
* The method `Word()` on `Stem` that fully resolves the `varying` characters
* The method `Append(Suffix)` on `Stem` that produces a new stem with the suffix attached
* The functions `ParseRoot`, `ParseSuffix`, and `ParseRootSuffixes` take in a string and parse a `Root`, a `Suffix`, and a `Root` followed by a variable number of `Suffix`es

#### Examples
`yap + Iyor + (y)sA + (I)m` which should produce `yapıyorsam`

`bu(n) lAr (n)In ki lAr DAn` which should produce  `bunlarinkilerden`

Example implementation:
```
root, suffixes, _ := ParseRootSuffixes("bu(n) lAr (n)In ki lAr DAn")
stem := Stem(root)
for _, suf := range sufs {
	stem = stem.Append(suf)
}
fmt.Printf("%v\n", stem)	// prints bunlarinkilerden
```

