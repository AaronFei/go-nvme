// Copyright 2017-2022 Daniel Swarbrick. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nvme

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"

	"golang.org/x/sys/unix"
)

var (
	// Defined in <linux/nvme_ioctl.h>
	NVME_IOCTL_ADMIN_CMD = ioctl.Iowr('N', 0x41, unsafe.Sizeof(nvmeAdminCmd{}))
	NVME_IOCTL_SUBMIT_IO = ioctl.Iow('N', 0x42, unsafe.Sizeof(nvmeUserIo{}))
	NVME_IOCTL_IO_CMD    = ioctl.Iowr('N', 0x43, unsafe.Sizeof(nvmePassthruCommand{}))
)

// NVMeController encapsulates the attributes of an NVMe controller.
type NvmeController struct {
	VendorID        uint16
	ModelNumber     string
	SerialNumber    string
	FirmwareVersion string
	OUI             uint32 // IEEE OUI identifier
	MaxDataXferSize uint
}

type NVMeDevice struct {
	Name      string
	fd        int
	ModelInfo NvmeController
}

func NewNVMeDevice(name string) *NVMeDevice {
	return &NVMeDevice{Name: name}
}

func (d *NVMeDevice) Open() (err error) {
	d.fd, err = unix.Open(d.Name, unix.O_RDWR, 0600)
	return err
}

func (d *NVMeDevice) Close() error {
	return unix.Close(d.fd)
}

// Print outputs the attributes of an NVMe controller in a pretty-print style.
func (c *NVMeDevice) IdentPrint(w io.Writer) {
	c.IdentifyController()
	fmt.Fprintf(w, "Vendor ID          : %#04x\n", c.ModelInfo.VendorID)
	fmt.Fprintf(w, "Model number       : %s\n", c.ModelInfo.ModelNumber)
	fmt.Fprintf(w, "Serial number      : %s\n", c.ModelInfo.SerialNumber)
	fmt.Fprintf(w, "Firmware version   : %s\n", c.ModelInfo.FirmwareVersion)
	fmt.Fprintf(w, "IEEE OUI identifier: %#06x\n", c.ModelInfo.OUI)
	fmt.Fprintf(w, "Max. data xfer size: %d pages\n", c.ModelInfo.MaxDataXferSize)
}
