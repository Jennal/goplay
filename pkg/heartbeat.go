package pkg

func NewHeartBeatHeader() *Header {
	return NewHeader(PKG_HEARTBEAT, ENCODING_NONE, "")
}

func NewHeartBeatResponseHeader(h *Header) *Header {
	newHB := &Header{}
	*newHB = *h
	newHB.Type = PKG_HEARTBEAT_RESPONSE
	return newHB
}
