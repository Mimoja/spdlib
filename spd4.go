package spdutil

import (
	"bytes"
	"encoding/binary"
	"github.com/sigurn/crc16"
	"log"
	"strings"
)

type SPD struct {
	SPDStatus                                         uint8
	SPDRevision                                       uint8
	DRAMDeviceType                                    uint8
	ModuleType                                        uint8
	DensityAndBanks                                   uint8
	Addressing                                        uint8
	PackageType                                       uint8
	OptionFeatures                                    uint8
	ThermalAndRefresh                                 uint8
	OtherOptionalFeatures                             uint8
	Reserved00_0                                      uint8
	NomalVoltage                                      uint8
	ModuleOrganization                                uint8
	ModuleMemoryBusWidth                              uint8
	ModuleThermalSensor                               uint8
	ExtendedModuleType                                uint8
	Reserved00_1                                      uint8
	Timebases                                         uint8
	MinimumCycleTime                                  uint8
	MaximumCycleTime                                  uint8
	CASLatencies                                      uint32
	MinimumCasLatency                                 uint8
	MinimumRAStoCASDelay                              uint8
	MinimumRowRechargeDelay                           uint8
	RASandRCminUppper                                 uint8
	ActiveToPrecharge_RAS_Lower                       uint8
	ActiveToActiveRefresh_RC_Lower                    uint8
	RefreshRecoveryDelay1                             uint16
	RefreshRecoveryDelay2                             uint16
	RefreshRecoveryDelay4                             uint16
	MinimumFourActiveWindow                           uint16
	MinimumActiveToActiveDelaySameBankGroup           uint8
	MinimumActiveToActiveDelayDifferentBankGroup      uint8
	MinimumCAStoCASDelay                              uint8
	Reserved00_2                                      [19]uint8
	ConnectorToRamMap                                 [18]uint8
	Reserved00_3                                      [116 - 78 + 1]uint8
	FineMinimumCAStoCASDelaySameBankGroup             uint8
	FineMinimumActiveToActiveDelaySameBankGroup       uint8
	FineMinimumAActiveToActiveDelayDifferentBankGroup uint8
	FineMinimumActiveToActiveRefreshDelay             uint8
	FineMinimumRowPrechargeDelay                      uint8
	FineMinimumRAStoCASDelay                          uint8
	FineMinimumCASLatencyTime                         uint8
	FineMaximumCycleTime                              uint8
	FineMinimumCycleTime                              uint8
	CRC                                               uint16
	ModuleSpecificParameter                           [191 - 127]byte
	HybridMemoryParameter                             [255 - 191]byte
	ExtendedFunctionParameter                         [319 - 255]byte
	ModuleManufactoringID                             uint16
	ModuleManufactoringLocation                       uint8
	ModuleManufactoringDate                           [2]byte
	ModuleManufactoringSerial                         uint32
	ModuleManufactoringPartNR                         [0x15C - 0x149 + 1]byte
	ModuleRevisionCode                                uint8
	ModuleManufactorID                                uint16
	ModuleStepping                                    uint8
	ModuleManufactoringData                           [0x17D - 0x161 + 1]byte
	ModuleReserved                                    [0x17F - 0x17E + 1]byte
	EndUserProgrammable                               [511 - 383]byte
}

type ParsedSPD struct {
	Raw                      SPD
	RawBytes                 []byte
	BytesTotal               uint
	BytesUsed                uint
	Revision                 uint8
	RamType                  string
	Vendor                   string
	CalculatedChecksum       uint16
	ManufactoringInformation struct {
	}
	ModulePartNumber string
}

