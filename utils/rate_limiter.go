// Copyright (c) 2017, 2018 Uber Technologies, Inc.
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

package utils

import (
	"sync"
	"time"
)

// RateLimiter is a filter used to check if a message that is worth itemCost units is within the rate limits.
type RateLimiter interface {
	CheckCredit(itemCost float64) bool
	Drain() float64
}

type rateLimiter struct {
	sync.Mutex

	creditsPerSecond float64
	balance          float64
	maxBalance       float64
	lastTick         time.Time

	timeNow func() time.Time
}

// NewRateLimiter creates a new rate limiter based on leaky bucket algorithm, formulated in terms of a
// credits balance that is replenished every time CheckCredit() method is called (tick) by the amount proportional
// to the time elapsed since the last tick, up to max of creditsPerSecond. A call to CheckCredit() takes a cost
// of an item we want to pay with the balance. If the balance exceeds the cost of the item, the item is "purchased"
// and the balance reduced, indicated by returned value of true. Otherwise the balance is unchanged and return false.
//
// This can be used to limit a rate of messages emitted by a service by instantiating the Rate Limiter with the
// max number of messages a service is allowed to emit per second, and calling CheckCredit(1.0) for each message
// to determine if the message is within the rate limit.
//
// It can also be used to limit the rate of traffic in bytes, by setting creditsPerSecond to desired throughput
// as bytes/second, and calling CheckCredit() with the actual message size.
func NewRateLimiter(creditsPerSecond, maxBalance float64) RateLimiter {
	return &rateLimiter{
		creditsPerSecond: creditsPerSecond,
		balance:          maxBalance,
		maxBalance:       maxBalance,
		lastTick:         time.Now(),
		timeNow:          time.Now}
}

// N.B. Call while holding lock
func (r *rateLimiter) updateBalance() {
	// calculate how much time passed since the last tick, and update current tick
	currentTime := r.timeNow()
	elapsedTime := currentTime.Sub(r.lastTick)
	r.lastTick = currentTime
	// calculate how much credit have we accumulated since the last tick
	r.balance += elapsedTime.Seconds() * r.creditsPerSecond
	if r.balance > r.maxBalance {
		r.balance = r.maxBalance
	}
}

func (r *rateLimiter) CheckCredit(itemCost float64) bool {
	r.Lock()
	defer r.Unlock()
	r.updateBalance()
	// if we have enough credits to pay for current item, then reduce balance and allow
	if r.balance >= itemCost {
		r.balance -= itemCost
		return true
	}
	return false
}

// Drain returns the balance and sets it to zero internally
func (r *rateLimiter) Drain() float64 {
	r.Lock()
	defer r.Unlock()
	r.updateBalance()
	result := r.balance
	r.balance = 0
	return result
}
