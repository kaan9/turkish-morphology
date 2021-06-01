package main

import (
	"bufio"
	"fmt"
	"os"
	inf "github.com/kaan9/turkish-morphology/inflection"
)



var suffixes = map[string]inf.Suffix{
	/* tense/aspect (does not include -makta, which can be encoded as as -mak + -ta */
	"known past": inf.Suffix{Head: 0, Tail: 0, Body: []rune("DI")},
	"infer past": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mIş")},
	"aorist a":   inf.Suffix{Head: 'A', Tail: 0, Body: []rune("r")},
	"aorist i":   inf.Suffix{Head: 'I', Tail: 0, Body: []rune("r")},
	"aorist neg": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAz")},
	"pres cont":  inf.Suffix{Head: 0, Tail: 0, Body: []rune("Iyor")},
	"fut":        inf.Suffix{Head: 'y', Tail: 0, Body: []rune("AcAK")},
	/* verb negation */
	"neg": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mA")},
	/* infinitive */
	"inf": inf.Suffix{Head: 0, Tail: 0, Body: []rune("mAK")},
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
