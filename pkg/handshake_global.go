package pkg

import "github.com/jennal/goplay/log"

type HandShakeImpl struct {
	serverResponse *HandShakeResponse
	routesMap      RouteMap
	rpcRoutesMap   RouteMap
	specRoutesMap  RouteMap
}

type HandShake interface {
	RegistSpecRoute(route string, idx RouteIndex) error

	MergeRpcRoutesMap(data RouteMap)

	AddRoutesMap(data string)
	UpdateHandShakeResponse(resp *HandShakeResponse)

	RoutesMap() RouteMap
	RpcRoutesMap() RouteMap
	SpecRoutesMap() RouteMap

	GetIndexRoute(str string) (RouteIndex, bool)
	GetStringRoute(idx RouteIndex) (string, bool)

	ConvertRouteIndexToRpc(route string) (RouteIndex, bool)
	ConvertRouteIndexFromRpc(route string) (RouteIndex, bool)
}

var _handShakeInstance HandShake

func init() {
	_handShakeInstance = NewHandShakeImpl()
}

func SetHandShakeImpl(hs HandShake) error {
	if hs == nil {
		return log.NewError("handshake can't be nil")
	}

	//restore spec
	if _handShakeInstance != nil {
		for k, v := range _handShakeInstance.SpecRoutesMap() {
			hs.RegistSpecRoute(k, v)
		}
	}

	_handShakeInstance = hs
	return nil
}

func DefaultHandShake() HandShake {
	return _handShakeInstance
}

func NewHandShakeImpl() *HandShakeImpl {
	return &HandShakeImpl{
		serverResponse: nil,
		routesMap:      make(RouteMap),
		rpcRoutesMap:   make(RouteMap),
		specRoutesMap:  make(RouteMap),
	}
}

func (r *HandShakeImpl) RegistSpecRoute(route string, idx RouteIndex) error {
	if _, ok := r.specRoutesMap[route]; ok {
		return log.NewErrorf("route already exists: %v", route)
	}

	r.specRoutesMap[route] = idx
	return nil
}

func (r *HandShakeImpl) MergeRpcRoutesMap(data RouteMap) {
	r.rpcRoutesMap.Merge(data)
}

func (r *HandShakeImpl) AddRoutesMap(route string) {
	count := len(r.routesMap)
	if r.routesMap == nil {
		r.routesMap = RouteMap{}
	}

	log.Tracef("Build route map: %v %v", count, route)
	r.routesMap[route] = RouteIndex(count)
}

func (r *HandShakeImpl) UpdateHandShakeResponse(resp *HandShakeResponse) {
	r.serverResponse = resp

	routesMap := make(RouteMap)
	for k, v := range resp.Routes {
		routesMap[k] = RouteIndex(v)
	}

	r.rpcRoutesMap.Merge(routesMap)
}

func (r *HandShakeImpl) RoutesMap() RouteMap {
	return r.routesMap
}

func (r *HandShakeImpl) RpcRoutesMap() RouteMap {
	return r.rpcRoutesMap
}

func (r *HandShakeImpl) SpecRoutesMap() RouteMap {
	return r.specRoutesMap
}

func (r *HandShakeImpl) GetIndexRoute(str string) (RouteIndex, bool) {
	if val, ok := r.specRoutesMap.GetIndexRoute(str); ok {
		return val, ok
	}

	if val, ok := r.routesMap.GetIndexRoute(str); ok {
		return val, ok
	}

	return r.rpcRoutesMap.GetIndexRoute(str)
}

func (r *HandShakeImpl) GetStringRoute(idx RouteIndex) (string, bool) {
	if val, ok := r.specRoutesMap.GetStringRoute(idx); ok {
		return val, ok
	}

	if val, ok := r.routesMap.GetStringRoute(idx); ok {
		return val, ok
	}

	return r.rpcRoutesMap.GetStringRoute(idx)
}

func (r *HandShakeImpl) ConvertRouteIndexToRpc(route string) (RouteIndex, bool) {
	return r.rpcRoutesMap.GetIndexRoute(route)
}

func (r *HandShakeImpl) ConvertRouteIndexFromRpc(route string) (RouteIndex, bool) {
	return r.routesMap.GetIndexRoute(route)
}
