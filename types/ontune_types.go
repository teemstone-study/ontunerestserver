package types

const (
	DATAKEY_CODE  = 0x00000001
	HOST_CODE     = 0x00000002
	LASTPERF_CODE = 0x00000004
	BASIC_CODE    = 0x00000008
	CPU_CODE      = 0x00000010
	MEM_CODE      = 0x00000020
	NET_CODE      = 0x00000040
	DISK_CODE     = 0x00000080
)

const (
	HOST_KEY     = 2
	LASTPERF_KEY = 4
	BASIC_KEY    = 8
	CPU_KEY      = 16
	MEM_KEY      = 32
	NET_KEY      = 64
	DISK_KEY     = 128
	//ALL_KEY = 254
)

type DataCode struct {
	Code uint32 `json:"code"`
}

type DataKey struct {
	Code uint32  `json:"code"`
	Key  Bitmask `json:"key"`
}

type RealData struct {
	Code uint32 `json:"code"`
	Data string `json:"data"`
}
