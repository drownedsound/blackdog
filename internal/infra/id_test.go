package infra

import (
	"strings"
	"testing"
)

func Test_SnowflakeIdGenerator_New(t *testing.T) {
	testCases := []struct {
		desc        string
		nodeId      int64
		expectedErr string
	}{
		{
			"NodeId Within Range Successfully Creates Snowflake Id",
			10,
			"",
		},
		{
			"Negative NodeId Fails to Create Snowflake Id",
			-1,
			"Node number must be between 0 and 1023",
		},
		{
			"Large NodeId Fails to Create Snowflake Id",
			1024,
			"Node number must be between 0 and 1023",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := NewSnowflakeIdGenerator(tC.nodeId)

			if tC.expectedErr == "" {
				if err != nil {
					t.Errorf("NewSnowflakeIdGenerator Error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf(
					"NewSnowflakeIdGenerator Error: Expected = %v, Actual = nil",
					tC.expectedErr,
				)
			}

			if !strings.Contains(err.Error(), tC.expectedErr) {
				t.Errorf(
					"NewSnowflakeIdGenerator Error: Expected = %v, Actual = %v",
					tC.expectedErr, err,
				)
			}
		})
	}
}

func Test_SnowflakeIdGenerator_Generate(t *testing.T) {
	gen, err := NewSnowflakeIdGenerator(1)
	if err != nil {
		t.Fatalf("Failed to Create Generator: %v", err)
	}

	id := gen.Generate()

	if id == 0 {
		t.Error("Generate Error: Expected > 0, Actual = 0")
	}
}
