// Copyright 2016 Christopher Brown. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tracker

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Query interface {
	Query() url.Values
}

type StoriesQuery struct {
	State  StoryState
	Label  string
	Filter []string

	AcceptedBefore time.Time
	AcceptedAfter  time.Time
	CreatedBefore  time.Time
	CreatedAfter   time.Time
	UpdatedBefore  time.Time
	UpdatedAfter   time.Time

	Limit  int
	Offset int
}

func (query StoriesQuery) Query() url.Values {
	params := url.Values{}

	if query.State != "" {
		params.Set("with_state", string(query.State))
	}

	if query.Label != "" {
		params.Set("with_label", query.Label)
	}

	if len(query.Filter) != 0 {
		params.Set("filter", strings.Join(query.Filter, " "))
	}

	if !query.AcceptedBefore.IsZero() {
		params.Set("accepted_before", query.AcceptedBefore.Format(time.RFC3339))
	}
	if !query.AcceptedAfter.IsZero() {
		params.Set("accepted_after", query.AcceptedAfter.Format(time.RFC3339))
	}
	if !query.CreatedBefore.IsZero() {
		params.Set("created_before", query.CreatedBefore.Format(time.RFC3339))
	}
	if !query.CreatedAfter.IsZero() {
		params.Set("created_after", query.CreatedAfter.Format(time.RFC3339))
	}
	if !query.UpdatedBefore.IsZero() {
		params.Set("updated_before", query.UpdatedBefore.Format(time.RFC3339))
	}
	if !query.UpdatedAfter.IsZero() {
		params.Set("updated_after", query.UpdatedAfter.Format(time.RFC3339))
	}

	if query.Limit != 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}

	if query.Offset != 0 {
		params.Set("offset", fmt.Sprintf("%d", query.Offset))
	}

	return params
}

type ActivityQuery struct {
	Limit          int
	Offset         int
	OccurredBefore int64
	OccurredAfter  int64
	SinceVersion   int
}

func (query ActivityQuery) Query() url.Values {
	params := url.Values{}

	if query.Limit != 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}

	if query.Offset != 0 {
		params.Set("offset", fmt.Sprintf("%d", query.Offset))
	}

	if query.OccurredBefore != 0 {
		params.Set("occurred_before", fmt.Sprintf("%d", query.OccurredBefore))
	}

	if query.OccurredAfter != 0 {
		params.Set("occurred_after", fmt.Sprintf("%d", query.OccurredAfter))
	}

	if query.SinceVersion != 0 {
		params.Set("since_version", fmt.Sprintf("%d", query.SinceVersion))
	}

	return params
}

type TaskQuery struct {}

func (query TaskQuery) Query() url.Values {
	return url.Values{}
}