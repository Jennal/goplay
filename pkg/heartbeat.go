package pkg

func NewHeartBeatHeader(idGen *IDGen) *Header {
	return NewHeader(PKG_HEARTBEAT, ENCODING_NONE, idGen, "")
}

func NewHeartBeatResponseHeader(h *Header) *Header {
	newHB := &Header{}
	*newHB = *h
	newHB.Type = PKG_HEARTBEAT_RESPONSE
	return newHB
}
