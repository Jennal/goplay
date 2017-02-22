package protocol

import "testing"
import "github.com/jennal/goplay/handler/pkg"
import "fmt"

func TestJsonDecode(t *testing.T) {
	encoder := JsonEncoder{}
	decoder := JsonDecoder{}

	content := []int{1, 2, 3}
	pack := pkg.Header{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_JSON,
		ID:       2,
	}
	buffer := encoder.Marshal(&pack, content)

	newPack := pkg.Header{}
	var newContent []int
	decoder.Unmarshal(buffer, &newPack, &newContent)

	fmt.Println(pack, newPack)
	fmt.Println(content, newContent)

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

	if content[0] != newContent[0] {
		t.Errorf("package.Content[0] are not equal %v != %v",
			content[0],
			newContent[0])
		t.Fail()
	}
}

func BenchmarkJsonDecode(b *testing.B) {
	encoder := JsonEncoder{}
	decoder := JsonDecoder{}
	pack := pkg.Header{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_JSON,
		ID:       2,
	}
	content := []int{1, 2, 3, 4}

	buffer := encoder.Marshal(&pack, content)
	newHeader := pkg.Header{}
	var newContent []int
	for i := 0; i < b.N; i++ {
		decoder.Unmarshal(buffer, &newHeader, &newContent)
	}
}
