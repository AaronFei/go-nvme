package nvme

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"
)

const (
	IDENTIFY_CNS_NSID           uint8 = 0x00
	IDENTIFY_CNS_CTRL           uint8 = 0x01
	IDENTIFY_CNS_ACTIVE_NS_LIST uint8 = 0x02
	IDENTIFY_CNS_NS_ID_DESC     uint8 = 0x03
	IDENTIFY_CNS_NVM_SET_LIST   uint8 = 0x04
	IDENTIFY_CNS_IOCS_NS        uint8 = 0x05
	IDENTIFY_CNS_IOCS_CTRL      uint8 = 0x06
	IDENTIFY_CNS_IOCS_ACTIVE_NS uint8 = 0x07
	IDENTIFY_CNS_IOCS_INDEP_NS  uint8 = 0x08
	IDENTIFY_CNS_ALLOC_NS_LIST  uint8 = 0x10
	IDENTIFY_CNS_ALLOC_NS       uint8 = 0x11
	IDENTIFY_CNS_CTRL_LIST_NS   uint8 = 0x12
	IDENTIFY_CNS_CTRL_LIST      uint8 = 0x13
	IDENTIFY_CNS_PRIMARY_CTRL   uint8 = 0x14
	IDENTIFY_CNS_SECONDARY_CTRL uint8 = 0x15
	IDENTIFY_CNS_NS_GRAN_LIST   uint8 = 0x16
	IDENTIFY_CNS_UUID_LIST      uint8 = 0x17
	IDENTIFY_CNS_DOMAIN_LIST    uint8 = 0x18
	IDENTIFY_CNS_ENDURANCE_LIST uint8 = 0x19
	IDENTIFY_CNS_IOCS_ALLOC_NS  uint8 = 0x1a
	IDENTIFY_CNS_IOCS_ALLOC     uint8 = 0x1b
	IDENTIFY_CNS_IOCS           uint8 = 0x1c
)

type NvmeIdentNamespace struct {
	Nsze    uint64
	Ncap    uint64
	Nuse    uint64
	Nsfeat  uint8
	Nlbaf   uint8
	Flbas   uint8
	Mc      uint8
	Dpc     uint8
	Dps     uint8
	Nmic    uint8
	Rescap  uint8
	Fpi     uint8
	Rsvd33  uint8
	Nawun   uint16
	Nawupf  uint16
	Nacwu   uint16
	Nabsn   uint16
	Nabo    uint16
	Nabspf  uint16
	Rsvd46  [2]byte
	Nvmcap  [16]byte
	Rsvd64  [40]byte
	Nguid   [16]byte
	EUI64   [8]byte
	Lbaf    [16]lbaf
	Rsvd192 [192]byte
	Vs      [3712]byte
} // 4096 bytes

type lbaf struct {
	Ms    uint16
	Lbads uint8
	Rp    uint8
}

func (d *NVMeDevice) IdentifyRaw(cns uint8, nsid uint32, cdw10 uint32, cdw11 uint32, cdw14 uint32, buf []byte) error {
	cmd := nvmePassthruCommand{
		opcode:   NVME_ADMIN_IDENTIFY,
		nsid:     nsid, // Namespace 0, since we are identifying the controller
		addr:     uint64(uintptr(unsafe.Pointer(&(buf[0])))),
		data_len: uint32(len(buf)),
		cdw10:    cdw10, // Identify controller
		cdw11:    cdw11,
		cdw14:    cdw14,
	}

	if err := ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_ADMIN_CMD, uintptr(unsafe.Pointer(&cmd))); err != nil {
		return err
	}

	return nil
}

func (d *NVMeDevice) IdentifyController() (NvmeIdentController, error) {
	buf := make([]byte, 4096)

	if err := d.IdentifyRaw(IDENTIFY_CNS_CTRL, 0, 1, 0, 0, buf); err != nil {
		return NvmeIdentController{}, err
	}

	var idCtrlr NvmeIdentController

	binary.Read(bytes.NewBuffer(buf[:]), NativeEndian, &idCtrlr)

	d.ModelInfo = NvmeController{
		VendorID:        idCtrlr.VendorID,
		ModelNumber:     string(idCtrlr.ModelNumber[:]),
		SerialNumber:    string(bytes.TrimSpace(idCtrlr.SerialNumber[:])),
		FirmwareVersion: string(idCtrlr.Firmware[:]),
		MaxDataXferSize: 1 << idCtrlr.Mdts,
		// Convert IEEE OUI ID from big-endian
		OUI: uint32(idCtrlr.IEEE[0]) | uint32(idCtrlr.IEEE[1])<<8 | uint32(idCtrlr.IEEE[2])<<16,
	}

	return idCtrlr, nil
}

func (d *NVMeDevice) IdentifyNamespace(nsid uint32) (NvmeIdentNamespace, error) {
	buf := make([]byte, 4096)

	if err := d.IdentifyRaw(IDENTIFY_CNS_NSID, nsid, 0, 0, 0, buf); err != nil {
		return NvmeIdentNamespace{}, err
	}

	var ns NvmeIdentNamespace
	binary.Read(bytes.NewBuffer(buf[:]), NativeEndian, &ns)

	return ns, nil
}
