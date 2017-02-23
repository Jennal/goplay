package pkg

type PackageType byte

const (
	PKG_NOTIFY PackageType = iota
	PKG_NOTIFY_RESPONSE
	PKG_REQUEST
	PKG_RESPONSE
)

type EncodingType byte

const (
	ENCODING_GOB EncodingType = iota
	ENCODING_JSON
	ENCODING_BSON
	ENCODING_PROTOBUF
)

type PackageID byte

type Header struct {
	Type        PackageType
	Encoding    EncodingType
	ID          PackageID
	ContentSize uint32
}
