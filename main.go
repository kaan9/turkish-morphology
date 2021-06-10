package main

import (
	"bufio"
	"fmt"
	inf "github.com/kaan9/turkish-morphology/inflection"
	"os"
)

// idea: construct an FSA where the keys are the types of suffixes (e.g. PL, NOM, ACC), etc. The nodes of the FSA
// are the types of words (noun, verb, adj, adverb, etc.) and intermediate forms. With each transition, the
// suffixation (as described by in inflection.Append) is performed on the stem. Then stemming can be achieved
// by reversing the FSA and attempting to 'consume' the ending (the opposite of how it was suffixed). FSA should
// be handled using

/* from https://www.dnathan.com/language/turkish/tsd/index.htm */
var suffixes = map[string]inf.Suffix{
	/* plural*/
	"PL": inf.Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},

	/* possessive (iyelik) */
	"POS.1sg": suffix("(I)m"),
	"POS.1pl": suffix("(I)mIz"),
	"POS.2sg": suffix("(I)n"),
	"POS.2pl": suffix("(I)nIz"),
	"POS.3sg": suffix("(s)I(n)"),
	"POS.3pl": suffix("lArI(n)"),

	/* the familial (kinship) -gil and -ler suffix: e.g. teyzemler, karıncayiyengiller */
	"KIN":    suffix("gil"), /* no consonant/vowel harmony */
	"KIN.PL": suffix("lAr"), /* not the same as PL */

	/* case (all except def. accusative can be come before predicative personal suffixes?) */
	"ABSL":    suffix(""),     /* Absolute (yalın) case */
	"ACC":     suffix("(y)I"), /* Definite accusative */
	"DAT":     suffix("(y)A"), /* dative-directional/lative */
	"GEN":     suffix("(n)In"),
	"LOC":     suffix("DA"),
	"ABL":     suffix("DAn"),
	"INS":     suffix("(y)lA"), /* also postposition 'ile' */

	/* Personal Suffixes (kişi ekleri) */
	/* Predicative Personal Suffix - type I (Copular and after -mIş -AcAK -(A/I)r -Iyor ... other forms) */
	"PRED.1sg": suffix("(y)Im"),
	"PRED.1pl": suffix("(y)Iz"),
	"PRED.2sg": suffix("sIn"),
	"PRED.2pl": suffix("sInIz"),
	"PRED.3sg": suffix(""),
	"PRED.3pl": suffix("lAr"),
	/* Verbal Personal Suffix -- type II (after -DI and -sA) */
	"VB.1sg": suffix("m"),
	"VB.1pl": suffix("k"),
	"VB.2sg": suffix("n"),
	"VB.2pl": suffix("nIz"),
	"VB.3sg": suffix(""),
	"VB.3pl": suffix("lAr"),
	/* Optative Personal Suffix -- type III (optative mood) */
	"OPT.1sg": suffix("(y)AyIm"),
	"OPT.1pl": suffix("(y)AlIm"),
	"OPT.2sg": suffix("(y)AsIn"),
	"OPT.2pl": suffix("(y)AsInIz"),
	"OPT.3sg": suffix("(y)A"),
	"OPT.3pl": suffix("(y)AlAr"),
	/* Imperative Personal Suffix -- type IV (imperative mood) */
	"IMP.2sg":  suffix(""),
	"IMP.2pl":  suffix("(y)In"),
	"IMP.2pl2": suffix("(y)InIz"), /* more formal */
	"IMP.3sg":  suffix("sIn"),
	"IMP.3pl":  suffix("sInlAr"),

	/* tense/aspect/mood */
	"TAM.PPFV.KNWN": suffix("DI"),   /* past perfective */
	"TAM.PPFV.INFR": suffix("mIş"),  /* inferred past perfective */
	"TAM.AOR.A":     suffix("(A)r"), /* aorist low vowel */
	"TAM.AOR.I":     suffix("(I)r"), /* aorist high vowel */
	"TAM.AOR.NEG":   suffix("z"),    /* aorist negative/impotential */
	/* AOR.NEG always comes after -mA or -(y)AmA (NEG/INAB); is irregular with 1sg, 1pl:
	yapmam, yapamam, yapmayız, yapamayız (rather than yapmazım, yapamazım, yapmazız, yapamazız),
	but the forms are correct with the interrogative: yapamaz mıyım, yapamaz mıyız, etc. */
	"TAM.PRS.IPFV": suffix("Iyor"),    /* present imperfective */
	"TAM.PRS.PROG": suffix("mAktA"),   /* pres. progressive: -mAK + -DA */
	"TAM.FUT":      suffix("(y)AcAK"), /* future */
	"TAM.COND": suffix("sA"),   /* conditional mood */
	"TAM.NEC":  suffix("mAlI"), /* necessitative mood: -mA + -lI */

	/* copula (comes after the same suffixes as the predicative personal suffixes (type I)) */
	"COP": suffix("DIr"), /* alethic modality */
	/* negative copula indicated with 'değil' which takes copula suffixes */
	"COP.PST":      suffix("(y)DI"),  /* alethic past tense */
	"COP.PST.INFR": suffix("(y)mIş"), /* alethic inferred tense */
	"COP.COND":     suffix("(y)sA"),  /* conditional mood copula */

	/* verbal noun */
	"INF": suffix("mAK"),   /* infinitive */
	"GER": suffix("mA"),    /* gerund */
	"WAY": suffix("(y)Iş"), /* 'way/act of doing' verb */

	/* interrogative particle */
	"INT": suffix("mI"), /* written separate by convention*/

	/* grammatical voice */
	"REFL":   suffix("(I)n"), /* reflexive voice (or pass.) */
	"RECP":   suffix("(I)ş"), /* reciprocal voice */
	"PASS":   suffix("(I)l"), /* passive voice */
	"CAUS.1": suffix("t"),    /* causative type I */
	"CAUS.2": suffix("DIr"),  /* causative type II */
	/* REFL is used as PASS in some cases (when the verb ends in 'lV' with V a vowel),
	CAUS.1 is used afer -l, -r, or a vowel in stems with multiple syllables
	CAUS.2 is used elsewhere but there are many irregular forms which should be anaylzed as their own roots
	The resultant meaning is constructive and not always apparent from the suffixes
	These suffixes can be chained: REFL+PASS, CAUS+CAUS=FAC (factitive), RECP+CAUS=REP (repetitive),
	CAUS.1+CAUS.2+CAUS.1 (causatives can be chained arbitrarily, alternatingly), REFL+PASS+CAUS, etc. */

	/* verb negation and potential, these precede tense/aspect/mood and must precede aorist negative */
	"NEG":  suffix("mA"),
	"INAB": suffix("(y)AmA"), /* impotential */

	/* Participles (separated as personal (always takes suffix of possession) versus impersonal) */
	"PTCP.IMPRS.AOR.A":    suffix("(A)r"),    /* aorist low vowel */
	"PTCP.IMPRS.AOR.I":    suffix("(I)r"),    /* aorist high vowel */
	"PTCP.IMPRS.AOR.NEG":  suffix("z"),     /* aorist negative/impotential (used with -mA/-(y)AmA) */
	"PTCP.IMPRS.IPFV":     suffix("(y)An"),   /* imperfective */
	"PTCP.IMPRS.FUT":      suffix("(y)AcAK"), /* impersonal future */
	"PTCP.PERS.FUT":       suffix("(y)AcAK"), /* personal future */
	"PTCP.IMPRS.PPFV":      suffix("mIş"),     /* impersonal inferred past perfective */
	"PTCP.PERS.PPFV":       suffix("DIK"),     /* personal (known) past perfective */

	/* Converbs  --  verb to adverb suffixes */
	/* converb occurs simultaneously with verb */
	"CVB.1": suffix("(y)A"),
	/* converb while or before main verb (konuşarak bekledik, düşünerek buldum),'olarak' means 'as' */
	"CVB.2": suffix("(y)ArAK"),
	/* NOT a GER+ABL (maybe comes from it), action not occurring or action following main verb */
	"CVB.3": suffix("mAdAn"),
	/* simultaneous, only comes after tenses: not yapken, yaparken/yapacakken/yapmışken etc. */
	"CVB.4": suffix("(y)ken"),
	"CVB.5": suffix("(y)Ip"), /* converb completed before verb */


	/* Verbs used as suffixes -- typically by combining with Converb -(y)A- */
	"VSX.ABIL":	suffix("(y)Abil"),	/* ability, opposite of INAB */
	"VSX.REPT":	suffix("(y)Agel"),	/* repetitive aspect */
	"VSX.SWFT":	suffix("(y)Iver"),	/* "swiftness" aspect */
	"VSX.CONT":	suffix("(y)Adur"),	/* continuous aspect */
	"VSX.NEXP":	suffix("(y)Akal"),	/* continuous aspect, unexpected (e.g. bakakalmak) */
	"VSX.NEAR":	suffix("(y)Ayaz"),	/* "almost happened" */


	/* The ki suffix -- acts as relative pronoun to create relative clause? */
	"REL": suffix("ki"), /* personal perfective */

	/* head marker -- attached to modified noun when a noun modifies another noun (same as POS.3sg) */
	"HD": suffix("(s)I(n)"),

	/* V from N/ADJ */
	"V.N.LA": suffix("lA"), /* kuru -> kurula  (dry -> to (make) dry) */

	/* N/ADJ from N/ADJ */
	"N.N.CI":  suffix("CI"),  /* person involved with noun */
	"N.N.LIK": suffix("lIK"), /* abstraction/object involved with noun */

	/* N/ADJ from V */
}

