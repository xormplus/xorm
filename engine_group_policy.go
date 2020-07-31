// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"math/rand"
	"sync"
	"time"
)

// GroupPolicy is be used by chosing the current subordinate from subordinates
type GroupPolicy interface {
	Subordinate(*EngineGroup) *Engine
}

// GroupPolicyHandler should be used when a function is a GroupPolicy
type GroupPolicyHandler func(*EngineGroup) *Engine

// Subordinate implements the chosen of subordinates
func (h GroupPolicyHandler) Subordinate(eg *EngineGroup) *Engine {
	return h(eg)
}

// RandomPolicy implmentes randomly chose the subordinate of subordinates
func RandomPolicy() GroupPolicyHandler {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(g *EngineGroup) *Engine {
		return g.Subordinates()[r.Intn(len(g.Subordinates()))]
	}
}

// WeightRandomPolicy implmentes randomly chose the subordinate of subordinates
func WeightRandomPolicy(weights []int) GroupPolicyHandler {
	var rands = make([]int, 0, len(weights))
	for i := 0; i < len(weights); i++ {
		for n := 0; n < weights[i]; n++ {
			rands = append(rands, i)
		}
	}
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))

	return func(g *EngineGroup) *Engine {
		var subordinates = g.Subordinates()
		idx := rands[r.Intn(len(rands))]
		if idx >= len(subordinates) {
			idx = len(subordinates) - 1
		}
		return subordinates[idx]
	}
}

// RoundRobinPolicy returns a group policy handler
func RoundRobinPolicy() GroupPolicyHandler {
	var pos = -1
	var lock sync.Mutex
	return func(g *EngineGroup) *Engine {
		var subordinates = g.Subordinates()

		lock.Lock()
		defer lock.Unlock()
		pos++
		if pos >= len(subordinates) {
			pos = 0
		}

		return subordinates[pos]
	}
}

// WeightRoundRobinPolicy returns a group policy handler
func WeightRoundRobinPolicy(weights []int) GroupPolicyHandler {
	var rands = make([]int, 0, len(weights))
	for i := 0; i < len(weights); i++ {
		for n := 0; n < weights[i]; n++ {
			rands = append(rands, i)
		}
	}
	var pos = -1
	var lock sync.Mutex

	return func(g *EngineGroup) *Engine {
		var subordinates = g.Subordinates()
		lock.Lock()
		defer lock.Unlock()
		pos++
		if pos >= len(rands) {
			pos = 0
		}

		idx := rands[pos]
		if idx >= len(subordinates) {
			idx = len(subordinates) - 1
		}
		return subordinates[idx]
	}
}

// LeastConnPolicy implements GroupPolicy, every time will get the least connections subordinate
func LeastConnPolicy() GroupPolicyHandler {
	return func(g *EngineGroup) *Engine {
		var subordinates = g.Subordinates()
		connections := 0
		idx := 0
		for i := 0; i < len(subordinates); i++ {
			openConnections := subordinates[i].DB().Stats().OpenConnections
			if i == 0 {
				connections = openConnections
				idx = i
			} else if openConnections <= connections {
				connections = openConnections
				idx = i
			}
		}
		return subordinates[idx]
	}
}
