package nvme

import (
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"
)

func (d *NVMeDevice) Write(lba uint64, length uint16, write_hint uint32, buf []byte) error {

	cmd := nvmeUserIo{
		opcode:  NVME_NVM_CMD_WRITE,
		slba:    lba,
		addr:    uint64(uintptr(unsafe.Pointer(&(buf)[0]))),
		nblocks: length - 1,
	}

	return ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_SUBMIT_IO, uintptr(unsafe.Pointer(&cmd)))
}
