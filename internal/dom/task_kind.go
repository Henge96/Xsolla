package dom

type (
	// TaskKind represents kind of task.
	TaskKind uint8
)

//go:generate stringer -output=stringer.TaskKind.go -type=TaskKind -trimprefix=TaskKind
const (
	_ TaskKind = iota
	TaskKindEventAdd
	TaskKindEventUpdate
)
