package writePolicy

type WritePolicy string

func (r WritePolicy) Value() string {
	return string(r)
}

const (
	WriteAround  WritePolicy = "write-around"
	WriteThrough WritePolicy = "write-through"
	WriteBack    WritePolicy = "write-back"
)
