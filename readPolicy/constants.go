package readPolicy

type ReadPolicy string

func (r ReadPolicy) Value() string {
	return string(r)
}

const (
	ReadThrough ReadPolicy = "read-through"
	CacheAside  ReadPolicy = "cache-aside"
)
