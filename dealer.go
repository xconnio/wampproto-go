package wampproto

import (
	"fmt"
	"math/rand"
	"path"
	"sync"

	"github.com/hashicorp/go-immutable-radix/v2"
	log "github.com/sirupsen/logrus"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/util"
)

const (
	OptionReceiveProgress = "receive_progress"
	OptionProgress        = "progress"
	OptionMatch           = "match"
	OptionInvoke          = "invoke"

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
	FeatureCallCancelling             = "call_canceling"
	FeaturePublisherExclusion         = "publisher_exclusion"
)

type PendingInvocation struct {
	RequestID       uint64
	CallerID        uint64
	CalleeID        uint64
	Progress        bool
	ReceiveProgress bool
}

type Registration struct {
	ID               uint64
	Procedure        string
	Registrants      map[uint64]uint64
	InvocationPolicy string
	nextCallee       int
	callees          []uint64
	Match            string
	Created          string
}

type CallMap struct {
	CallerID uint64
	CallID   uint64
}

type RegistrationEvent struct {
	SessionID      uint64
	RegistrationID uint64
}

type Dealer struct {
	sessions                         map[uint64]*SessionDetails
	registrationsByProcedure         map[string]*Registration
	registrationsBySession           map[uint64]map[uint64]*Registration
	prefixTree                       *iradix.Tree[*Registration]
	wildcardRegistrationsByProcedure map[string]*Registration
	pendingCalls                     map[uint64]*PendingInvocation
	invocationIDbyCall               map[CallMap]uint64
	details                          bool

	idGen *SessionScopeIDGenerator
	sync.Mutex

	RegistrationCreated chan *Registration
	CalleeAdded         chan RegistrationEvent
	CalleeRemoved       chan RegistrationEvent
	RegistrationDeleted chan RegistrationEvent

	exactRegistrationsByID    map[uint64]*Registration
	prefixRegistrationsByID   map[uint64]*Registration
	wildcardRegistrationsByID map[uint64]*Registration
	metaAPi                   bool
}

func NewDealer() *Dealer {
	return &Dealer{
		sessions:                         make(map[uint64]*SessionDetails),
		registrationsByProcedure:         make(map[string]*Registration),
		registrationsBySession:           make(map[uint64]map[uint64]*Registration),
		pendingCalls:                     make(map[uint64]*PendingInvocation),
		invocationIDbyCall:               make(map[CallMap]uint64),
		idGen:                            &SessionScopeIDGenerator{},
		prefixTree:                       iradix.New[*Registration](),
		wildcardRegistrationsByProcedure: make(map[string]*Registration),
		RegistrationCreated:              make(chan *Registration),
		CalleeAdded:                      make(chan RegistrationEvent),
		CalleeRemoved:                    make(chan RegistrationEvent),
		RegistrationDeleted:              make(chan RegistrationEvent),
		exactRegistrationsByID:           make(map[uint64]*Registration),
		prefixRegistrationsByID:          make(map[uint64]*Registration),
		wildcardRegistrationsByID:        make(map[uint64]*Registration),
	}
}

func (d *Dealer) ExactRegistrationsByID() map[uint64]*Registration {
	d.Lock()
	defer d.Unlock()
	copyMap := make(map[uint64]*Registration, len(d.exactRegistrationsByID))
	for id, reg := range d.exactRegistrationsByID {
		copyMap[id] = reg
	}

	return copyMap
}

func (d *Dealer) PrefixRegistrationsByID() map[uint64]*Registration {
	d.Lock()
	defer d.Unlock()
	copyMap := make(map[uint64]*Registration, len(d.prefixRegistrationsByID))
	for id, reg := range d.prefixRegistrationsByID {
		copyMap[id] = reg
	}

	return copyMap
}

func (d *Dealer) WildCardRegistrationsByID() map[uint64]*Registration {
	d.Lock()
	defer d.Unlock()
	copyMap := make(map[uint64]*Registration, len(d.wildcardRegistrationsByID))
	for id, reg := range d.wildcardRegistrationsByID {
		copyMap[id] = reg
	}

	return copyMap
}

func (d *Dealer) RegistrationsByProcedure() map[string]*Registration {
	d.Lock()
	defer d.Unlock()
	copyMap := make(map[string]*Registration, len(d.registrationsByProcedure))
	for id, reg := range d.registrationsByProcedure {
		copyMap[id] = reg
	}
	return copyMap
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

	registrations, ok := d.registrationsBySession[id]
	if ok {
		for _, reg := range registrations {
			registration, ok := d.registrationsByProcedure[reg.Procedure]
			if !ok {
				continue
			}
			d.removeRegistration(registration.ID, id)
		}
	} else {
		log.Debugf("Dealer: no registrations found for session %d", id)
	}

	delete(d.registrationsBySession, id)
	delete(d.sessions, id)

	return nil
}

