package readPolicy

type ReadPolicy string

func (r ReadPolicy) Value() string {
	return string(r)
}

const (
	CacheAside ReadPolicy = "cache-aside"
)
