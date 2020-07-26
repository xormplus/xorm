// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"time"

	"github.com/xormplus/xorm/caches"
	"github.com/xormplus/xorm/contexts"
	"github.com/xormplus/xorm/dialects"
	"github.com/xormplus/xorm/log"
	"github.com/xormplus/xorm/names"
)

// EngineGroup defines an engine group
type EngineGroup struct {
	*Engine
	subordinates []*Engine
	policy GroupPolicy
}

// NewEngineGroup creates a new engine group
func NewEngineGroup(args1 interface{}, args2 interface{}, policies ...GroupPolicy) (*EngineGroup, error) {
	var eg EngineGroup
	if len(policies) > 0 {
		eg.policy = policies[0]
	} else {
		eg.policy = RoundRobinPolicy()
	}

	driverName, ok1 := args1.(string)
	conns, ok2 := args2.([]string)
	if ok1 && ok2 {
		engines := make([]*Engine, len(conns))
		for i, conn := range conns {
			engine, err := NewEngine(driverName, conn)
			if err != nil {
				return nil, err
			}
			engine.engineGroup = &eg
			engines[i] = engine
		}

		eg.Engine = engines[0]
		eg.subordinates = engines[1:]
		return &eg, nil
	}

	main, ok3 := args1.(*Engine)
	subordinates, ok4 := args2.([]*Engine)
	if ok3 && ok4 {
		main.engineGroup = &eg
		for i := 0; i < len(subordinates); i++ {
			subordinates[i].engineGroup = &eg
		}
		eg.Engine = main
		eg.subordinates = subordinates
		return &eg, nil
	}
	return nil, ErrParamsType
}

// Close the engine
func (eg *EngineGroup) Close() error {
	err := eg.Engine.Close()
	if err != nil {
		return err
	}

	for i := 0; i < len(eg.subordinates); i++ {
		err := eg.subordinates[i].Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// ContextHook returned a group session
func (eg *EngineGroup) Context(ctx context.Context) *Session {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.Context(ctx)
}

// NewSession returned a group session
func (eg *EngineGroup) NewSession() *Session {
	sess := eg.Engine.NewSession()
	sess.sessionType = groupSession
	return sess
}

// Main returns the main engine
func (eg *EngineGroup) Main() *Engine {
	return eg.Engine
}

// Ping tests if database is alive
func (eg *EngineGroup) Ping() error {
	if err := eg.Engine.Ping(); err != nil {
		return err
	}

	for _, subordinate := range eg.subordinates {
		if err := subordinate.Ping(); err != nil {
			return err
		}
	}
	return nil
}

// SetColumnMapper set the column name mapping rule
func (eg *EngineGroup) SetColumnMapper(mapper names.Mapper) {
	eg.Engine.SetColumnMapper(mapper)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetColumnMapper(mapper)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (eg *EngineGroup) SetConnMaxLifetime(d time.Duration) {
	eg.Engine.SetConnMaxLifetime(d)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetConnMaxLifetime(d)
	}
}

// SetDefaultCacher set the default cacher
func (eg *EngineGroup) SetDefaultCacher(cacher caches.Cacher) {
	eg.Engine.SetDefaultCacher(cacher)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetDefaultCacher(cacher)
	}
}

// SetLogger set the new logger
func (eg *EngineGroup) SetLogger(logger interface{}) {
	eg.Engine.SetLogger(logger)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetLogger(logger)
	}
}

func (eg *EngineGroup) AddHook(hook contexts.Hook) {
	eg.Engine.AddHook(hook)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].AddHook(hook)
	}
}

// SetLogLevel sets the logger level
func (eg *EngineGroup) SetLogLevel(level log.LogLevel) {
	eg.Engine.SetLogLevel(level)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetLogLevel(level)
	}
}

// SetMapper set the name mapping rules
func (eg *EngineGroup) SetMapper(mapper names.Mapper) {
	eg.Engine.SetMapper(mapper)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetMapper(mapper)
	}
}

// SetMaxIdleConns set the max idle connections on pool, default is 2
func (eg *EngineGroup) SetMaxIdleConns(conns int) {
	eg.Engine.DB().SetMaxIdleConns(conns)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].DB().SetMaxIdleConns(conns)
	}
}

// SetMaxOpenConns is only available for go 1.2+
func (eg *EngineGroup) SetMaxOpenConns(conns int) {
	eg.Engine.DB().SetMaxOpenConns(conns)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].DB().SetMaxOpenConns(conns)
	}
}

// SetPolicy set the group policy
func (eg *EngineGroup) SetPolicy(policy GroupPolicy) *EngineGroup {
	eg.policy = policy
	return eg
}

// SetQuotePolicy sets the special quote policy
func (eg *EngineGroup) SetQuotePolicy(quotePolicy dialects.QuotePolicy) {
	eg.Engine.SetQuotePolicy(quotePolicy)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetQuotePolicy(quotePolicy)
	}
}

// SetTableMapper set the table name mapping rule
func (eg *EngineGroup) SetTableMapper(mapper names.Mapper) {
	eg.Engine.SetTableMapper(mapper)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].SetTableMapper(mapper)
	}
}

// ShowSQL show SQL statement or not on logger if log level is great than INFO
func (eg *EngineGroup) ShowSQL(show ...bool) {
	eg.Engine.ShowSQL(show...)
	for i := 0; i < len(eg.subordinates); i++ {
		eg.subordinates[i].ShowSQL(show...)
	}
}

// Subordinate returns one of the physical databases which is a subordinate according the policy
func (eg *EngineGroup) Subordinate() *Engine {
	switch len(eg.subordinates) {
	case 0:
		return eg.Engine
	case 1:
		return eg.subordinates[0]
	}
	return eg.policy.Subordinate(eg)
}

// Subordinates returns all the subordinates
func (eg *EngineGroup) Subordinates() []*Engine {
	return eg.subordinates
}

func (eg *EngineGroup) RegisterSqlTemplate(sqlt SqlTemplate, Cipher ...Cipher) error {
	err := eg.Engine.RegisterSqlTemplate(sqlt, Cipher...)
	if err != nil {
		return err
	}
	for i := 0; i < len(eg.subordinates); i++ {
		err = eg.subordinates[i].RegisterSqlTemplate(sqlt, Cipher...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (eg *EngineGroup) RegisterSqlMap(sqlm SqlM, Cipher ...Cipher) error {
	err := eg.Engine.RegisterSqlMap(sqlm, Cipher...)
	if err != nil {
		return err
	}
	for i := 0; i < len(eg.subordinates); i++ {
		err = eg.subordinates[i].RegisterSqlMap(sqlm, Cipher...)
		if err != nil {
			return err
		}
	}
	return nil
}