func removeCallee(callees []uint64, id uint64) []uint64 {
	for i, callee := range callees {
		if callee == id {
			if len(callees) == 1 {
				return make([]uint64, 0)
			}
			return append(callees[:i], callees[i+1:]...)
		}
	}
	return callees
}

func (d *Dealer) removeRegistration(registrationID uint64, sessionID uint64) {
	registrations, exists := d.registrationsBySession[sessionID]
	if !exists || len(registrations) == 0 {
		return
	}

	registration, ok := registrations[registrationID]
	if !ok {
		log.Debugf("Dealer: registration %d not found for session %d", registrationID, sessionID)
		return
	}
	delete(registration.Registrants, sessionID)
	registration.callees = removeCallee(registration.callees, sessionID)

	if len(registration.Registrants) == 0 {
		delete(registrations, registrationID)
		delete(d.registrationsByProcedure, registration.Procedure)
		switch registration.Match {
		case MatchPrefix:
			delete(d.prefixRegistrationsByID, registration.ID)
			d.prefixTree, _, _ = d.prefixTree.Delete([]byte(registration.Procedure))
		case MatchWildcard:
			delete(d.wildcardRegistrationsByID, registration.ID)
			delete(d.wildcardRegistrationsByProcedure, registration.Procedure)
		default:
			delete(d.exactRegistrationsByID, registration.ID)
		}
		if d.metaAPi {
			select {
			case d.CalleeRemoved <- RegistrationEvent{SessionID: sessionID, RegistrationID: registration.ID}:
			default:
			}
			select {
			case d.RegistrationDeleted <- RegistrationEvent{SessionID: sessionID, RegistrationID: registration.ID}:
			default:
			}
		}
	} else {
		registrations[registrationID] = registration
		d.registrationsByProcedure[registration.Procedure] = registration
		switch registration.Match {
		case MatchPrefix:
			d.prefixTree, _, _ = d.prefixTree.Insert([]byte(registration.Procedure), registration)
			d.prefixRegistrationsByID[registration.ID] = registration
		case MatchWildcard:
			d.wildcardRegistrationsByProcedure[registration.Procedure] = registration
			d.wildcardRegistrationsByID[registration.ID] = registration
		default:
			d.exactRegistrationsByID[registration.ID] = registration
		}
		if d.metaAPi {
			select {
			case d.CalleeRemoved <- RegistrationEvent{SessionID: sessionID, RegistrationID: registration.ID}:
			default:
			}
		}
	}
	d.registrationsBySession[sessionID] = registrations
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

		regs, found = d.matchRegistration(call.Procedure())
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
		callee, ok := d.sessions[calleeID]
		if !ok || callee == nil {
			log.Debugf("Dealer: callee %d gone before sending invocation", calleeID)
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
		caller, ok := d.sessions[pending.CallerID]
		if !ok || caller == nil {
			log.Debugf("Dealer: caller %d gone before receiving result", pending.CallerID)
			return nil, fmt.Errorf("yield: caller %d gone before receiving result", pending.CallerID)
		}

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
				Created:          auth.NowISO8601(),
			}
			if d.metaAPi {
				select {
				case d.RegistrationCreated <- registration:
				default:
				}
			}
		}

		if d.metaAPi {
			select {
			case d.CalleeAdded <- RegistrationEvent{SessionID: sessionID, RegistrationID: registration.ID}:
			default:
			}
		}

		match := util.ToString(register.Options()[OptionMatch])
		switch match {
		case MatchPrefix:
			registration.Match = match
			d.prefixTree, _, _ = d.prefixTree.Insert([]byte(registration.Procedure), registration)
			d.prefixRegistrationsByID[registration.ID] = registration
		case MatchWildcard:
			registration.Match = match
			d.wildcardRegistrationsByProcedure[registration.Procedure] = registration
			d.wildcardRegistrationsByID[registration.ID] = registration
		default:
			registration.Match = MatchExact
			d.exactRegistrationsByID[registration.ID] = registration
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

		d.removeRegistration(unregister.RegistrationID(), sessionID)

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

func (d *Dealer) matchRegistration(procedure string) (reg *Registration, found bool) {
	if r, ok := d.registrationsByProcedure[procedure]; ok && len(r.Registrants) > 0 {
		return r, true
	}

	if d.prefixTree.Len() > 0 {
		_, reg, ok := d.prefixTree.Root().LongestPrefix([]byte(procedure))
		if ok {
			return reg, true
		}
	}

	for pattern, reg := range d.wildcardRegistrationsByProcedure {
		if wildcardMatch(procedure, pattern) {
			return reg, true
		}
	}

	return nil, false
}
func (d *Dealer) MatchRegistration(procedure string) (reg *Registration, found bool) {
	d.Lock()
	defer d.Unlock()
	return d.matchRegistration(procedure)
}

func (d *Dealer) EnableMetaAPI() {
	d.Lock()
	defer d.Unlock()
	d.metaAPi = true
}

func wildcardMatch(str, pattern string) bool {
	matched, err := path.Match(pattern, str)
	return err == nil && matched
}
