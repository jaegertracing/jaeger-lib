// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xkit

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// logger wraps a go-kit logger instance in a Jaeger client compatible one.
type logger struct {
	infoLogger  log.Logger
	errorLogger log.Logger
}

// NewLogger creates a new Jaeger client logger from a go-kit one.
func NewLogger(l log.Logger) *logger {
	return &logger{
		infoLogger:  level.Info(l),
		errorLogger: level.Error(l),
	}
}

// Error implements the github.com/uber/jaeger-client-go/log.Logger interface.
func (l *logger) Error(msg string) {
	l.errorLogger.Log("msg", msg)
}

// Infof implements the github.com/uber/jaeger-client-go/log.Logger interface.
func (l *logger) Infof(msg string, args ...interface{}) {
	l.infoLogger.Log("msg", fmt.Sprintf(msg, args...))
}
