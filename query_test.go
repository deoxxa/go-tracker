package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xoebus/go-tracker"
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
