package pkg

const maxID PackageIDType = 255

var nextID PackageIDType = 0

func NextID() PackageIDType {
	if nextID == maxID {
		defer func() { nextID = 0 }()
	} else {
		defer func() { nextID++ }()
	}

	return nextID
}
