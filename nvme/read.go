package nvme

import (
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"
)

func (d *NVMeDevice) Read(lba uint64, length uint16, buf []byte) error {

	cmd := nvmeUserIo{
		opcode:  NVME_NVM_CMD_READ,
		slba:    lba,
		addr:    uint64(uintptr(unsafe.Pointer(&(buf)[0]))),
		nblocks: length - 1,
	}

	return ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_SUBMIT_IO, uintptr(unsafe.Pointer(&cmd)))
}
