package infra

import (
	// "errors"
	"testing"
	"strings"
	// app "github.com/drownedsound/blackdog/internal/features/application"
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
					t.Errorf("NewSnowflakeIdGenerator Error: %v", err,
						)
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

// func Test_SnowflakeIdGenerator_Concurrency(t *testing.T) {
// 	gen, err := NewSnowflakeIdGenerator(2)
// 	if err != nil {
// 		t.Fatalf("Failed to Create Id Generator: %v", err)
// 	}
//
// 	const numGoroutines = 50
// 	const idsPerGoroutine = 1000
//
// 	var wg sync.WaitGroup
// 	wg.Add(numGoroutines)
//
// 	// Use a sharded data structure or sync.Map to collect IDs concurrently.
// 	// For simplicity in testing, a mutex-protected map is used here to verify uniqueness
// 	// across threads.
// 	generatedIds := make(map[int64]bool)
// 	var mu sync.Mutex
//
// 	for i := 0; i < numGoroutines; i++ {
// 		go func() {
// 			defer wg.Done()
// 			for j := 0; j < idsPerGoroutine; j++ {
// 				id := gen.Generate()
//
// 				mu.Lock()
// 				if generatedIds[id] {
// 					// We use t.Error inside a protected block, but note that
// 					// testing.T is not strictly thread-safe for FailNow calls.
// 					// Error is generally acceptable in simple concurrency tests.
// 					t.Errorf("Duplicate ID found in concurrency test: %d", id)
// 				}
// 				generatedIds[id] = true
// 				mu.Unlock()
// 			}
// 		}()
// 	}
//
// 	wg.Wait()
// }
