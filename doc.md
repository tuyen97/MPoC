### Struct

Mỗi block được tạo ra sau blockIntervalMs, 

```go
var (
	// blockIntervalMs is the block genration interval in milli-seconds.
	blockIntervalMs int64
	// bpMinTimeLimitMs is the minimum block generation time limit in milli-sconds.
	bpMinTimeLimitMs int64
	// bpMaxTimeLimitMs is the maximum block generation time limit in milli-seconds.
	bpMaxTimeLimitMs int64
)

// Slot is a DPoS slot implmentation.
type Slot struct {
	timeNs    int64 // nanosecond
	timeMs    int64 // millisecond
	prevIndex int64
	nextIndex int64
}

```


