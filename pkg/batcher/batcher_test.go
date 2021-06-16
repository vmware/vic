// Copyright 2018 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batcher

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func single(batcher Batcher, group string) error {
	batcher.Start(context.Background())
	err := batcher.QueueSync(context.Background(), group, 0, "hello")
	if err == nil {
		return nil
	}
	return err.(error)
}

func multipleMembersSingleGroup(batcher Batcher, group string, data ...interface{}) error {
	var wg sync.WaitGroup

	var err error
	for i := range data {
		wg.Add(1)
		batcher.QueueAsync(context.Background(), group, 0, data[i], func(v interface{}) {
			if err == nil && v != nil {
				err = v.(error)
			}
			wg.Done()
		})
	}

	batcher.Start(context.Background())
	wg.Wait()

	return err
}

func multipleMembersAndGroups(batcher Batcher, data map[string][]interface{}) error {
	var wg sync.WaitGroup

	var err error
	for k, v := range data {
		for i := range v {
			wg.Add(1)
			batcher.QueueAsync(context.Background(), k, 0, v[i], func(x interface{}) {
				if err == nil && x != nil {
					err = x.(error)
				}
				wg.Done()
			})
		}
	}

	batcher.Start(context.Background())
	wg.Wait()

	return err
}

type accounting struct {
	totoalMemberCount int
	groupCounts       map[string]int
	invocationCount   int
	groupInvocations  map[string]int
	groupState        map[string]interface{}
	result            []interface{}
	groupResult       map[string][]interface{}
}

func (a *accounting) init() *accounting {
	if a != nil {
		return a
	}

	return &accounting{
		groupCounts:      make(map[string]int),
		groupInvocations: make(map[string]int),
		groupResult:      make(map[string][]interface{}),
		groupState:       make(map[string]interface{}),
	}
}

// Accepts all candidates and returns the growing member set as state
func newBasicAssessor(accounts *accounting) (*accounting, Assessor) {
	accounts = accounts.init()

	assessor := func(ctx context.Context, groupID string, candidate interface{}, members []interface{}, state interface{}) (Assessment, interface{}) {
		fmt.Printf("assessing member %v for group %s\n", candidate, groupID)
		if accounts.invocationCount != 0 && state == nil {
			panic("expected state to be passed through from prior invocation")
		}
		accounts.totoalMemberCount++
		accounts.groupCounts[groupID]++
		accounts.invocationCount++
		accounts.groupInvocations[groupID]++
		accounts.result = append(accounts.result, candidate)
		accounts.groupResult[groupID] = append(accounts.groupResult[groupID], candidate)
		gresult := append(accounts.groupResult[groupID], members...)
		accounts.groupResult[groupID] = gresult
		accounts.groupState[groupID] = gresult
		return Accept, gresult
	}

	return accounts, assessor
}

// Rejects the second candidate for all groups, otherwise calls through to basic assessor
func rejectSecondAssessor(accounts *accounting) (*accounting, Assessor) {
	accounts = accounts.init()

	_, basic := newBasicAssessor(accounts)
	assessor := func(ctx context.Context, groupID string, candidate interface{}, members []interface{}, state interface{}) (Assessment, interface{}) {
		if accounts.groupInvocations[groupID] == 1 {
			fmt.Printf("rejecting second member %v for group %s\n", candidate, groupID)

			accounts.groupInvocations[groupID]++
			accounts.invocationCount++

			return RejectImmediate, nil
		}

		return basic(ctx, groupID, candidate, members, state)
	}

	return accounts, assessor
}

// Does simple accounts and aggregation for the members that were approved
func newBasicProcessor(accounts *accounting) (*accounting, Processor) {
	accounts = accounts.init()

	processor := func(groupID string, members []interface{}, state interface{}, ctxs Contexts) interface{} {
		fmt.Printf("processing %d members for group %s\n", len(members), groupID)
		accounts.totoalMemberCount = accounts.totoalMemberCount + len(members)
		accounts.groupCounts[groupID] += len(members)
		accounts.invocationCount++
		accounts.result = append(accounts.result, members...)
		gresult := append(accounts.groupResult[groupID], members...)
		accounts.groupResult[groupID] = gresult
		accounts.groupState[groupID] = gresult
		return nil
	}

	return accounts, processor
}

func TestSingleItem(t *testing.T) {
	accounts, processor := newBasicProcessor(nil)

	batcher := NewBatcher(context.Background(), nil, processor, false)
	err := single(batcher, t.Name())

	require.Nil(t, err, "Expected no error")
	require.Equal(t, 1, accounts.totoalMemberCount, "Expected only one group member to have been processed")
}

