package wampproto

import (
	"fmt"
	"math/rand"
	"path"
	"sync"

	"github.com/hashicorp/go-immutable-radix/v2"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/util"
)

const (
	OptionReceiveProgress = "receive_progress"
	OptionProgress        = "progress"
	OptionMatch           = "match"
	OptionInvoke          = "invoke"
	OptionMode            = "mode"
	OptionReason          = "reason"

	CancelModeKill       = "kill"
	CancelModeKillNoWait = "killnowait"
	CancelModeSkip       = "skip"

	MatchExact    = "exact"
	MatchPrefix   = "prefix"
	MatchWildcard = "wildcard"

	InvokeSingle     = "single"
	InvokeFirst      = "first"
	InvokeLast       = "last"
	InvokeRoundRobin = "roundrobin"
	InvokeRandom     = "random"
)

const (
	FeatureProgressiveCallInvocations = "progressive_call_invocations"
	FeatureProgressiveCallResults     = "progressive_call_results"
	FeatureCallCanceling              = "call_canceling"
	FeaturePublisherExclusion         = "publisher_exclusion"
)

type PendingInvocation struct {
	RequestID       uint64
	CallerID        uint64
	CalleeID        uint64
	Progress        bool
	ReceiveProgress bool
	CancelMode      string
}

type Registration struct {
	ID               uint64
	Procedure        string
	Registrants      map[uint64]uint64
	InvocationPolicy string
	nextCallee       int
	callees          []uint64
	Match            string
}

type CallMap struct {
	CallerID uint64
	CallID   uint64
}

type Dealer struct {
	sessions                   map[uint64]*SessionDetails
	registrationsByProcedure   map[string]*Registration
	registrationsBySession     map[uint64]map[uint64]*Registration
	prefixTree                 *iradix.Tree[*Registration]
	wcRegistrationsByProcedure map[string]*Registration
	pendingCalls               map[uint64]*PendingInvocation
	invocationIDbyCall         map[CallMap]uint64
	details                    bool

	idGen *SessionScopeIDGenerator
	sync.Mutex
}

func NewDealer() *Dealer {
	return &Dealer{
		sessions:                   make(map[uint64]*SessionDetails),
		registrationsByProcedure:   make(map[string]*Registration),
		registrationsBySession:     make(map[uint64]map[uint64]*Registration),
		pendingCalls:               make(map[uint64]*PendingInvocation),
		invocationIDbyCall:         make(map[CallMap]uint64),
		idGen:                      &SessionScopeIDGenerator{},
		prefixTree:                 iradix.New[*Registration](),
		wcRegistrationsByProcedure: make(map[string]*Registration),
	}
}

func (d *Dealer) AddSession(details *SessionDetails) error {
	d.Lock()
	defer d.Unlock()

	_, exists := d.sessions[details.ID()]
	if exists {
		return fmt.Errorf("cannot attach an already attached client %d", details.ID())
	}

	d.registrationsBySession[details.ID()] = map[uint64]*Registration{}
	d.sessions[details.ID()] = details
	return nil
}

func (d *Dealer) RemoveSession(id uint64) error {
	d.Lock()
	defer d.Unlock()

	_, exists := d.sessions[id]
	if !exists {
		return fmt.Errorf("cannot remove client with id %d not attached", id)
	}

	registrations := d.registrationsBySession[id]
	for _, reg := range registrations {
		registration, ok := d.registrationsByProcedure[reg.Procedure]
		if !ok {
			continue
		}
		delete(registration.Registrants, id)
		if len(registration.Registrants) == 0 {
			delete(d.registrationsByProcedure, registration.Procedure)
		}
		if registration.Match == MatchPrefix {
			d.prefixTree.Delete([]byte(registration.Procedure))
		}

		if registration.Match == MatchWildcard {
			delete(d.wcRegistrationsByProcedure, registration.Procedure)
		}

		for i, callee := range registration.callees {
			if callee == id {
				if len(registration.callees) == 1 {
					registration.callees = make([]uint64, 0)
				} else {
					registration.callees = append(registration.callees[:i], registration.callees[i+1:]...)
				}
			}
		}
	}

	delete(d.registrationsBySession, id)
	delete(d.sessions, id)

	return nil
}

func (d *Dealer) HasProcedure(procedure string) bool {
	d.Lock()
	defer d.Unlock()

	reg, exists := d.registrationsByProcedure[procedure]
	return exists && len(reg.Registrants) > 0
}

