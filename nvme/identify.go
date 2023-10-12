package nvme

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"
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

func (d *NVMeDevice) IdentifyController() (NvmeIdentController, error) {
	buf := make([]byte, 4096)

	cmd := nvmePassthruCommand{
		opcode:   NVME_ADMIN_IDENTIFY,
		nsid:     0, // Namespace 0, since we are identifying the controller
		addr:     uint64(uintptr(unsafe.Pointer(&(buf[0])))),
		data_len: uint32(len(buf)),
		cdw10:    1, // Identify controller
	}

	if err := ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_ADMIN_CMD, uintptr(unsafe.Pointer(&cmd))); err != nil {
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

func (d *NVMeDevice) IdentifyNamespace(namespace uint32) (NvmeIdentNamespace, error) {
	buf := make([]byte, 4096)

	cmd := nvmePassthruCommand{
		opcode:   NVME_ADMIN_IDENTIFY,
		nsid:     namespace,
		addr:     uint64(uintptr(unsafe.Pointer(&buf[0]))),
		data_len: uint32(len(buf)),
		cdw10:    0,
	}

	if err := ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_ADMIN_CMD, uintptr(unsafe.Pointer(&cmd))); err != nil {
		return NvmeIdentNamespace{}, err
	}

	var ns NvmeIdentNamespace
	binary.Read(bytes.NewBuffer(buf[:]), NativeEndian, &ns)

	d.namespace[namespace] = ns

	return ns, nil
}

func (d *NVMeDevice) getLbaSize(nsid uint32) (uint64, error) {
	ns, err := d.IdentifyNamespace(nsid)
	if err != nil {
		return 0, err
	}

	lbaf := ns.Lbaf[getBitsValue(uint64(ns.Flbas), 3, 0)]

	return uint64(1 << lbaf.Lbads), nil
}
