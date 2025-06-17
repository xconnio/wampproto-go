package wampproto

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-immutable-radix/v2"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/util"
)

const (
	OptionReceiveProgress = "receive_progress"
	OptionProgress        = "progress"
	OptionMatch           = "match"

	MatchExact  = "exact"
	MatchPrefix = "prefix"
)

const (
	FeatureProgressiveCallInvocations = "progressive_call_invocations"
	FeatureProgressiveCallResults     = "progressive_call_results"
	FeatureCallCancelling             = "call_canceling"
)

type PendingInvocation struct {
	RequestID       int64
	CallerID        int64
	CalleeID        int64
	Progress        bool
	ReceiveProgress bool
}

type Registration struct {
	ID               int64
	Procedure        string
	Registrants      map[int64]int64
	InvocationPolicy string
	Match            string
}

type CallMap struct {
	CallerID int64
	CallID   int64
}

type Dealer struct {
	sessions                 map[int64]*SessionDetails
	registrationsByProcedure map[string]*Registration
	registrationsBySession   map[int64]map[int64]*Registration
	prefixTree               *iradix.Tree[*Registration]
	pendingCalls             map[int64]*PendingInvocation
	invocationIDbyCall       map[CallMap]int64

	idGen *SessionScopeIDGenerator
	sync.Mutex
}

func NewDealer() *Dealer {
	return &Dealer{
		sessions:                 make(map[int64]*SessionDetails),
		registrationsByProcedure: make(map[string]*Registration),
		registrationsBySession:   make(map[int64]map[int64]*Registration),
		pendingCalls:             make(map[int64]*PendingInvocation),
		invocationIDbyCall:       make(map[CallMap]int64),
		idGen:                    &SessionScopeIDGenerator{},
		prefixTree:               iradix.New[*Registration](),
	}
}

func (d *Dealer) AddSession(details *SessionDetails) error {
	d.Lock()
	defer d.Unlock()

	_, exists := d.sessions[details.ID()]
	if exists {
		return fmt.Errorf("cannot attach an already attached client %d", details.ID())
	}

	d.registrationsBySession[details.ID()] = map[int64]*Registration{}
	d.sessions[details.ID()] = details
	return nil
}

func (d *Dealer) RemoveSession(id int64) error {
	d.Lock()
	defer d.Unlock()

	_, exists := d.sessions[id]
	if !exists {
		return fmt.Errorf("cannot remove client with id %d not attached", id)
	}

	registrations := d.registrationsBySession[id]
	for _, registration := range registrations {
		registration = d.registrationsByProcedure[registration.Procedure]
		delete(registration.Registrants, id)
		if len(registration.Registrants) == 0 {
			delete(d.registrationsByProcedure, registration.Procedure)
		}
		if registration.Match == MatchPrefix {
			d.prefixTree.Delete([]byte(registration.Procedure))
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

func (d *Dealer) ReceiveMessage(sessionID int64, msg messages.Message) (*MessageWithRecipient, error) {
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
				found = ok
				if ok {
					regs = reg
				}
			}
		}

		if !found || len(regs.Registrants) == 0 {
			callErr := messages.NewError(messages.MessageTypeCall, call.RequestID(), map[string]any{},
				"wamp.error.no_such_procedure", nil, nil)
			return &MessageWithRecipient{Message: callErr, Recipient: sessionID}, nil
		}

		var callee int64
		for session := range regs.Registrants {
			callee = session
			break
		}
		receiveProgress, _ := call.Options()[OptionReceiveProgress].(bool)
		progress, _ := call.Options()[OptionProgress].(bool)

		invocationID, ok := d.invocationIDbyCall[CallMap{CallerID: sessionID, CallID: call.RequestID()}]
		if !ok || !progress {
			invocationID = d.idGen.NextID()
			d.pendingCalls[invocationID] = &PendingInvocation{
				RequestID:       call.RequestID(),
				CallerID:        sessionID,
				CalleeID:        callee,
				ReceiveProgress: receiveProgress,
				Progress:        progress,
			}
			d.invocationIDbyCall[CallMap{CallerID: sessionID, CallID: call.RequestID()}] = invocationID
		}

		var invocation *messages.Invocation
		if call.PayloadIsBinary() && d.sessions[callee].StaticSerializer() {
			invocation = messages.NewInvocationBinary(invocationID, regs.ID, nil, call.Payload(),
				call.PayloadSerializer())
		} else {
			details := map[string]any{}
			if receiveProgress {
				details[OptionReceiveProgress] = receiveProgress
			}
			if progress {
				details[OptionProgress] = progress
			}
			invocation = messages.NewInvocation(invocationID, regs.ID, details, call.Args(), call.KwArgs())
		}

		return &MessageWithRecipient{Message: invocation, Recipient: callee}, nil
	case messages.MessageTypeYield:
		yield := msg.(*messages.Yield)
		pending, exists := d.pendingCalls[yield.RequestID()]
		if !exists {
			return nil, fmt.Errorf("yield: not pending calls for session %d", sessionID)
		}

		progress, _ := yield.Options()[OptionProgress].(bool)
		var details map[string]any
		if pending.ReceiveProgress && progress {
			details = map[string]any{OptionProgress: progress}
		} else {
			delete(d.pendingCalls, yield.RequestID())
		}

		var result *messages.Result
		if yield.PayloadIsBinary() && d.sessions[pending.CallerID].StaticSerializer() {
			result = messages.NewResultBinary(pending.RequestID, nil, yield.Payload(), yield.PayloadSerializer())
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

		registration, exists := d.registrationsByProcedure[register.Procedure()]
		if exists {
			// TODO: implement shared registrations
			err := messages.NewError(messages.MessageTypeRegister, register.RequestID(), map[string]any{},
				"wamp.error.procedure_already_exists", nil, nil)
			return &MessageWithRecipient{Message: err, Recipient: sessionID}, nil
		} else {
			registration = &Registration{
				ID:          d.idGen.NextID(),
				Procedure:   register.Procedure(),
				Registrants: map[int64]int64{sessionID: sessionID},
			}
			match := util.ToString(register.Options()[OptionMatch])
			switch match {
			case MatchPrefix:
				registration.Match = match
				d.prefixTree, _, _ = d.prefixTree.Insert([]byte(registration.Procedure), registration)
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

		wErr = messages.NewError(messages.MessageTypeCall, pending.RequestID, wErr.Details(), wErr.URI(),
			wErr.Args(), wErr.KwArgs())
		return &MessageWithRecipient{Message: wErr, Recipient: pending.CallerID}, nil
	default:
		return nil, fmt.Errorf("dealer: received unexpected message of type %T", msg)
	}
}
