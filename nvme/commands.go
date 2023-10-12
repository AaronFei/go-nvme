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

const (
	// cf. NVM Express Base Specification 2.0c , section 5: Admin Command Set
	NVME_ADMIN_GET_LOG_PAGE           uint8 = 0x02
	NVME_ADMIN_IDENTIFY               uint8 = 0x06
	NVME_ADMIN_ABORT                  uint8 = 0x08
	NVME_ADMIN_SET_FEATURES           uint8 = 0x09
	NVME_ADMIN_GET_FEATURES           uint8 = 0x0a
	NVME_ADMIN_ASYNC_EV_REQ           uint8 = 0x0c
	NVME_ADMIN_NS_MANAGEMENT          uint8 = 0x0d
	NVME_ADMIN_FIRMWARE_COMMIT        uint8 = 0x10
	NVME_ADMIN_FIRMWARE_DOWNLOAD      uint8 = 0x11
	NVME_ADMIN_DEVICE_SELF_TEST       uint8 = 0x14
	NVME_ADMIN_NS_ATTACHMENT          uint8 = 0x15
	NVME_ADMIN_KEEP_ALIVE             uint8 = 0x18
	NVME_ADMIN_DIRECTIVE_SEND         uint8 = 0x19
	NVME_ADMIN_DIRECTIVE_RECV         uint8 = 0x1a
	NVME_ADMIN_VIRTUALIZATION_MGMT    uint8 = 0x1c
	NVME_ADMIN_NVME_MI_SEND           uint8 = 0x1d
	NVME_ADMIN_NVME_MI_RECV           uint8 = 0x1e
	NVME_ADMIN_CAPACITY_MGMT          uint8 = 0x20
	NVME_ADMIN_LOCKDOWN               uint8 = 0x24
	NVME_ADMIN_DOORBELL_BUFFER_CONFIG uint8 = 0x7c
	NVME_ADMIN_FORMAT_NVM             uint8 = 0x80
	NVME_ADMIN_SECURITY_SEND          uint8 = 0x81
	NVME_ADMIN_SECURITY_RECV          uint8 = 0x82
	NVME_ADMIN_SANITIZE               uint8 = 0x84
	NVME_ADMIN_GET_LBA_STATUS         uint8 = 0x86

	NVME_NVM_CMD_FLUSH                uint8 = 0x00
	NVME_NVM_CMD_WRITE                uint8 = 0x01
	NVME_NVM_CMD_READ                 uint8 = 0x02
	NVME_NVM_CMD_WRITE_UNCORRECTABLE  uint8 = 0x04
	NVME_NVM_CMD_COMPARE              uint8 = 0x05
	NVME_NVM_CMD_WRITE_ZEROES         uint8 = 0x08
	NVME_NVM_CMD_DATASET_MANAGEMENT   uint8 = 0x09
	NVME_NVM_CMD_VERIFY               uint8 = 0x0c
	NVME_NVM_CMD_RESERVATION_REGISTER uint8 = 0x0d
	NVME_NVM_CMD_RESERVATION_REPORT   uint8 = 0x0e
	NVME_NVM_CMD_RESERVATION_ACQUIRE  uint8 = 0x11
	NVME_NVM_CMD_RESERVATION_RELEASE  uint8 = 0x15
	NVME_NVM_CMD_COPY                 uint8 = 0x19
)

type nvmeAdminCmd nvmePassthruCommand

// Defined in <linux/nvme_ioctl.h> (first 64 bytes refer to NVM Express Base Specification 2.0c,
// figure 88: Common Command Format - Admin and NVM Vendor Specific Commands)
type nvmePassthruCommand struct {
	opcode       uint8
	flags        uint8
	rsvd1        uint16
	nsid         uint32
	cdw2         uint32
	cdw3         uint32
	metadata     uint64
	addr         uint64
	metadata_len uint32
	data_len     uint32
	cdw10        uint32
	cdw11        uint32
	cdw12        uint32
	cdw13        uint32
	cdw14        uint32
	cdw15        uint32
	timeout_ms   uint32
	result       uint32
} // 72 bytes

type nvmeUserIo struct {
	opcode   uint8
	flags    uint8
	control  uint16
	nblocks  uint16
	rsvd     uint16
	metadata uint64
	addr     uint64
	slba     uint64
	dsmgmt   uint32
	reftag   uint32
	apptag   uint16
	appmask  uint16
}
