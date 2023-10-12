package nvme

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
)

type nvmeSMARTLog struct {
	CritWarning      uint8
	Temperature      [2]uint8
	AvailSpare       uint8
	SpareThresh      uint8
	PercentUsed      uint8
	Rsvd6            [26]byte
	DataUnitsRead    [16]byte
	DataUnitsWritten [16]byte
	HostReads        [16]byte
	HostWrites       [16]byte
	CtrlBusyTime     [16]byte
	PowerCycles      [16]byte
	PowerOnHours     [16]byte
	UnsafeShutdowns  [16]byte
	MediaErrors      [16]byte
	NumErrLogEntries [16]byte
	WarningTempTime  uint32
	CritCompTime     uint32
	TempSensor       [8]uint16
	Rsvd216          [296]byte
} // 512 bytes

func (d *NVMeDevice) GetLogPageSmart(buf []byte) error {
	cdw10 := buildCdw(LogPageCdw10BitInfo, LogPageCdw10{
		LID:   uint32(LOGPAGE_SMART_HEALTH_INFO),
		NUMDL: ((uint32(len(buf)) / 4) - 1),
	})

	if err := d.GetLogPageRaw(0xffffffff, cdw10, 0, 0, 0, 0, buf); err != nil {
		return err
	}

	return nil
}

func (d *NVMeDevice) PrintSMART(w io.Writer) error {
	buf := make([]byte, 512)

	// Read SMART log
	if err := d.GetLogPageSmart(buf); err != nil {
		return err
	}

	var sl nvmeSMARTLog

	binary.Read(bytes.NewBuffer(buf[:]), NativeEndian, &sl)

	unitsRead := le128ToBigInt(sl.DataUnitsRead)
	unitsWritten := le128ToBigInt(sl.DataUnitsWritten)
	unit := big.NewInt(512 * 1000)

	fmt.Fprintln(w, "\nSMART data follows:")
	fmt.Fprintf(w, "Critical warning: %#02x\n", sl.CritWarning)
	fmt.Fprintf(w, "Temperature: %d° Celsius\n",
		(uint16(sl.Temperature[0])|uint16(sl.Temperature[1])<<8)-273) // Kelvin to degrees Celsius
	fmt.Fprintf(w, "Avail. spare: %d%%\n", sl.AvailSpare)
	fmt.Fprintf(w, "Avail. spare threshold: %d%%\n", sl.SpareThresh)
	fmt.Fprintf(w, "Percentage used: %d%%\n", sl.PercentUsed)
	fmt.Fprintf(w, "Data units read: %d [%s]\n",
		unitsRead, formatBigBytes(new(big.Int).Mul(unitsRead, unit)))
	fmt.Fprintf(w, "Data units written: %d [%s]\n",
		unitsWritten, formatBigBytes(new(big.Int).Mul(unitsWritten, unit)))
	fmt.Fprintf(w, "Host read commands: %d\n", le128ToBigInt(sl.HostReads))
	fmt.Fprintf(w, "Host write commands: %d\n", le128ToBigInt(sl.HostWrites))
	fmt.Fprintf(w, "Controller busy time: %d\n", le128ToBigInt(sl.CtrlBusyTime))
	fmt.Fprintf(w, "Power cycles: %d\n", le128ToBigInt(sl.PowerCycles))
	fmt.Fprintf(w, "Power on hours: %d\n", le128ToBigInt(sl.PowerOnHours))
	fmt.Fprintf(w, "Unsafe shutdowns: %d\n", le128ToBigInt(sl.UnsafeShutdowns))
	fmt.Fprintf(w, "Media & data integrity errors: %d\n", le128ToBigInt(sl.MediaErrors))
	fmt.Fprintf(w, "Error information log entries: %d\n", le128ToBigInt(sl.NumErrLogEntries))

	return nil
}
