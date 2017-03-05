package pkg

type PackageType byte

const (
	PKG_NOTIFY PackageType = iota
	PKG_NOTIFY_RESPONSE
	PKG_REQUEST
	PKG_RESPONSE
	PKG_HEARTBEAT
	PKG_HEARTBEAT_RESPONSE
)

type EncodingType byte

const (
	ENCODING_GOB EncodingType = iota
	ENCODING_JSON
	ENCODING_BSON
	ENCODING_PROTOBUF
)

type PackageIDType byte
type PackageSizeType uint16

type Header struct {
	Type        PackageType
	Encoding    EncodingType
	ID          PackageIDType
	ContentSize PackageSizeType
	Route       string
}

func NewHeader(t PackageType, e EncodingType, r string) *Header {
	return &Header{
		Type:        t,
		Encoding:    e,
		ID:          NextID(),
		ContentSize: 0,
		Route:       r,
	}
}
