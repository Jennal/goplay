package pkg

func NewRouteMap() RouteMap {
	return make(RouteMap)
}

func (rm RouteMap) Add(route string, encoded RouteIndex) {
	rm[route] = encoded
}

func (rm RouteMap) Merge(data RouteMap) {
	for k, v := range data {
		rm[k] = v
	}
}

func (rm RouteMap) IsEmpty() bool {
	if rm == nil {
		return true
	}

	return len(rm) == 0
}

func (r RouteMap) GetIndexRoute(str string) (RouteIndex, bool) {
	if r.IsEmpty() {
		return ROUTE_INDEX_NONE, false
	}

	val, ok := r[str]
	if !ok {
		return ROUTE_INDEX_NONE, false
	}

	return val, ok
}

func (r RouteMap) GetStringRoute(idx RouteIndex) (string, bool) {
	if r.IsEmpty() {
		return "", false
	}

	for path, i := range r {
		if i == idx {
			return path, true
		}
	}

	return "", false
}
