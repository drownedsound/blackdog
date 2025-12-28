package infra

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

type IdGenerator interface {
	Generate() int64
}

// SnowflakeGenerator implements IdGenerator using Twitter Snowflake algorithm.
type SnowflakeIdGenerator struct {
	node *snowflake.Node
	// mu ensures safety if we ever need to reconfigure the node, 
	// though the library's Generate() is already thread-safe.
	mu   sync.Mutex
}

// NewSnowflakeIdGenerator initializes a new node.
// nodeId should be unique per running instance of the application (0-1023).
func NewSnowflakeIdGenerator(nodeId int64) (*SnowflakeIdGenerator, error) {
	// FIXME: Get nodeId from configuration file
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize snowflake node: %w", err)
	}
	return &SnowflakeIdGenerator{node: node}, nil
}

// Generate creates a new unique int64 Id.
// This operation is thread-safe and highly performant.
func (g *SnowflakeIdGenerator) Generate() int64 {
	return g.node.Generate().Int64()
}

