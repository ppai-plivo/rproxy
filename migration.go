package main

type phaseType int

const (
	WriteBothReadSrc phaseType = iota
	WriteBothReadDst
	WriteDstReadDst
)

func (p phaseType) String() string {

	switch p {
	case WriteBothReadSrc:
		return "Write to both redis; Read from src redis"
	case WriteBothReadDst:
		return "Write to both redis; Read from dst redis"
	case WriteDstReadDst:
		return "Write to and read from dst redis only"
	}

	return "Invalid migration phase"
}