func suffix(s string) inf.Suffix {
	suf, ok := inf.ParseSuffix(s)
	if !ok {
		panic("failed to parse suffix")
	}
	return suf
}

/*
Returns the syllables of a word. Syllables are of the form CVCC where the onset always has priority.
Input should be lowercase.
*/
func Syllables(w []rune) [][]rune {
	syl_starts := []int{}
	if len(w) == 0 {
		return [][]rune{}
	}
	if inf.Vowel[w[0]] {
		syl_starts = append(syl_starts, 0)
	}

	for i := 1; i < len(w)-1; i++ {
		if (inf.Vowel[w[i]] && inf.Vowel[w[i-1]]) || (!inf.Vowel[w[i]] && inf.Vowel[w[i+1]]) {
			syl_starts = append(syl_starts, i)
		}
	}

	if len(w) >= 2 && inf.Vowel[w[len(w)-1]] && inf.Vowel[w[len(w)-2]] {
		syl_starts = append(syl_starts, len(w)-1)
	}

	syls := [][]rune{}

	for i := 0; i < len(syl_starts)-1; i++ {
		syls = append(syls, w[syl_starts[i]:syl_starts[i+1]])
	}
	syls = append(syls, w[syl_starts[len(syl_starts)-1]:])

	return syls
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Input root and suffixes:\n")
	if scanner.Scan() {
		root, sufs, ok := inf.ParseRootSuffixes(scanner.Text())
		if !ok {
			fmt.Printf("Error: failed to parse input\n")
			return
		}
		s := inf.Stem(root)
		for _, suf := range sufs {
			fmt.Printf("Stem: %s\nAdding suffix %s\n", s, suf)
			s = s.Append(suf)
		}
		fmt.Printf("Stem: %s\nWord: %s\n\n", s, s.Word())

	}
}
