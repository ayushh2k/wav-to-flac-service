// utils/logger.go
package utils

import (
	"log"
)

func LogError(context string, err error) {
	log.Printf("Error [%s]: %v", context, err)
}
