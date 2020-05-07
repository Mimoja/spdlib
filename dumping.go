package spdutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func WriteSPD4(pspd ParsedSPD, filename string) {
	file, err := os.OpenFile(filename+".spd.hex", os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		log.Fatal("Could not open file for writing")
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("# TotalBytes: %d ; BytesUsed: %d\n", pspd.BytesTotal, pspd.BytesUsed))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.SPDStatus))

	file.WriteString(fmt.Sprintf("# SPD Revision %X.%X\n", pspd.Revision>>4, pspd.Revision&0xF))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.SPDRevision))

	file.WriteString(fmt.Sprintf("# DDR Ramtype: %s\n", pspd.RamType))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.DRAMDeviceType))

	file.WriteString(fmt.Sprintf("# Config Rest\n"))
	file.WriteString(Dump(pspd.RawBytes[3:126]))

	matching := "Not Matching"
	if pspd.CalculatedChecksum == pspd.Raw.CRC {
		matching = "Match!"
	}

	file.WriteString(fmt.Sprintf("# CRC Is: 0x%X Calculated: 0x%X %s\n",
		pspd.Raw.CRC, pspd.CalculatedChecksum, matching))
	file.WriteString(fmt.Sprintf("%04X\n", pspd.Raw.CRC))

	file.WriteString(fmt.Sprintf("\n# ModuleSpecificParameter\n"))
	file.WriteString(Dump(pspd.Raw.ModuleSpecificParameter[:]))

	file.WriteString(fmt.Sprintf("# HybridMemoryParameter\n"))
	file.WriteString(Dump(pspd.Raw.HybridMemoryParameter[:]))

	file.WriteString(fmt.Sprintf("# ExtendedFunctionParameter\n"))
	file.WriteString(Dump(pspd.Raw.ExtendedFunctionParameter[:]))

	file.WriteString(fmt.Sprintf("# ManufactoringInformation\n"))

	file.WriteString(fmt.Sprintf("## Module Manufactoring ID\n"))
	file.WriteString(fmt.Sprintf("%04X\n", pspd.Raw.ModuleManufactoringID))

	file.WriteString(fmt.Sprintf("## Module Manufactoring Location and Date\n"))
	file.WriteString(fmt.Sprintf("%02X %s", pspd.Raw.ModuleManufactoringLocation, Dump(
		pspd.Raw.ModuleManufactoringDate[:])))

	file.WriteString(fmt.Sprintf("## Module Manufactoring Serial\n"))
	file.WriteString(fmt.Sprintf("%08X\n", pspd.Raw.ModuleManufactoringSerial))

	file.WriteString(fmt.Sprintf("## Module Manufactoring Part Number: \"%s\"\n",
		pspd.ModulePartNumber))
	file.WriteString(Dump(pspd.Raw.ModuleManufactoringPartNR[:]))

	file.WriteString(fmt.Sprintf("## Module Manufactoring Revision Code\n"))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.ModuleRevisionCode))

	file.WriteString(fmt.Sprintf("## Module Manufactor: \"%s\" (0x%04X)\n", pspd.Vendor, pspd.Raw.ModuleManufactorID))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.ModuleManufactorID))

	file.WriteString(fmt.Sprintf("## Module Stepping\n"))
	file.WriteString(fmt.Sprintf("%02X\n", pspd.Raw.ModuleStepping))

	file.WriteString(fmt.Sprintf("## Module Manufactoring Data\n"))
	file.WriteString(Dump(pspd.Raw.ModuleManufactoringData[:]))

	file.WriteString(fmt.Sprintf("## Module Reserved\n"))
	file.WriteString(Dump(pspd.Raw.ModuleReserved[:]))

	file.WriteString(fmt.Sprintf("\n# EndUserProgrammable\n"))
	file.WriteString(Dump(pspd.Raw.EndUserProgrammable[:]))

	ioutil.WriteFile(filename+".spd.bin", pspd.RawBytes, os.ModePerm)

	validityChecks := []bool{
		pspd.Raw.CRC == pspd.CalculatedChecksum,
		pspd.Raw.ModuleManufactorID != 0,
		string(pspd.Raw.ModuleManufactoringPartNR[:]) != "                    ",
	}

	validity := 0.0
	for _, v := range validityChecks {
		if v {
			validity += 1.0
		}
	}
	validity = validity / float64(len(validityChecks))

	fmt.Printf("%s has a validity of %03f\n", filename, validity)
}

func Dump(bs []byte) string {
	dump := ""
	pos := 1
	for _, b := range bs {
		dump += fmt.Sprintf("%02X", b)
		if pos > 15 {
			dump += "\n"
			pos = 0
		} else {
			dump += " "
		}
		pos++
	}
	return dump + "\n"
}
