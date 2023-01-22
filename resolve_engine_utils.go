package resolve

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
)

func LogMessage(msg string) {
	yellow := color.New(color.FgHiBlue).SprintFunc()
	white := color.New(color.FgHiWhite).SprintFunc()
	fmt.Printf("%s %s\n", yellow("[Resolve]"), white(msg))
}

func ErrorMessage(msg string, tracer int) {
	red := color.New(color.FgHiRed).SprintFunc()
	white := color.New(color.FgHiWhite).SprintFunc()
	fmt.Printf("%s %s\n", red("[Resolve Error]"), white(msg))

	if tracer != -1 {
		pc := make([]uintptr, 15)
		n := runtime.Callers(tracer, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()

		fmt.Printf("%s %s:%s\n", red("[Resolve Error]"), white(frame.File), white(frame.Line))
	}
}
