package main

import (
	"flag"
	"fmt"
	"github.com/mimoja/spdutil"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	outputFolder := flag.String("o", ".", "Output folder")
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("No input file provided")
		os.Exit(1)
	}
	fileName := flag.Args()[0]
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		println("Could not read " + fileName + ": " + err.Error())
		os.Exit(2)
	}
	ds := spdutil.FindSPDs(bs)

	if len(ds) == 0 {
		println("No DDR4 SPDs found!")
		os.Exit(0)
	}

	for _, d := range ds {
		pspd := spdutil.ParseSPD4(d.SPD)
		spdutil.WriteSPD4(pspd, filepath.Join(*outputFolder, fmt.Sprintf("%08X", d.Offset)))
	}
}
