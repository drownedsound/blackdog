package infra

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	// app "github.com/drownedsound/blackdog/internal/features/application"
)

type IdGenerator interface {
	Generate() int64
}

// SnowflakeIdGenerator implements IdGenerator using the
// Twitter Snowflake algorithm.
type SnowflakeIdGenerator struct {
	node *snowflake.Node
}

// NewSnowflakeIdGenerator initializes a new node.
// nodeId should be unique per running instance of the application (0-1023).
func NewSnowflakeIdGenerator(nodeId int64) (
	idGen *SnowflakeIdGenerator,
	err error,
) {
	// FIXME: Get nodeId from configuration file
	// if nodeId < 0 || nodeId >= 1023 {
	// 	return nil, app.ErrInvalidNodeId
	// }

	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize snowflake node: %w", err)
	}

	return &SnowflakeIdGenerator{node: node}, nil
}

// Generate creates a new unique int64 Id.
func (g *SnowflakeIdGenerator) Generate() int64 {
	return g.node.Generate().Int64()
}
