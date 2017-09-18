package pkg

import (
	"github.com/jennal/goplay/log"
)

type HandShake struct {
	serverResponse *HandShakeResponse
	routesMap      RouteMap
}

var HandShakeInstance *HandShake

func init() {
	HandShakeInstance = &HandShake{
		serverResponse: nil,
		routesMap:      nil,
	}
}

func (r *HandShake) UpdateRoutesMap(data RouteMap) {
	r.routesMap = data
}

func (r *HandShake) UpdateHandShakeResponse(resp *HandShakeResponse) {
	r.serverResponse = resp
	r.routesMap = make(RouteMap)

	for k, v := range resp.Routes {
		r.routesMap[k] = RouteIndex(v)
	}
}

func (r *HandShake) IsInited() bool {
	return r.routesMap != nil
}

func (r *HandShake) RoutesMap() RouteMap {
	return r.routesMap
}

func (r *HandShake) GetIndexRoute(str string) (RouteIndex, bool) {
	if !r.IsInited() {
		log.Errorf("pkg.HandShakeInstance.routesMap not inited")
		return ROUTE_INDEX_NONE, false
	}

	val, ok := r.routesMap[str]
	if !ok {
		return ROUTE_INDEX_NONE, false
	}

	return val, ok
}

func (r *HandShake) GetStringRoute(idx RouteIndex) (string, bool) {
	if !r.IsInited() {
		log.Errorf("pkg.HandShakeInstance.routesMap not inited")
		return "", false
	}

	for path, i := range r.routesMap {
		if i == idx {
			return path, true
		}
	}

	return "", false
}
