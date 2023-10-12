package nvme

import (
	"fmt"
	"unsafe"

	"github.com/AaronFei/go-nvme/ioctl"
)

const (
	LOGPAGE_SUPPORTED_LOG_PAGES                 uint8 = 0x00
	LOGPAGE_ERROR_INFO                          uint8 = 0x01
	LOGPAGE_SMART_HEALTH_INFO                   uint8 = 0x02
	LOGPAGE_FIRMWARE_SLOT_INFO                  uint8 = 0x03
	LOGPAGE_CHANGED_NS_LIST                     uint8 = 0x04
	LOGPAGE_CMD_SUPPORTED_EFFECTS               uint8 = 0x05
	LOGPAGE_DEVICE_SELF_TEST                    uint8 = 0x06
	LOGPAGE_TELEMETRY_HOST                      uint8 = 0x07
	LOGPAGE_TELEMETRY_CTRL                      uint8 = 0x08
	LOGPAGE_ENDURANCE_GROUP_INFO                uint8 = 0x09
	LOGPAGE_PREDICTABLE_LATENCY_PER_NVM_SET     uint8 = 0x0a
	LOGPAGE_PREDICTABLE_LATENCY_EVENT_AGGREGATE uint8 = 0x0b
	LOGPAGE_ASYMMETRIC_NS_ACCESS                uint8 = 0x0c
	LOGPAGE_PERSISTENT_EVENT_LOG                uint8 = 0x0d
	LOGPAGE_ENDURANCE_GROUP_EVENT_AGGREGATE     uint8 = 0x0f
	LOGPAGE_MEDIA_UNIT_STATUS                   uint8 = 0x10
	LOGPAGE_SUPPORTED_CAPACITY_CONFIG_LIST      uint8 = 0x11
	LOGPAGE_FEATURE_IDENTIFIERS                 uint8 = 0x12
	LOGPAGE_NVME_MI_CMD_SUPPORTED_EFFECTS       uint8 = 0x13
	LOGPAGE_CMD_FEATURE_LOCKDOWN                uint8 = 0x14
	LOGPAGE_BOOT_PARTITION                      uint8 = 0x15
	LOGPAGE_ROTATIONAL_MEDIA_INFO               uint8 = 0x16
	LOGPAGE_DISCOVERY                           uint8 = 0x70
	LOGPAGE_RESERVATION_NOTIFICATION            uint8 = 0x80
	LOGPAGE_SANITIZE_STATUS                     uint8 = 0x81
)

var LogPageCdw10BitInfo = cdwBitInfo{
	{
		name: "LID", bitStart: 0,
	},
	{
		name: "LSP", bitStart: 8,
	},
	{
		name: "RAE", bitStart: 15,
	},
	{
		name: "NUMDL", bitStart: 16,
	},
}

type LogPageCdw10 struct {
	LID   uint32
	LSP   uint32
	RAE   uint32
	NUMDL uint32
}

var LogPageCdw11BitInfo = cdwBitInfo{
	{
		name: "NUMDU", bitStart: 0,
	},
	{
		name: "LSID", bitStart: 16,
	},
}

type LogPageCdw11 struct {
	NUMDU uint32
	LSID  uint32
}

var LogPageCdw12BitInfo = cdwBitInfo{
	{
		name: "LPOL", bitStart: 0,
	},
}

type LogPageCdw12 struct {
	LPOL uint32
}

var LogPageCdw13BitInfo = cdwBitInfo{
	{
		name: "LPOU", bitStart: 0,
	},
}

type LogPageCdw13 struct {
	LPOU uint32
}

var LogPageCdw14BitInfo = cdwBitInfo{
	{
		name: "UUID", bitStart: 0,
	},
	{
		name: "OT", bitStart: 23,
	},
	{
		name: "CSI", bitStart: 24,
	},
}

type LogPageCdw14 struct {
	UUID uint32
	OT   uint32
	CSI  uint32
}

func (d *NVMeDevice) GetLogPageRaw(nsid, cdw10, cdw11, cdw12, cdw13, cdw14 uint32, buf []byte) error {
	cmd := nvmePassthruCommand{
		opcode:   NVME_ADMIN_GET_LOG_PAGE,
		nsid:     nsid,
		addr:     uint64(uintptr(unsafe.Pointer(&(buf)[0]))),
		data_len: uint32(len(buf)),
		cdw10:    cdw10,
		cdw11:    cdw11,
		cdw12:    cdw12,
		cdw13:    cdw13,
		cdw14:    cdw14,
	}

	return ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_ADMIN_CMD, uintptr(unsafe.Pointer(&cmd)))
}

func (d *NVMeDevice) readLogPage(logID uint8, buf []byte) error {
	bufLen := len(buf)

	if (bufLen < 4) || (bufLen > 0x4000) || (bufLen%4 != 0) {
		return fmt.Errorf("invalid buffer size")
	}

	cmd := nvmePassthruCommand{
		opcode:   NVME_ADMIN_GET_LOG_PAGE,
		nsid:     0xffffffff, // FIXME
		addr:     uint64(uintptr(unsafe.Pointer(&(buf)[0]))),
		data_len: uint32(bufLen),
		cdw10:    uint32(logID) | (((uint32(bufLen) / 4) - 1) << 16),
	}

	return ioctl.Ioctl(uintptr(d.fd), NVME_IOCTL_ADMIN_CMD, uintptr(unsafe.Pointer(&cmd)))
}
