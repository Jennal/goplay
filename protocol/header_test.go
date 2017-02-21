package protocol

import "testing"
import "github.com/jennal/goplay/handler/pkg"

func TestHeaderDecode(t *testing.T) {
	encoder := HeaderEncoder{}
	decoder := HeaderDecoder{}
	pack := pkg.Header{
		Type:        pkg.PKG_NOTIFY,
		Encoding:    pkg.ENCODING_GOB,
		ID:          2,
		ContentSize: 10,
	}
	buffer := encoder.MarshalHeader(&pack)
	newPack := pkg.Header{}
	decoder.UnmarshalHeader(buffer, &newPack)

	if pack.Type != newPack.Type {
		t.Errorf("package.Type are not equal %v != %v", pack.Type, newPack.Type)
		t.Fail()
	}

	if pack.Encoding != newPack.Encoding {
		t.Errorf("package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
		t.Fail()
	}

	if pack.ID != newPack.ID {
		t.Errorf("package.ID are not equal %v != %v", pack.ID, newPack.ID)
		t.Fail()
	}

	if pack.ContentSize != newPack.ContentSize {
		t.Errorf("package.ContentSize are not equal %v != %v", pack.ContentSize, newPack.ContentSize)
		t.Fail()
	}
}

func BenchmarkHeaderDecode(b *testing.B) {
	encoder := HeaderEncoder{}
	decoder := HeaderDecoder{}
	pack := pkg.Header{
		Type:        pkg.PKG_NOTIFY,
		Encoding:    pkg.ENCODING_GOB,
		ID:          2,
		ContentSize: 3,
	}
	buffer := encoder.MarshalHeader(&pack)
	newHeader := pkg.Header{}
	for i := 0; i < b.N; i++ {
		decoder.UnmarshalHeader(buffer, &newHeader)
	}
}