func TestMultipleItems(t *testing.T) {
	accounts, processor := newBasicProcessor(nil)

	batcher := NewBatcher(context.Background(), nil, processor, false)
	err := multipleMembersSingleGroup(batcher, t.Name(), "hello", "world")

	require.Nil(t, err, "Expected no error")
	require.Equal(t, 2, accounts.totoalMemberCount, "Expected exactly two group members to have been processed")
	require.Equal(t, 1, accounts.invocationCount, "Expected exactly one batch to be processed")

	// There is no guaranteed ordering in the batch so we cannot check for the resulting string directly
	resultString := ""
	for _, v := range accounts.result {
		resultString += v.(string)
	}
	require.True(t, strings.Contains(resultString, "hello"), "Expected values to propagate correctly")
	require.True(t, strings.Contains(resultString, "world"), "Expected values to propagate correctly")
}

func TestMultipleItemsWithAssessor(t *testing.T) {
	paccounts, processor := newBasicProcessor(nil)
	aaccounts, assessor := newBasicAssessor(nil)

	batcher := NewBatcher(context.Background(), assessor, processor, false)
	err := multipleMembersSingleGroup(batcher, t.Name(), "hello", "world")

	require.Nil(t, err, "Expected no error")

	require.Equal(t, 2, aaccounts.totoalMemberCount, "Expected exactly two group members to have been assessed")
	require.Equal(t, 2, aaccounts.invocationCount, "Expected exactly one invocation of assessor per member")

	require.Equal(t, 2, paccounts.totoalMemberCount, "Expected exactly two group members to have been processed")
	require.Equal(t, 1, paccounts.invocationCount, "Expected exactly one batch to be processed")
	require.NotNil(t, paccounts.groupState, "Expected state to be passed through to processor from assessor")

	// There is no guaranteed ordering in the batch so we cannot check for the resulting string directly
	resultString := ""
	for _, v := range paccounts.result {
		resultString += v.(string)
	}
	require.True(t, strings.Contains(resultString, "hello"), "Expected values to propagate correctly")
	require.True(t, strings.Contains(resultString, "world"), "Expected values to propagate correctly")
}

func TestMultipleItemsAndGroups(t *testing.T) {
	accounts, processor := newBasicProcessor(nil)

	batcher := NewBatcher(context.Background(), nil, processor, false)
	colours := []interface{}{"red", "green", "blue"}
	greetings := []interface{}{"hello", "bonjour", "hola", "konnichiwa"}
	denied := []interface{}{"no", "non", "no", "iie"}
	aggregate := map[string][]interface{}{
		"colours":   colours,
		"greetings": greetings,
		"denied":    denied,
	}

	err := multipleMembersAndGroups(batcher, aggregate)

	require.Nil(t, err, "Expected no error")
	require.Equal(t, len(colours)+len(greetings)+len(denied), accounts.totoalMemberCount, "Expected all members in all groups to be processed")
	require.Equal(t, len(aggregate), accounts.invocationCount, "Expected one batch per group")

	// There is no guaranteed ordering in the batch so we cannot check for the resulting string directly
	for k, v := range aggregate {
		require.Equal(t, len(v), accounts.groupCounts[k], "Expected group memeber count to match input for group %s", k)
	}
}

func TestMultipleItemsRejectFirst(t *testing.T) {
	paccounts, processor := newBasicProcessor(nil)
	aaccounts, assessor := rejectSecondAssessor(nil)

	members := []interface{}{"hello", "cruel", "world"}
	target := []interface{}{members[0]}
	target = append(target, members[2:]...)

	batcher := NewBatcher(context.Background(), assessor, processor, false)
	err := multipleMembersSingleGroup(batcher, t.Name(), members...)

	require.Nil(t, err, "Expected no error")

	require.Equal(t, len(members), aaccounts.invocationCount, "Expected exactly one invocation of assessor per member")

	require.Equal(t, aaccounts.invocationCount-1, paccounts.totoalMemberCount, "Expected one rejected members")
	require.Equal(t, 1, paccounts.invocationCount, "Expected exactly one batch to be processed")
	require.NotNil(t, paccounts.groupState, "Expected state to be passed through to processor from assessor")

	// There is no guaranteed ordering in the batch so we cannot check for the resulting string directly
	require.Equal(t, len(members)-1, len(paccounts.result), "Expected all members present except single reject")
	require.Equal(t, target, paccounts.result, "Expected only second member to be dropped")
}
