package pkg

const maxID PackageIDType = 255

type IDGen struct {
	nextID PackageIDType
}

func NewIDGen() *IDGen {
	return &IDGen{
		nextID: 0,
	}
}

func (self *IDGen) NextID() PackageIDType {
	if self.nextID == maxID {
		defer func() { self.nextID = 0 }()
	} else {
		defer func() { self.nextID++ }()
	}

	return self.nextID
}