func (d *Dealer) AutoDiscloseCaller(disclose bool) {
	d.Lock()
	defer d.Unlock()
	d.details = disclose
}

func (d *Dealer) ReceiveMessage(sessionID uint64, msg messages.Message) (*MessageWithRecipient, error) {
	d.Lock()
	defer d.Unlock()

	switch msg.Type() {
	case messages.MessageTypeCall:
		call := msg.(*messages.Call)
		var regs *Registration
		var found bool

		regs, found = d.registrationsByProcedure[call.Procedure()]
		if !found || len(regs.Registrants) == 0 {
			if d.prefixTree.Len() > 0 {
				_, reg, ok := d.prefixTree.Root().LongestPrefix([]byte(call.Procedure()))
				if ok {
					regs, found = reg, true
				}
			}

			if !found {
				for procedure, reg := range d.wcRegistrationsByProcedure {
					if wildcardMatch(call.Procedure(), procedure) {
						regs, found = reg, true
						break
					}
				}
			}
		}

		if !found || len(regs.Registrants) == 0 {
			callErr := messages.NewError(messages.MessageTypeCall, call.RequestID(), map[string]any{},
				"wamp.error.no_such_procedure", nil, nil)
			return &MessageWithRecipient{Message: callErr, Recipient: sessionID}, nil
		}

		var calleeID uint64
		if len(regs.callees) > 1 {
			switch regs.InvocationPolicy {
			case InvokeFirst:
				calleeID = regs.callees[0]
			case InvokeLast:
				calleeID = regs.callees[len(regs.callees)-1]
			case InvokeRoundRobin:
				if regs.nextCallee >= len(regs.callees) {
					regs.nextCallee = 0
				}
				calleeID = regs.callees[regs.nextCallee]
				regs.nextCallee++
			case InvokeRandom:
				idx := rand.Intn(len(regs.callees)) // #nosec
				calleeID = regs.callees[idx]
			default:
				fmt.Printf("multiple callees registered with '%s' policy", InvokeSingle)
				calleeID = regs.callees[0]
			}
		} else {
			calleeID = regs.callees[0]
		}

		receiveProgress, _ := call.Options()[OptionReceiveProgress].(bool)
		progress, _ := call.Options()[OptionProgress].(bool)

		invocationID, ok := d.invocationIDbyCall[CallMap{CallerID: sessionID, CallID: call.RequestID()}]
		if !ok || !progress {
			invocationID = d.idGen.NextID()
			d.pendingCalls[invocationID] = &PendingInvocation{
				RequestID:       call.RequestID(),
				CallerID:        sessionID,
				CalleeID:        calleeID,
				ReceiveProgress: receiveProgress,
				Progress:        progress,
			}
			d.invocationIDbyCall[CallMap{CallerID: sessionID, CallID: call.RequestID()}] = invocationID
		}

		details := map[string]any{}
		if receiveProgress {
			details[OptionReceiveProgress] = receiveProgress
		}

		if progress {
			details[OptionProgress] = progress
		}

		if d.details {
			caller := d.sessions[sessionID]
			details["procedure"] = call.Procedure()
			details["caller"] = sessionID
			details["caller_authid"] = caller.AuthID()
			details["caller_authrole"] = caller.AuthRole()
		}

		var invocation *messages.Invocation
		callee := d.sessions[calleeID]
		if callee == nil {
			return nil, fmt.Errorf("call: callee %d gone before sending invocation", calleeID)
		}

		if call.PayloadIsBinary() && callee.StaticSerializer() {
			invocation = messages.NewInvocationBinary(invocationID, regs.ID, details, call.Payload(),
				call.PayloadSerializer())
		} else {
			invocation = messages.NewInvocation(invocationID, regs.ID, details, call.Args(), call.KwArgs())
		}

		return &MessageWithRecipient{Message: invocation, Recipient: calleeID}, nil
	case messages.MessageTypeYield:
		yield := msg.(*messages.Yield)
		pending, exists := d.pendingCalls[yield.RequestID()]
		if !exists {
			return nil, fmt.Errorf("yield: no pending call for invocation %d", yield.RequestID())
		}

		switch pending.CancelMode {
		case CancelModeSkip, CancelModeKillNoWait:
			// Caller already received ERROR; discard. For skip, the callee never got
			// INTERRUPT so keep the pending call alive until the final YIELD.
			progress, _ := yield.Options()[OptionProgress].(bool)
			if !progress {
				delete(d.pendingCalls, yield.RequestID())
				delete(d.invocationIDbyCall, CallMap{CallerID: pending.CallerID, CallID: pending.RequestID})
			}
			return nil, nil
		case CancelModeKill:
			// Callee ignored the INTERRUPT and sent YIELD; deliver canceled error to caller.
			delete(d.pendingCalls, yield.RequestID())
			delete(d.invocationIDbyCall, CallMap{CallerID: pending.CallerID, CallID: pending.RequestID})
			if d.sessions[pending.CallerID] == nil {
				return nil, nil
			}
			errMsg := messages.NewError(messages.MessageTypeCall, pending.RequestID, nil, ErrCanceled, nil, nil)
			return &MessageWithRecipient{Message: errMsg, Recipient: pending.CallerID}, nil
		}

		caller := d.sessions[pending.CallerID]
		if caller == nil {
			// Caller disconnected; interrupt the callee. Mark as killnowait so any
			// further yields from a progressive call are discarded rather than erroring.
			pending.CancelMode = CancelModeKillNoWait
			return &MessageWithRecipient{Message: messages.NewInterrupt(yield.RequestID(),
				map[string]any{OptionReason: ErrCanceled, OptionMode: CancelModeKillNoWait}),
				Recipient: sessionID,
			}, nil
		}

		progress, _ := yield.Options()[OptionProgress].(bool)
		var details map[string]any
		if pending.ReceiveProgress && progress {
			details = map[string]any{OptionProgress: progress}
		} else {
			delete(d.pendingCalls, yield.RequestID())
			delete(d.invocationIDbyCall, CallMap{CallerID: pending.CallerID, CallID: pending.RequestID})
		}

		var result *messages.Result
		if yield.PayloadIsBinary() && caller.StaticSerializer() {
			result = messages.NewResultBinary(pending.RequestID, details, yield.Payload(), yield.PayloadSerializer())
		} else {
			result = messages.NewResult(pending.RequestID, details, yield.Args(), yield.KwArgs())
		}

		return &MessageWithRecipient{Message: result, Recipient: pending.CallerID}, nil
	case messages.MessageTypeRegister:
		register := msg.(*messages.Register)
		_, exists := d.registrationsBySession[sessionID]
		if !exists {
			return nil, fmt.Errorf("cannot register procedure for non-existent session %d", sessionID)
		}

		invokePolicy := util.ToString(register.Options()[OptionInvoke])
		registration, exists := d.registrationsByProcedure[register.Procedure()]
		if exists {
			if registration.InvocationPolicy == "" || registration.InvocationPolicy == InvokeSingle ||
				registration.InvocationPolicy != invokePolicy {
				err := messages.NewError(messages.MessageTypeRegister, register.RequestID(), map[string]any{},
					"wamp.error.procedure_already_exists", nil, nil)
				return &MessageWithRecipient{Message: err, Recipient: sessionID}, nil
			}
			registration.Registrants[sessionID] = sessionID
			registration.callees = append(registration.callees, sessionID)

		} else {
			registration = &Registration{
				ID:               d.idGen.NextID(),
				Procedure:        register.Procedure(),
				Registrants:      map[uint64]uint64{sessionID: sessionID},
				callees:          []uint64{sessionID},
				InvocationPolicy: invokePolicy,
			}

			match := util.ToString(register.Options()[OptionMatch])
			switch match {
			case MatchPrefix:
				registration.Match = match
				d.prefixTree, _, _ = d.prefixTree.Insert([]byte(registration.Procedure), registration)
			case MatchWildcard:
				registration.Match = match
				d.wcRegistrationsByProcedure[registration.Procedure] = registration
			default:
				registration.Match = MatchExact
			}
		}

		d.registrationsByProcedure[register.Procedure()] = registration
		d.registrationsBySession[sessionID][registration.ID] = registration

		registered := messages.NewRegistered(register.RequestID(), registration.ID)
		return &MessageWithRecipient{Message: registered, Recipient: sessionID}, nil
	case messages.MessageTypeUnregister:
		unregister := msg.(*messages.Unregister)
		registrations, exists := d.registrationsBySession[sessionID]
		if !exists || len(registrations) == 0 {
			return nil, fmt.Errorf("unregister: session %d has no registration %d", sessionID,
				unregister.RegistrationID())
		}

		registration := registrations[unregister.RegistrationID()]
		delete(registration.Registrants, sessionID)

		if len(registration.Registrants) == 0 {
			delete(registrations, unregister.RegistrationID())
			delete(d.registrationsByProcedure, registration.Procedure)
			if registration.Match == MatchPrefix {
				d.prefixTree.Delete([]byte(registration.Procedure))
			}
			if registration.Match == MatchWildcard {
				delete(d.wcRegistrationsByProcedure, registration.Procedure)
			}
		}

		unregistered := messages.NewUnregistered(unregister.RequestID())
		return &MessageWithRecipient{Message: unregistered, Recipient: sessionID}, nil
	case messages.MessageTypeError:
		wErr := msg.(*messages.Error)
		if wErr.MessageType() != messages.MessageTypeInvocation {
			return nil, fmt.Errorf("dealer: only expected to receive error in response to invocation")
		}

		pending, exists := d.pendingCalls[wErr.RequestID()]
		if !exists {
			return nil, fmt.Errorf("dealer: no pending invocation for %d", wErr.RequestID())
		}

		delete(d.pendingCalls, wErr.RequestID())
		delete(d.invocationIDbyCall, CallMap{CallerID: pending.CallerID, CallID: pending.RequestID})

		if pending.CancelMode == CancelModeSkip || pending.CancelMode == CancelModeKillNoWait {
			// Caller already received ERROR.
			return nil, nil
		}

		wErr = messages.NewError(messages.MessageTypeCall, pending.RequestID, wErr.Details(), wErr.URI(),
			wErr.Args(), wErr.KwArgs())
		return &MessageWithRecipient{Message: wErr, Recipient: pending.CallerID}, nil
	default:
		return nil, fmt.Errorf("dealer: received unexpected message of type %T", msg)
	}
}

