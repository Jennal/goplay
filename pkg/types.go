package pkg

type Status int16

const (
	STAT_OK Status = iota
	STAT_ERR_TIMEOUT
)

type PackageType byte

const (
	PKG_NOTIFY PackageType = iota
	PKG_REQUEST
	PKG_RESPONSE
	PKG_HEARTBEAT
	PKG_HEARTBEAT_RESPONSE

	PKG_RPC_NOTIFY
	PKG_RPC_REQUEST
	PKG_RPC_RESPONSE
)

type EncodingType byte

const (
	ENCODING_NONE EncodingType = iota
	ENCODING_GOB
	ENCODING_JSON
	ENCODING_BSON
	ENCODING_PROTOBUF
)

type PackageIDType byte
type PackageSizeType uint16
