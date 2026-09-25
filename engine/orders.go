package engine

import "gocnc/interfaces"

var (
	NewMoveOrder   = interfaces.NewMoveOrder
	NewAttackOrder = interfaces.NewAttackOrder
	NewPatrolOrder = interfaces.NewPatrolOrder
	NewGuardOrder  = interfaces.NewGuardOrder
)
