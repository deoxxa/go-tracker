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
package tracker_test

import (
	"github.com/deoxxa/go-tracker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queries", func() {
	queryString := func(query tracker.Query) string {
		return query.Query().Encode()
	}

	Describe("StoriesQuery", func() {
		It("only has date_format by default", func() {
			query := tracker.StoriesQuery{}
			Expect(queryString(query)).To(Equal(""))
		})

		It("can query by story state", func() {
			query := tracker.StoriesQuery{
				State: tracker.StoryStateRejected,
			}
			Expect(queryString(query)).To(Equal("with_state=rejected"))
		})

		It("can query by story labels", func() {
			query := tracker.StoriesQuery{
				Label: "blocked",
			}
			Expect(queryString(query)).To(Equal("with_label=blocked"))
		})

		Describe("query by filter", func() {
			It("handles a single attribute", func() {
				query := tracker.StoriesQuery{
					Filter: []string{
						"owner:dv",
					},
				}
				Expect(queryString(query)).To(Equal("filter=owner%3Adv"))
			})

			It("handles multiple attributes", func() {
				query := tracker.StoriesQuery{
					Filter: []string{
						"owner:dv",
						"state:started",
					},
				}
				Expect(queryString(query)).To(Equal("filter=owner%3Adv+state%3Astarted"))
			})
		})

		It("can limit the numer of results", func() {
			query := tracker.StoriesQuery{
				Limit: 33,
			}
			Expect(queryString(query)).To(Equal("limit=33"))
		})
	})
})
