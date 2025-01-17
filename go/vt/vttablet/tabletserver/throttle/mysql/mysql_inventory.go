/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This codebase originates from https://github.com/github/freno, See https://github.com/github/freno/blob/master/LICENSE
/*
	MIT License

	Copyright (c) 2017 GitHub

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package mysql

import (
	"vitess.io/vitess/go/vt/vttablet/tabletserver/throttle/base"
)

// TabletResultMap maps a tablet to a result
type TabletResultMap map[string]base.MetricResultMap

func (m TabletResultMap) Split(alias string) (withAlias TabletResultMap, all TabletResultMap) {
	withAlias = make(TabletResultMap)
	if val, ok := m[alias]; ok {
		withAlias[alias] = val
	}
	return withAlias, m
}

// Inventory has the operational data about probes, their metrics, and relevant configuration
type Inventory struct {
	ClustersProbes       Probes
	IgnoreHostsCount     int
	IgnoreHostsThreshold float64
	TabletMetrics        TabletResultMap
}

// NewInventory creates a Inventory
func NewInventory() *Inventory {
	inventory := &Inventory{
		TabletMetrics: make(TabletResultMap),
	}
	return inventory
}
