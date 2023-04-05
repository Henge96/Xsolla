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

func DomTaskKind(txt string) TaskKind {
	switch {
	case TaskKindEventAdd.String() == txt:
		return TaskKindEventAdd
	case TaskKindEventUpdate.String() == txt:
		return TaskKindEventUpdate
	default:
		return 0
	}
}