func ParseSPD4(spdBytes []byte) ParsedSPD {
	var spd SPD
	var pspd ParsedSPD

	reader := bytes.NewReader(spdBytes)
	binary.Read(reader, binary.LittleEndian, &spd)

	pspd.RawBytes = spdBytes
	pspd.Raw = spd

	/** Byte 0 **/
	switch spd.SPDStatus >> 4 & 0b111 {
	case 0b00:
		break
	case 0b10:
		pspd.BytesTotal = 512
		break
	case 0b01:
		pspd.BytesTotal = 256
		break
	default:
		log.Fatal("BytesTotal is faulty")
	}

	switch spd.SPDStatus & 0b1111 {
	case 0b0000:
		break
	case 0b0001:
		pspd.BytesUsed = 128
		break
	case 0b0010:
		pspd.BytesUsed = 256
		break
	case 0b0011:
		pspd.BytesUsed = 384
		break
	case 0b0100:
		pspd.BytesUsed = 512
		break

	default:
		log.Fatal("BytesTotal is faulty")
	}

	/** Byte 1 */
	pspd.Revision = spd.SPDRevision

	/** Byte 2 */
	switch spd.DRAMDeviceType {
	case 0x0B:
		log.Fatal("This is a DDR3 SPD File")
		break
	case 0x0C:
		pspd.RamType = "SD-DDR4"
		break
	case 0x10:
		pspd.RamType = "LPDDR4"
		break
	case 0x11:
		pspd.RamType = "LPDDR4X"
	default:
		log.Fatalf("This is an unknown SPD File: 0x%02X", spd.DRAMDeviceType)
	}

	/*
		ModuleType uint8
		DensityAndBanks uint8
		Addressing uint8
		PackageType uint8
		OptionFeatures uint8
		ThermalAndRefresh uint8
		OtherOptionalFeatures uint8
		Reserved00_0 uint8
		NomalVoltage uint8
		ModuleOrganization uint8
		ModuleMemoryBusWidth uint8
		ModuleThermalSensor uint8
		ExtendedModuleType uint8
		Reserved00_1 uint8
		Timebases uint8
		MinimumCycleTime uint8
		MaximumCycleTime uint8
		CASLatencies uint32
		MinimumCasLatency uint8
		MinimumRAStoCASDelay uint8
		MinimumRowRechargeDelay uint8
		RASandRCminUppper uint8
		ActiveToPrecharge_RAS_Lower uint8
		ActiveToActiveRefresh_RC_Lower uint8
		RefreshRecoveryDelay1 uint16
		RefreshRecoveryDelay2 uint16
		RefreshRecoveryDelay4 uint16
		MinimumFourActiveWindow uint16
		MinimumActiveToActiveDelaySameBankGroup uint8
		MinimumActiveToActiveDelayDifferentBankGroup uint8
		MinimumCAStoCASDelay uint8
		Reserved00_2 [19]uint8
		ConnectorToRamMap [18]uint8
		Reserved00_3 [116 - 78 +1]uint8
		FineMinimumCAStoCASDelaySameBankGroup uint8
		FineMinimumActiveToActiveDelaySameBankGroup uint8
		FineMinimumAActiveToActiveDelayDifferentBankGroup uint8
		FineMinimumActiveToActiveRefreshDelay uint8
		FineMinimumRowPrechargeDelay uint8
		FineMinimumRAStoCASDelay uint8
		FineMinimumCASLatencyTime uint8
		FineMaximumCycleTime uint8
		FineMinimumCycleTime uint8
		CRC uint16
	*/

	pspd.CalculatedChecksum = crc16.Checksum(spdBytes[0:126], crc16.MakeTable(crc16.CRC16_XMODEM))

	vendor := "Unknown"
	switch pspd.Raw.ModuleManufactorID {
	case 0x2c80:
		vendor = "Crucial/Micron"
		break
	case 0x4304:
		vendor = "Ramaxel"
		break
	case 0x4f01:
		vendor = "Transcend"
		break
	case 0x9801:
		vendor = "Kingston"
		break
	case 0x987f:
		vendor = "Hynix"
		break
	case 0x9e02:
		vendor = "Corsair"
		break
	case 0xb004:
		vendor = "OCZ"
		break
	case 0xad80:
		vendor = "Hynix/Hyundai"
		break
	case 0xb502:
		vendor = "SuperTalent"
		break
	case 0xcd04:
		vendor = "GSkill"
		break
	case 0xce80:
		vendor = "Samsung"
		break
	case 0xfe02:
		vendor = "Elpida"
		break
	case 0xff2c:
		vendor = "Micron"
		break
	}
	pspd.Vendor = vendor

	pspd.ModulePartNumber = string(pspd.Raw.ModuleManufactoringPartNR[:])
	pspd.ModulePartNumber = strings.Trim(pspd.ModulePartNumber, "\000")
	return pspd
}
