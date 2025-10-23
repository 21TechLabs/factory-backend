package dto

import (
	"time"

	"github.com/21TechLabs/factory-backend/utils"
)

type MinMax[T utils.Number | time.Time] struct {
	Min *T
	Max *T
}
