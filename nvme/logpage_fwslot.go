package nvme

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type LogPageFwSlotInfo struct {
	AFI        uint8
	Reserved1  [7]byte
	FWRevision [7][8]byte
	Reserved2  [448]byte
}

func (d *NVMeDevice) GetLogPageFwSlotInfo(buf []byte) error {
	cdw10 := buildCdw(LogPageCdw10BitInfo, LogPageCdw10{
		LID:   uint32(LOGPAGE_FIRMWARE_SLOT_INFO),
		NUMDL: ((uint32(len(buf)) / 4) - 1),
	})

	if err := d.GetLogPageRaw(0, cdw10, 0, 0, 0, 0, buf); err != nil {
		return err
	}

	return nil
}

func (d *NVMeDevice) PrintFwSlotInfo() error {
	buf := make([]byte, 512)

	// Read SMART log
	if err := d.GetLogPageFwSlotInfo(buf); err != nil {
		return err
	}

	var sl LogPageFwSlotInfo

	binary.Read(bytes.NewBuffer(buf[:]), NativeEndian, &sl)

	fmt.Printf("\nActive Firmware Slot follows:\n")

	slotNum := (sl.AFI & 0x7)
	if slotNum == 0 {
		fmt.Println("Active Firmware Info: Invalid")
	} else {
		fmt.Printf("Active Firmware Info: %#02b\n", sl.AFI)
		fmt.Printf("Firmware current Revision: %s\n", string(sl.FWRevision[slotNum-1][:]))
		for i := 0; i < 7; i++ {
			fmt.Printf("Firmware Revision for slot %d: %s\n", i+1, string(sl.FWRevision[i][:]))
		}
	}

	return nil
}
