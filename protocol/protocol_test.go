package protocol

import "testing"
import "github.com/jennal/goplay/handler/pkg"

func TestBaseDecode(t *testing.T) {
	encoder := BaseEncoder{}
	decoder := BaseDecoder{}
	pack := pkg.Package{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_GOB,
		ID:       2,
	}
	buffer := encoder.Marshal(&pack)
	newPack := decoder.Unmarshal(buffer)

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
}

func BenchmarkBaseDecode(b *testing.B) {
	encoder := BaseEncoder{}
	decoder := BaseDecoder{}
	pack := pkg.Package{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_GOB,
		ID:       2,
	}
	buffer := encoder.Marshal(&pack)
	for i := 0; i < b.N; i++ {
		_ = decoder.Unmarshal(buffer)
	}
}
