package redis_state

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	osqueue "github.com/inngest/inngest/pkg/execution/queue"
	"github.com/inngest/inngest/pkg/execution/state"
	"github.com/oklog/ulid/v2"
)

type legacyKeyGenerator interface {
	// Workflow returns the key for the current workflow ID and version.
	Workflow(ctx context.Context, workflowID uuid.UUID, version int) string

	// Idempotency stores the idempotency key for atomic lookup.
	Idempotency(context.Context, state.Identifier) string

	// RunMetadata stores state regarding the current run identifier, such
	// as the workflow version, the time the run started, etc.
	RunMetadata(ctx context.Context, runID ulid.ULID) string

	// Event returns the key used to store the specific event for the
	// given workflow run.
	Event(context.Context, state.Identifier) string

	// Events returns the key used to store the specific batch for the
	// given workflow run.
	Events(context.Context, state.Identifier) string

	// Actions returns the key used to store the action response map used
	// for given workflow run - ie. the results for individual steps.
	Actions(context.Context, state.Identifier) string

	// Errors returns the key used to store the error hash map used
	// for given workflow run.
	Errors(context.Context, state.Identifier) string

	// PauseLease stores the key which references a pause's lease.
	//
	// This is stored independently as we may store more than one copy of a pause
	// for easy iteration.
	PauseLease(context.Context, uuid.UUID) string

	// PauseID returns the key used to store an individual pause from its ID.
	PauseID(context.Context, uuid.UUID) string

	// PauseEvent returns the key used to store data for loading pauses by events.
	PauseEvent(context.Context, uuid.UUID, string) string

	// PauseIndex is a key that's used to index added/expired times for pauses.
	//
	// Added times are necessary to load pauses after a specific point in time,
	// which is used when caching pauses in-memory to only load the subset of pauses
	// added after the cache was last updated.
	PauseIndex(ctx context.Context, kind string, wsID uuid.UUID, event string) string

	// RunPauses stores pause IDs for each run as a zset
	RunPauses(ctx context.Context, runID ulid.ULID) string

	// Invoke returns the key used to store the correlation key associated with invoke functions
	Invoke(ctx context.Context, wsID uuid.UUID) string

	// History returns the key used to store a log entry for run hisotry
	History(ctx context.Context, runID ulid.ULID) string

	// Stack returns the key used to store the stack for a given run
	Stack(ctx context.Context, runID ulid.ULID) string
}

type legacyDefaultKeyFunc struct {
	Prefix string
}

func (d legacyDefaultKeyFunc) Idempotency(ctx context.Context, id state.Identifier) string {
	return fmt.Sprintf("%s:key:%s", d.Prefix, id.IdempotencyKey())
}

func (d legacyDefaultKeyFunc) RunMetadata(ctx context.Context, runID ulid.ULID) string {
	return fmt.Sprintf("%s:metadata:%s", d.Prefix, runID)
}

func (d legacyDefaultKeyFunc) Workflow(ctx context.Context, id uuid.UUID, version int) string {
	return fmt.Sprintf("%s:workflows:%s-%d", d.Prefix, id, version)
}

func (d legacyDefaultKeyFunc) Event(ctx context.Context, id state.Identifier) string {
	return fmt.Sprintf("%s:events:%s:%s", d.Prefix, id.WorkflowID, id.RunID)
}

func (d legacyDefaultKeyFunc) Events(ctx context.Context, id state.Identifier) string {
	return fmt.Sprintf("%s:bulk-events:%s:%s", d.Prefix, id.WorkflowID, id.RunID)
}

func (d legacyDefaultKeyFunc) Actions(ctx context.Context, id state.Identifier) string {
	return fmt.Sprintf("%s:actions:%s:%s", d.Prefix, id.WorkflowID, id.RunID)
}

func (d legacyDefaultKeyFunc) Errors(ctx context.Context, id state.Identifier) string {
	return fmt.Sprintf("%s:errors:%s:%s", d.Prefix, id.WorkflowID, id.RunID)
}

func (d legacyDefaultKeyFunc) PauseID(ctx context.Context, id uuid.UUID) string {
	return fmt.Sprintf("%s:pauses:%s", d.Prefix, id.String())
}

func (d legacyDefaultKeyFunc) PauseLease(ctx context.Context, id uuid.UUID) string {
	return fmt.Sprintf("%s:pause-lease:%s", d.Prefix, id.String())
}

func (d legacyDefaultKeyFunc) PauseEvent(ctx context.Context, workspaceID uuid.UUID, event string) string {
	return fmt.Sprintf("%s:pause-events:%s:%s", d.Prefix, workspaceID, event)
}

