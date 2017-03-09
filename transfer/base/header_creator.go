package base

import "github.com/jennal/goplay/pkg"

type IHeaderCreator interface {
	NewHeader(t pkg.PackageType, e pkg.EncodingType, r string) *pkg.Header
	NewHeartBeatHeader() *pkg.Header
	NewHeartBeatResponseHeader(h *pkg.Header) *pkg.Header
}

type HeaderCreator struct {
	idGen *pkg.IDGen
}

func NewHeaderCreator() *HeaderCreator {
	return &HeaderCreator{
		idGen: pkg.NewIDGen(),
	}
}

func (self *HeaderCreator) NewHeader(t pkg.PackageType, e pkg.EncodingType, r string) *pkg.Header {
	return pkg.NewHeader(t, e, self.idGen, r)
}

func (self *HeaderCreator) NewHeartBeatHeader() *pkg.Header {
	return pkg.NewHeartBeatHeader(self.idGen)
}

func (self *HeaderCreator) NewHeartBeatResponseHeader(h *pkg.Header) *pkg.Header {
	return pkg.NewHeartBeatResponseHeader(h)
}
