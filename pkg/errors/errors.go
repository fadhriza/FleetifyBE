package errors

import (
	"fmt"
	"log"
	"runtime/debug"
)

func LogError(context string, err error) {
	log.Printf("\n❌ ERROR [%s]:\n", context)
	log.Printf("   Error: %v\n", err)
	
	stack := debug.Stack()
	log.Printf("   Stack Trace:\n%s\n", string(stack))
	
	// wrapped errors
	if err != nil {
		log.Printf("   Full Error Chain:\n")
		currentErr := err
		depth := 0
		for currentErr != nil && depth < 10 {
			log.Printf("      [%d] %v\n", depth, currentErr)
			if unwrapped, ok := currentErr.(interface{ Unwrap() error }); ok {
				currentErr = unwrapped.Unwrap()
			} else {
				break
			}
			depth++
		}
	}
	
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// LogErrorf logs
func LogErrorf(context string, format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	LogError(context, err)
}