func (d legacyDefaultKeyFunc) PauseIndex(ctx context.Context, kind string, wsID uuid.UUID, event string) string {
	if event == "" {
		return fmt.Sprintf("%s:pause-idx:%s:%s:-", d.Prefix, kind, wsID)
	}
	return fmt.Sprintf("%s:pause-idx:%s:%s:%s", d.Prefix, kind, wsID, event)
}

func (d legacyDefaultKeyFunc) RunPauses(ctx context.Context, runID ulid.ULID) string {
	return fmt.Sprintf("%s:pr:%s", d.Prefix, runID)
}

func (d legacyDefaultKeyFunc) Invoke(ctx context.Context, wsID uuid.UUID) string {
	return fmt.Sprintf("%s:invoke:%s", d.Prefix, wsID)
}

func (d legacyDefaultKeyFunc) History(ctx context.Context, runID ulid.ULID) string {
	return fmt.Sprintf("%s:history:%s", d.Prefix, runID)
}

func (d legacyDefaultKeyFunc) Stack(ctx context.Context, runID ulid.ULID) string {
	return fmt.Sprintf("%s:stack:%s", d.Prefix, runID)
}

type legacyQueueKeyGenerator interface {
	// QueueItem returns the key for the hash containing all items within a
	// queue for a function.
	QueueItem() string
	// QueueIndex returns the key containing the sorted zset for a function
	// queue.
	QueueIndex(id string) string

	//
	// Partition keys
	//

	// Shards is a key to a hashmap of shards available.  The values of this
	// key are JSON-encoded shards.
	Shards() string
	// PartitionItem returns the key for the hash containing all partition items.
	PartitionItem() string
	// PartitionMeta returns the key to store metadata for partitions, eg.
	// the number of items enqueued, number in progress, etc.
	PartitionMeta(id string) string
	// GlobalPartitionIndex returns the sorted set for the partition queue;  the
	// earliest time that each function is available.  This is a global queue of
	// all functions across every partition, used for minimum latency.
	GlobalPartitionIndex() string
	// ShardPartitionIndex returns the sorted set for the shard's partition queue.
	ShardPartitionIndex(shard string) string
	// ThrottleKey returns the throttle key for a given queue item.
	ThrottleKey(t *osqueue.Throttle) string

	//
	// Queue metadata keys
	//

	// Sequential returns the key which allows a worker to claim sequential processing
	// of the partitions.
	Sequential() string
	// Scavenger returns the key which allows a worker to claim scavenger processing
	// of the partitions for lost jobs
	Scavenger() string
	// Idempotency stores the map for storing idempotency keys in redis
	Idempotency(key string) string
	// Concurrency returns a key for a given concurrency string.  This stores an ordered
	// zset of items that are in progress for the given concurrency key, giving us a total count
	// of in-progress leased items.
	Concurrency(prefix, key string) string
	// ConcurrencyIndex returns a key for storing pointers to partition concurrency queues that
	// have in-progress work.  This allows us to scan and scavenge jobs in concurrency queues where
	// leases have expired (in the case of failed workers)
	ConcurrencyIndex() string

	// RunIndex returns the index for storing job IDs associated with run IDs.
	RunIndex(runID ulid.ULID) string

	// Status returns the key used for status queue for the provided function.
	Status(status string, fnID uuid.UUID) string

	// ConcurrencyFnEWMA returns the key storing the amount of times of concurrency hits, used for
	// calculating the EWMA value for the function
	ConcurrencyFnEWMA(fnID uuid.UUID) string

	// ***************** Deprecated ************************
	BatchPointer(context.Context, uuid.UUID) string
	Batch(context.Context, ulid.ULID) string
	BatchMetadata(context.Context, ulid.ULID) string
}

type legacyDebounceKeyGenerator interface {
	// QueueItem returns the key for the hash containing all items within a
	// queue for a function.  This is used to check leases on debounce jobs.
	QueueItem() string
	// DebouncePointer returns the key which stores the pointer to the current debounce
	// for a given function.
	DebouncePointer(ctx context.Context, fnID uuid.UUID, key string) string
	// Debounce returns the key for storing debounce-related data given a debounce ID.
	Debounce(ctx context.Context) string
}

type legacyDefaultQueueKeyGenerator struct {
	Prefix string
}

