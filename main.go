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
	"PL":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},

	/* possessive (iyelik) */
	"POS.1sg":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("m")},
	"POS.1pl":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("mIz")},
	"POS.2sg":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("n")},
	"POS.2pl":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("nIz")},
	"POS.3sg":	inf.Suffix{Head: 's', Tail: 'n', Body: []rune("I")},
	"POS.3pl":	inf.Suffix{Head: 0, Tail: 'n', Body: []rune("lArI")},

	/* the familial -gil and -ler suffix: e.g. teyzemler, karıncayiyengiller */
	"FAML":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("gil")},	/* no consonant/vowel harmony */
	"FAML.PL":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},	/* not the same as PL */

	/* case (all except def. accusative can be come before predicative personal suffixes?) */
	"ABSL":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("")},		/* Absolute (yalın) case */
	"ACC.DEF":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("I")},	/* Definite accusative */
	"DAT":		inf.Suffix{Head: 'y', Tail: 0, Body: []rune("A")},	/* dative-directional/lative */
	"GEN":		inf.Suffix{Head: 'n', Tail: 0, Body: []rune("In")},
	"LOC":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("DA")},
	"ABL":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("DAn")},
	"INS":		inf.Suffix{Head: 'y', Tail: 0, Body: []rune("lA")},	/* also postposition 'ile' */


	/* Personal Suffixes (kişi ekleri) */
	/* Predicative Personal Suffix - type I (Copular and after -mIş -AcAK -(A/I)r -Iyor ... other forms) */
	"PRED.1sg":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("Im")},
	"PRED.1pl":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("Iz")},
	"PRED.2sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("sIn")},
	"PRED.2pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("sInIz")},
	"PRED.3sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("")},
	"PRED.3pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},
	/* Verbal Personal Suffix -- type II (after -DI and -sA) */
	"VB.1sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("m")},
	"VB.1pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("k")},
	"VB.2sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("n")},
	"VB.2pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("nIz")},
	"VB.3sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("")},
	"VB.3pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("lAr")},
	/* Optative Personal Suffix -- type III (optative mood) */
	"OPT.1sg":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AyIm")},
	"OPT.1pl":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AlIm")},
	"OPT.2sg":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AsIn")},
	"OPT.2pl":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AsInIz")},
	"OPT.3sg":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("A")},
	"OPT.3pl":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AlAr")},
	/* Imperative Personal Suffix -- type IV (imperative mood) */
	"IMP.2sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("")},
	"IMP.2pl":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("In")},
	"IMP.2pl2":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("InIz")},	/* more formal */
	"IMP.3sg":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("sIn")},
	"IMP.3pl":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("sInlAr")},

	/* tense/aspect (does not include PRS.PROG -makta, which can be encoded as as -mak + -ta) */
	"PPFV":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("DI")},	/* past perfective */
	"PPFV.INFR":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("mIş")},	/* inferred past perfective */
	"AOR.A":	inf.Suffix{Head: 'A', Tail: 0, Body: []rune("r")},	/* aorist low vowel */
	"AOR.I":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("r")},	/* aorist high vowel */
	"AOR.NEG":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAz")},	/* aorist negative */
	"AOR.INAB":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AmAz")},	/* aorist impotential */
	/* are these irregular? yapmam, yapmayız yapamam(yapamazım) yapamayız(yapamazız) */
	"PRS.IPFV":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("Iyor")},	/* present imperfective */
	"PRS.PROG":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAktA")},	/* pres. progressive: -mAK + -DA */
	"FUT":		inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},	/* future */
	/* mood */
	"COND":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("sA")},	/* conditional */
	"NEC":		inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAlI")},	/* necessitative */

	/* verb negation */
	"NEG": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mA")},

	/* copula (comes after the same suffixes as the predicative personal suffixes (type I)) */
	"COP": inf.Suffix{Head: 0, Tail: 0, Body: []rune("DIr")},

	/* verbal noun */
	"INF": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAK")},		/* infinitive */
	"GER": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mA")},		/* gerund */
	"WAY": inf.Suffix{Head: 'y', Tail: 0, Body: []rune("Iş")},		/* 'way/act of doing' verb */

	/* interrogative particle */
	"INT": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mI")},		/* written separate by convention*/


	/* grammatical voice */


	/* Participles (separated as personal (takes suffix of possession) versus impersonal) */
	"PTCP.IMPRS.AOR.A":	inf.Suffix{Head: 'A', Tail: 0, Body: []rune("r")},	/* aorist low vowel */
	"PTCP.IMPRS.AOR.I":	inf.Suffix{Head: 'I', Tail: 0, Body: []rune("r")},	/* aorist high vowel */
	"PTCP.IMPRS.AOR.NEG":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAz")},	/* aorist negative */
	"PTCP.IMPRS.AOR.INAB":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AmAz")},	/* aorist impotential */
	"PTCP.IMPRS.IPFV":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("An")},	/* imperfective */
	"PTCP.IMPRS.FUT":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},	/* impersonal future */
	"PTCP.PERS.FUT":	inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},	/* personal future */
	"PTCP.IMPRS.PFV":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("mIş")},	/* impersonal perfective */
	"PTCP.PERS.PFV":	inf.Suffix{Head: 0, Tail: 0, Body: []rune("DIK")},	/* personal perfective */

	// is the converb suffix -Ip also a participle?    yapıp



}


/*
Break word up into syllables. Syllables are of the form (C)V((G)C) where the onset always has
priority.
*/
func Syllables(s string) []string{
	return nil
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
