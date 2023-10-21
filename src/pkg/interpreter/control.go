package interpreter

type ControlType string

const (
	RETURN   ControlType = "RETURN"
	BREAK                = "BREAK"
	CONTINUE             = "CONTINUE"
)

type LoxControl struct {
	controlType ControlType
	value       any
}

func LoxReturn(value any) *LoxControl {
	return &LoxControl{controlType: RETURN, value: value}
}

var LoxBreak = &LoxControl{controlType: BREAK}

var LoxContinue = &LoxControl{controlType: CONTINUE}