func (d legacyDefaultQueueKeyGenerator) Shards() string {
	return fmt.Sprintf("%s:queue:shards", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) QueueItem() string {
	return fmt.Sprintf("%s:queue:item", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) QueueIndex(id string) string {
	return fmt.Sprintf("%s:queue:sorted:%s", d.Prefix, id)
}

func (d legacyDefaultQueueKeyGenerator) PartitionItem() string {
	return fmt.Sprintf("%s:partition:item", d.Prefix)
}

// GlobalPartitionIndex returns the sorted index for the partition group, which stores the earliest
// time for each function/queue in the partition.
//
// This is grouped so that we can make N partitions independently.
func (d legacyDefaultQueueKeyGenerator) GlobalPartitionIndex() string {
	return fmt.Sprintf("%s:partition:sorted", d.Prefix)
}

// GlobalPartitionIndex returns the sorted index for the partition group, which stores the earliest
// time for each function/queue in the partition.
//
// This is grouped so that we can make N partitions independently.
func (d legacyDefaultQueueKeyGenerator) ShardPartitionIndex(shard string) string {
	if shard == "" {
		return fmt.Sprintf("%s:shard:-", d.Prefix)
	}
	return fmt.Sprintf("%s:shard:%s", d.Prefix, shard)
}

func (d legacyDefaultQueueKeyGenerator) ThrottleKey(t *osqueue.Throttle) string {
	if t == nil || t.Key == "" {
		return fmt.Sprintf("%s:throttle:-", d.Prefix)
	}
	return fmt.Sprintf("%s:throttle:%s", d.Prefix, t.Key)
}

func (d legacyDefaultQueueKeyGenerator) PartitionMeta(id string) string {
	return fmt.Sprintf("%s:partition:meta:%s", d.Prefix, id)
}

func (d legacyDefaultQueueKeyGenerator) Sequential() string {
	return fmt.Sprintf("%s:queue:sequential", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) Scavenger() string {
	return fmt.Sprintf("%s:queue:scavenger", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) Idempotency(key string) string {
	return fmt.Sprintf("%s:queue:seen:%s", d.Prefix, key)
}

func (d legacyDefaultQueueKeyGenerator) Concurrency(prefix, key string) string {
	if key == "" {
		// None supplied; this means ignore.
		return fmt.Sprintf("%s:-", d.Prefix)
	}
	return fmt.Sprintf("%s:concurrency:%s:%s", d.Prefix, prefix, key)
}

func (d legacyDefaultQueueKeyGenerator) ConcurrencyIndex() string {
	return fmt.Sprintf("%s:concurrency:sorted", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) QueuePrefix() string {
	return d.Prefix
}

func (d legacyDefaultQueueKeyGenerator) BatchPointer(ctx context.Context, workflowID uuid.UUID) string {
	return fmt.Sprintf("%s:workflows:%s:batch", d.Prefix, workflowID)
}

func (d legacyDefaultQueueKeyGenerator) BatchPointerWithKey(ctx context.Context, workflowID uuid.UUID, batchKey string) string {
	return fmt.Sprintf("%s:%s", d.BatchPointer(ctx, workflowID), batchKey)
}

func (d legacyDefaultQueueKeyGenerator) Batch(ctx context.Context, batchID ulid.ULID) string {
	return fmt.Sprintf("%s:batches:%s", d.Prefix, batchID)
}

func (d legacyDefaultQueueKeyGenerator) BatchMetadata(ctx context.Context, batchID ulid.ULID) string {
	return fmt.Sprintf("%s:metadata", d.Batch(ctx, batchID))
}

// DebouncePointer returns the key which stores the pointer to the current debounce
// for a given function.
func (d legacyDefaultQueueKeyGenerator) DebouncePointer(ctx context.Context, fnID uuid.UUID, key string) string {
	return fmt.Sprintf("%s:debounce-ptrs:%s:%s", d.Prefix, fnID, key)
}

// Debounce returns the key for storing debounce-related data given a debounce ID.
// This is a hash of debounce IDs -> debounces.
func (d legacyDefaultQueueKeyGenerator) Debounce(ctx context.Context) string {
	return fmt.Sprintf("%s:debounce-hash", d.Prefix)
}

func (d legacyDefaultQueueKeyGenerator) RunIndex(runID ulid.ULID) string {
	return fmt.Sprintf("%s:idx:run:%s", d.Prefix, runID)
}

func (d legacyDefaultQueueKeyGenerator) Status(status string, fnID uuid.UUID) string {
	return fmt.Sprintf("%s:queue:status:%s:%s", d.Prefix, fnID, status)
}

func (d legacyDefaultQueueKeyGenerator) ConcurrencyFnEWMA(fnID uuid.UUID) string {
	return fmt.Sprintf("%s:queue:concurrency-ewma:%s", d.Prefix, fnID)
}
