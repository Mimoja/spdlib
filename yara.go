package spdutil

import (
	"github.com/hillu/go-yara"
	"log"
	"os"
)

var yaraRules *yara.Rules

func setupYara() {
	c, err := yara.NewCompiler()
	if err != nil {
		panic("Could not create yara compiler")
	}

	file, err := os.Open("rule.yara")
	if err != nil {
		log.Fatalf("Could not load rules: %v", err)
	}

	c.AddFile(file, "test")

	r, err := c.GetRules()
	if err != nil {
		log.Fatalf("Failed to compile rules: %s", err)
	}
	yaraRules = r
}

type SPDMatch struct {
	Offset uint64
	SPD    []byte
}

func FindSPDs(bs []byte) (spds []SPDMatch) {
	setupYara()

	matches, err := yaraRules.ScanMem(bs, 0, 0, nil)
	if err != nil {
		log.Fatal("could not scan with yara %v\n", err)
		return
	}

	if len(matches) == 0 {
		log.Println("Could not find any matches!")
	}

	for _, match := range matches {
		for _, m := range match.Strings {
			log.Printf("Found: %s : %s at 0x%X", match.Rule, m.Name[1:], m.Offset)
			spds = append(spds, SPDMatch{
				Offset: m.Offset,
				SPD:    bs[m.Offset : m.Offset+512],
			})
		}

	}
	return spds
}
