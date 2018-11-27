package snake

import "Wave/application/snake/core"

// ----------------| wall node

type wallNode struct {
	*core.Object // base object
}

func newWallNode(*wall) *wallNode {
	return nil
}

// ----------------| wall

type wall struct {
}
