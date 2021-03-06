package dht

import (
	"net"

	"github.com/optmzr/d7024e-dht/route"

	"github.com/optmzr/d7024e-dht/network"
	"github.com/optmzr/d7024e-dht/node"
	"github.com/optmzr/d7024e-dht/store"
)

type Call interface {
	Do(nw network.Network, address net.UDPAddr) (ch chan network.FindResult, err error)
	Result(result network.FindResult, callee route.Contact) (stop bool)
	Target() (target node.ID)
}

func NewFindNodesCall(target node.ID) *FindNodesCall {
	return &FindNodesCall{
		target: target,
	}
}

type FindNodesCall struct {
	target node.ID
}

func (q *FindNodesCall) Do(nw network.Network, address net.UDPAddr) (chan network.FindResult, error) {
	return nw.FindNodes(q.target, address)
}

func (q *FindNodesCall) Result(_ network.FindResult, _ route.Contact) (_ bool) { return }
func (q *FindNodesCall) Target() node.ID                                       { return q.target }

func NewFindValueCall(hash store.Key) *FindValueCall {
	return &FindValueCall{
		hash: hash,
	}
}

type FindValueCall struct {
	hash   store.Key
	value  string
	sender node.ID
}

func (q *FindValueCall) Do(nw network.Network, address net.UDPAddr) (chan network.FindResult, error) {
	return nw.FindValue(q.hash, address)
}

func (q *FindValueCall) Result(result network.FindResult, callee route.Contact) (stop bool) {
	// TODO: Value validation could be added here, where the value received is
	// checked towards the expected hash.

	q.value = result.Value()
	q.sender = callee.NodeID
	if q.value != "" {
		stop = true
	} else {
		stop = false
	}

	return
}

func (q *FindValueCall) Target() node.ID { return node.ID(q.hash) }
