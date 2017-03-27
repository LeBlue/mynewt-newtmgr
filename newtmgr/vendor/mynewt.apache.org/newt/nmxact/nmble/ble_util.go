package nmble

import (
	"fmt"
	"strconv"
	"sync/atomic"

	log "github.com/Sirupsen/logrus"

	. "mynewt.apache.org/newt/nmxact/bledefs"
	"mynewt.apache.org/newt/nmxact/nmxutil"
)

const NmpPlainSvcUuid = "8D53DC1D-1DB7-4CD3-868B-8A527460AA84"
const NmpPlainChrUuid = "DA2E7828-FBCE-4E01-AE9E-261174997C48"
const NmpOicSvcUuid = "ADE3D529-C784-4F63-A987-EB69F70EE816"
const NmpOicReqChrUuid = "AD7B334F-4637-4B86-90B6-9D787F03D218"
const NmpOicRspChrUuid = "E9241982-4580-42C4-8831-95048216B256"

const WRITE_CMD_BASE_SZ = 3
const NOTIFY_CMD_BASE_SZ = 3

var nextSeq uint32

func NextSeq() int {
	return int(atomic.AddUint32(&nextSeq, 1))
}

func ParseUuid(uuidStr string) (BleUuid, error) {
	bu := BleUuid{}

	if len(uuidStr) != 36 {
		return bu, fmt.Errorf("Invalid UUID: %s", uuidStr)
	}

	boff := 0
	for i := 0; i < 36; {
		switch i {
		case 8, 13, 18, 23:
			if uuidStr[i] != '-' {
				return bu, fmt.Errorf("Invalid UUID: %s", uuidStr)
			}
			i++

		default:
			u64, err := strconv.ParseUint(uuidStr[i:i+2], 16, 8)
			if err != nil {
				return bu, fmt.Errorf("Invalid UUID: %s", uuidStr)
			}
			bu.Bytes[boff] = byte(u64)
			i += 2
			boff++
		}
	}

	return bu, nil
}

func BhdTimeoutError(rspType MsgType) error {
	str := fmt.Sprintf("Timeout waiting for blehostd to send %s response",
		MsgTypeToString(rspType))

	log.Debug(str)
	return nmxutil.NewXportTimeoutError(str)
}

func StatusError(op MsgOp, msgType MsgType, status int) error {
	str := fmt.Sprintf("%s %s indicates error: %s (%d)",
		MsgOpToString(op),
		MsgTypeToString(msgType),
		ErrCodeToString(status),
		status)

	log.Debug(str)
	return nmxutil.NewBleHostError(status, str)
}

func NewBleConnectReq() *BleConnectReq {
	return &BleConnectReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_CONNECT,
		Seq:  NextSeq(),

		OwnAddrType:  BLE_ADDR_TYPE_PUBLIC,
		PeerAddrType: BLE_ADDR_TYPE_PUBLIC,
		PeerAddr:     BleAddr{},

		DurationMs:         30000,
		ScanItvl:           0x0010,
		ScanWindow:         0x0010,
		ItvlMin:            24,
		ItvlMax:            40,
		Latency:            0,
		SupervisionTimeout: 0x0200,
		MinCeLen:           0x0010,
		MaxCeLen:           0x0300,
	}
}

func NewBleTerminateReq() *BleTerminateReq {
	return &BleTerminateReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_TERMINATE,
		Seq:  NextSeq(),

		ConnHandle: 0,
		HciReason:  0,
	}
}

func NewBleConnCancelReq() *BleConnCancelReq {
	return &BleConnCancelReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_CONN_CANCEL,
		Seq:  NextSeq(),
	}
}

func NewBleDiscSvcUuidReq() *BleDiscSvcUuidReq {
	return &BleDiscSvcUuidReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_DISC_SVC_UUID,
		Seq:  NextSeq(),

		ConnHandle: 0,
		Uuid:       BleUuid{},
	}
}

func NewBleDiscAllChrsReq() *BleDiscAllChrsReq {
	return &BleDiscAllChrsReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_DISC_ALL_CHRS,
		Seq:  NextSeq(),
	}
}

func NewBleExchangeMtuReq() *BleExchangeMtuReq {
	return &BleExchangeMtuReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_EXCHANGE_MTU,
		Seq:  NextSeq(),

		ConnHandle: 0,
	}
}

func NewBleWriteCmdReq() *BleWriteCmdReq {
	return &BleWriteCmdReq{
		Op:   MSG_OP_REQ,
		Type: MSG_TYPE_WRITE_CMD,
		Seq:  NextSeq(),

		ConnHandle: 0,
		AttrHandle: 0,
		Data:       BleBytes{},
	}
}