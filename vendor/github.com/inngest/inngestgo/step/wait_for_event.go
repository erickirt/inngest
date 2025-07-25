package step

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/inngest/inngest/pkg/enums"
	"github.com/inngest/inngestgo/internal/sdkrequest"
	str2duration "github.com/xhit/go-str2duration/v2"
)

// ErrEventNotReceived is returned when a WaitForEvent call times out.  It indicates that a
// matching event was not received before the timeout.
var ErrEventNotReceived = fmt.Errorf("event not received")

type WaitForEventOpts struct {
	// Name represents the optional step name.
	Name string
	// Event is the event name to wait for.
	Event string
	// Timeout is how long to wait.  We must always timebound event lsiteners.
	Timeout time.Duration
	// If allows you to write arbitrary expressions to match against.
	If *string `json:"if"`
}

// WaitForEvent pauses function execution until a specific event is received or the wait times
// out.  You must pass in an event name within WaitForEventOpts.Event, and may pass an optional
// expression to filter events based off of data.
//
// For example:
//
//	step.waitForEvent(ctx, "wait-for-open", opts.WaitForEventOpts{
//		Event: "email/mail.opened",
//		If:	inngestgo.StrPtr(fmt.Sprintf("async.data.id == %s", strconv.Quote("my-id"))),
//		Timeout: 24 * time.Hour,
//	})
func WaitForEvent[T any](ctx context.Context, stepID string, opts WaitForEventOpts) (T, error) {
	targetID := getTargetStepID(ctx)
	mgr := preflight(ctx)
	args := map[string]any{
		"timeout": str2duration.String(opts.Timeout),
		"event":   opts.Event,
	}
	if opts.If != nil {
		args["if"] = *opts.If
	}
	if opts.Name == "" {
		opts.Name = stepID
	}
	op := mgr.NewOp(enums.OpcodeWaitForEvent, stepID, args)
	hashedID := op.MustHash()

	if val, ok := mgr.Step(ctx, op); ok {
		var output T
		if val == nil || bytes.Equal(val, []byte{0x6e, 0x75, 0x6c, 0x6c}) {
			return output, ErrEventNotReceived
		}
		if err := json.Unmarshal(val, &output); err != nil {
			mgr.SetErr(fmt.Errorf("error unmarshalling wait for event value in '%s': %w", opts.Event, err))
			panic(ControlHijack{})
		}
		return output, nil
	}

	if targetID != nil && *targetID != hashedID {
		// Don't report this step since targeting is happening and it isn't
		// targeted
		panic(ControlHijack{})
	}

	plannedOp := sdkrequest.GeneratorOpcode{
		ID:   hashedID,
		Op:   op.Op,
		Name: opts.Name,
		Opts: op.Opts,
	}
	plannedOp.SetParallelMode(parallelMode(ctx))
	mgr.AppendOp(plannedOp)
	panic(ControlHijack{})
}