func (d *Dealer) ReceiveCancel(sessionID uint64, cancel *messages.Cancel) ([]*MessageWithRecipient, error) {
	d.Lock()
	defer d.Unlock()

	mode := util.ToString(cancel.Options()[OptionMode])
	switch mode {
	case CancelModeSkip, CancelModeKill, CancelModeKillNoWait:
	case "":
		mode = CancelModeKillNoWait
	default:
		errMsg := messages.NewError(messages.MessageTypeCancel, cancel.RequestID(), nil, ErrInvalidArgument,
			[]any{fmt.Sprintf("invalid cancel mode: %s", mode)}, nil)
		return []*MessageWithRecipient{{Message: errMsg, Recipient: sessionID}}, nil
	}

	callMap := CallMap{CallerID: sessionID, CallID: cancel.RequestID()}
	invocationID, ok := d.invocationIDbyCall[callMap]
	if !ok {
		errMsg := messages.NewError(messages.MessageTypeCancel, cancel.RequestID(), nil, ErrInvalidArgument, nil, nil)
		return []*MessageWithRecipient{{Message: errMsg, Recipient: sessionID}}, nil
	}

	pendingCall, ok := d.pendingCalls[invocationID]
	if !ok {
		errMsg := messages.NewError(messages.MessageTypeCancel, cancel.RequestID(), nil, ErrInvalidArgument, nil, nil)
		return []*MessageWithRecipient{{Message: errMsg, Recipient: sessionID}}, nil
	}

	// Mark the pending call with the cancel mode so YIELD/ERROR handlers know how to
	// handle any future messages from the callee.
	pendingCall.CancelMode = mode

	var msgs []*MessageWithRecipient
	if mode != CancelModeSkip {
		msgs = append(msgs, &MessageWithRecipient{
			Message:   messages.NewInterrupt(invocationID, map[string]any{OptionReason: ErrCanceled, OptionMode: mode}),
			Recipient: pendingCall.CalleeID,
		})
	}

	if mode != CancelModeKill {
		// skip and killnowait: caller gets its error immediately.
		msgs = append(msgs, &MessageWithRecipient{
			Message:   messages.NewError(messages.MessageTypeCall, cancel.RequestID(), nil, ErrCanceled, nil, nil),
			Recipient: pendingCall.CallerID,
		})
	}

	return msgs, nil
}

func wildcardMatch(str, pattern string) bool {
	matched, err := path.Match(pattern, str)
	return err == nil && matched
}
