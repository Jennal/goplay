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
	ENCODING_PROTOBUF EncodingType = iota
	ENCODING_GOB
	ENCODING_JSON
	ENCODING_BSON
)

type PackageID byte

type Package struct {
	Type     PackageType
	Encoding EncodingType
	ID       PackageID
	Content  interface{}
}
