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
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"

	"github.com/pivotal-cf/go-tracker"
)

var _ = Describe("Tracker Client", func() {
	var (
		server *ghttp.Server
		client *tracker.Client
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		tracker.DefaultURL = server.URL()
		client = tracker.NewClient("api-token")
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("getting information about the current user", func() {
		var statusCode int

		It("works if everything goes to plan", func() {
			statusCode = http.StatusOK

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/v5/me"),
				verifyTrackerToken(),

				ghttp.RespondWith(statusCode, Fixture("me.json")),
			))

			me, err := client.Me()

			Expect(err).NotTo(HaveOccurred())
			Expect(me.Username).To(Equal("vader"))
		})

		It("returns an error if the response is not successful", func() {
			statusCode = http.StatusInternalServerError

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(statusCode, ""),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()
			Expect(err).To(MatchError("request failed (500)"))
		})

		It("returns a helpful error if the token is invalid", func() {
			statusCode = http.StatusUnauthorized

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(statusCode, ""),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()
			Expect(err).To(MatchError("invalid token"))
		})

		It("returns an error if the request fails", func() {
			server.Close()

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("failed to make request"))
			server = ghttp.NewServer()
		})

		It("returns an error if the request can't be created", func() {
			tracker.DefaultURL = "aaaaa)#Q&%*(*"

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("failed to create request"))
		})

		It("returns an error if the response JSON is broken", func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, `{"`),
			))

			client := tracker.NewClient("api-token")
			_, err := client.Me()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("invalid json response"))
		})
	})

	Describe("retrieving a story by ID", func() {
		It("gets one story", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/stories/560"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("story.json")),
				),
			)

			client := tracker.NewClient("api-token")

			story, err := client.Story(560)
			Expect(err).NotTo(HaveOccurred())
			Expect(story.ID).To(Equal(560))
			Expect(story.Name).To(Equal("Tractor beam loses power intermittently"))
		})
	})

	Describe("listing stories", func() {
		It("gets all the stories by default", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
				),
			)

			client := tracker.NewClient("api-token")

			stories, pagination, err := client.InProject(99).Stories(tracker.StoriesQuery{})
			Expect(stories).To(HaveLen(4))
			Expect(pagination).To(BeZero())
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns pagination info allowing the caller to follow through pages themselves", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json"), http.Header{
						"X-Tracker-Pagination-Total":    []string{"1"},
						"X-Tracker-Pagination-Offset":   []string{"2"},
						"X-Tracker-Pagination-Limit":    []string{"3"},
						"X-Tracker-Pagination-Returned": []string{"4"},
					}),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories", "offset=1234"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json"), http.Header{
						"X-Tracker-Pagination-Total":    []string{"5"},
						"X-Tracker-Pagination-Offset":   []string{"6"},
						"X-Tracker-Pagination-Limit":    []string{"7"},
						"X-Tracker-Pagination-Returned": []string{"8"},
					}),
				),
			)

			client := tracker.NewClient("api-token")

			stories, pagination, err := client.InProject(99).Stories(tracker.StoriesQuery{})
			Expect(stories).To(HaveLen(4))
			Expect(pagination).To(Equal(tracker.Pagination{
				Total:    1,
				Offset:   2,
				Limit:    3,
				Returned: 4,
			}))
			Expect(err).NotTo(HaveOccurred())

			stories, pagination, err = client.InProject(99).Stories(tracker.StoriesQuery{Offset: 1234})
			Expect(stories).To(HaveLen(4))
			Expect(pagination).To(Equal(tracker.Pagination{
				Total:    5,
				Offset:   6,
				Limit:    7,
				Returned: 8,
			}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows different queries to be made", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories", "with_state=finished"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("stories.json")),
				),
			)

			client := tracker.NewClient("api-token")

			query := tracker.StoriesQuery{
				State: tracker.StoryStateFinished,
			}
			stories, pagination, err := client.InProject(99).Stories(query)
			Expect(stories).To(HaveLen(4))
			Expect(pagination).To(BeZero())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("listing project memberships", func() {
		It("gets all the project memberships", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/memberships"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("project_memberships.json")),
				),
			)

			client := tracker.NewClient("api-token")

			memberships, err := client.InProject(99).ProjectMemberships()
			Expect(memberships).To(HaveLen(7))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("listing a story's activity", func() {
		It("gets the story's activity", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories/560/activity"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("activities.json")),
				),
			)

			client := tracker.NewClient("api-token")

			activities, err := client.InProject(99).StoryActivity(560, tracker.ActivityQuery{})
			Expect(activities).To(HaveLen(4))
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows different queries to be made", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						"/services/v5/projects/99/stories/560/activity",
						"limit=2&occurred_after=1000000000000&occurred_before=1433091819000&offset=1&since_version=1",
					),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("activities.json")),
				),
			)

			client := tracker.NewClient("api-token")

			query := tracker.ActivityQuery{
				Limit:          2,
				Offset:         1,
				OccurredBefore: 1433091819000,
				OccurredAfter:  1000000000000,
				SinceVersion:   1,
			}
			activities, err := client.InProject(99).StoryActivity(560, query)
			Expect(activities).To(HaveLen(4))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("listing a story's tasks", func() {
		It("gets the story's tasks", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories/560/tasks"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("tasks.json")),
				),
			)

			client := tracker.NewClient("api-token")

			tasks, err := client.InProject(99).StoryTasks(560, tracker.TaskQuery{})
			Expect(tasks).To(HaveLen(3))
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows different queries to be made", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						"/services/v5/projects/99/stories/560/activity",
						"limit=2&occurred_after=1000000000000&occurred_before=1433091819000&offset=1&since_version=1",
					),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("activities.json")),
				),
			)

			client := tracker.NewClient("api-token")

			query := tracker.ActivityQuery{
				Limit:          2,
				Offset:         1,
				OccurredBefore: 1433091819000,
				OccurredAfter:  1000000000000,
				SinceVersion:   1,
			}
			activities, err := client.InProject(99).StoryActivity(560, query)
			Expect(activities).To(HaveLen(4))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("listing a story's comments", func() {
		It("gets the story's comments", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/services/v5/projects/99/stories/560/comments"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("comments.json")),
				),
			)

			client := tracker.NewClient("api-token")

			tasks, err := client.InProject(99).StoryComments(560, tracker.CommentsQuery{})
			Expect(tasks).To(HaveLen(2))
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows different queries to be made", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						"/services/v5/projects/99/stories/560/activity",
						"limit=2&occurred_after=1000000000000&occurred_before=1433091819000&offset=1&since_version=1",
					),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, Fixture("activities.json")),
				),
			)

			client := tracker.NewClient("api-token")

			query := tracker.ActivityQuery{
				Limit:          2,
				Offset:         1,
				OccurredBefore: 1433091819000,
				OccurredAfter:  1000000000000,
				SinceVersion:   1,
			}
			activities, err := client.InProject(99).StoryActivity(560, query)
			Expect(activities).To(HaveLen(4))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("delivering a story", func() {
		It("HTTP PUTs it in its place", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/99/stories/15225523"),
					ghttp.VerifyJSON(`{"current_state":"delivered"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, ""),
				),
			)

			client := tracker.NewClient("api-token")

			err := client.InProject(99).DeliverStory(15225523)
			Expect(err).NotTo(HaveOccurred())
		})

		It("HTTP PUTs it in its place with a comment", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/99/stories/15225523"),
					ghttp.VerifyJSON(`{"current_state":"delivered"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, ""),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/services/v5/projects/99/stories/15225523/comments"),
					ghttp.VerifyJSON(`{"text":"some delive\"}ry comment with tricky text"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusCreated, ""),
				),
			)

			client := tracker.NewClient("api-token")

			comment := `some delive"}ry comment with tricky text`
			err := client.InProject(99).DeliverStoryWithComment(15225523, comment)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("creating a story", func() {
		It("POSTs", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/services/v5/projects/99/stories"),
					ghttp.VerifyJSON(`{"name":"Exhaust ports are ray shielded","blockers":[{"id":5,"description":"something"}]}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, `{
						"id": 1234,
						"project_id": 5678,
						"name": "Exhaust ports are ray shielded",
						"url": "https://some-url.biz/1234",
						"blockers":
						[
						 {"id": 5, "description": "something"}
						]
					}`),
				),
			)

			client := tracker.NewClient("api-token")

			story, err := client.InProject(99).CreateStory(tracker.Story{
				Name: "Exhaust ports are ray shielded",
				Blockers: []tracker.Blocker{
					{
						ID:          5,
						Description: "something",
					},
				},
			})
			Expect(story).To(Equal(tracker.Story{
				ID:        1234,
				ProjectID: 5678,

				Name: "Exhaust ports are ray shielded",

				URL: "https://some-url.biz/1234",
				Blockers: []tracker.Blocker{
					{
						ID:          5,
						Description: "something",
					},
				},
			}))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("udpating a story", func() {
		It("PUTs", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/services/v5/projects/99/stories/1234"),
					ghttp.VerifyJSON(`{"name":"The death star is approaching", "id": 1234 }`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, `{
						"id": 1234,
						"project_id": 5678,
						"name": "The death star is approaching"
					}`),
				),
			)

			client := tracker.NewClient("api-token")

			story, err := client.InProject(99).UpdateStory(tracker.Story{
				Name: "The death star is approaching",
				ID: 1234,
			})

			Expect(story).To(Equal(tracker.Story{
				ID:        1234,
				ProjectID: 5678,
				Name:      "The death star is approaching",
			}))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("deleting a story", func() {
		It("DELETES", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/services/v5/projects/99/stories/1234"),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, ""),
				),
			)
			client := tracker.NewClient("api-token")
			err := client.InProject(99).DeleteStory(1234)
			Expect(err).NotTo(HaveOccurred())
		})
		Context("when the delete is not successful", func() {
			It("returns error saying request failed", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/services/v5/projects/99/stories/1234"),
						verifyTrackerToken(),

						ghttp.RespondWith(http.StatusInternalServerError, ""),
					),
				)
				client := tracker.NewClient("api-token")
				err := client.InProject(99).DeleteStory(1234)
				Expect(err).To(Equal(errors.New("request failed (500)")))
			})
		})
	})

	Describe("creating a task", func() {
		It("POSTs", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/services/v5/projects/99/stories/560/tasks"),
					ghttp.VerifyJSON(`{"description":"some-tracker-task", "position":1}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, `{
					  "kind": "task",
					  "id": 1234,
					  "story_id": 560,
					  "description": "some-tracker-task",
					  "complete": false,
					  "position": 1
					}`),
				),
			)

			client := tracker.NewClient("api-token")

			task, err := client.InProject(99).CreateTask(560, tracker.Task{
				Description: "some-tracker-task",
				Position:    1,
			})

			Expect(task).To(Equal(tracker.Task{
				ID:          1234,
				StoryID:     560,
				Description: "some-tracker-task",
				Position:    1,
				IsComplete:  false,
			}))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("creating a comment", func() {
		It("POSTs", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/services/v5/projects/99/stories/560/comments"),
					ghttp.VerifyJSON(`{"text":"some-tracker-comment"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, `{
			          "kind": "comment",
                      "id": 111,
                      "story_id": 560,
                      "text": "some-tracker-comment",
                      "person_id": 101,
                      "created_at": "2017-03-07T12:00:00Z",
                      "updated_at": "2017-03-07T12:00:00Z"
					}`),
				),
			)

			client := tracker.NewClient("api-token")

			comment, err := client.InProject(99).CreateComment(560, tracker.Comment{
				Text: "some-tracker-comment",
			})

			Expect(comment).To(Equal(tracker.Comment{
				ID:   111,
				Text: "some-tracker-comment",
			}))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("creating a blocker", func() {
		It("POSTs", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/services/v5/projects/99/stories/560/blockers"),
					ghttp.VerifyJSON(`{"description":"some-tracker-blocker"}`),
					verifyTrackerToken(),

					ghttp.RespondWith(http.StatusOK, `{
			          "kind": "blocker",
                      "id": 111,
                      "story_id": 560,
                      "description": "some-tracker-blocker",
                      "person_id": 101,
                      "created_at": "2017-03-07T12:00:00Z",
                      "updated_at": "2017-03-07T12:00:00Z"
					}`),
				),
			)

			client := tracker.NewClient("api-token")

			blocker, err := client.InProject(99).CreateBlocker(560, tracker.Blocker{
				Description: "some-tracker-blocker",
			})

			Expect(blocker).Should(Equal(tracker.Blocker{
				ID:          111,
				Description: "some-tracker-blocker",
			}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

func verifyTrackerToken() http.HandlerFunc {
	headers := http.Header{
		"X-TrackerToken": {"api-token"},
	}

	return ghttp.VerifyHeader(headers)
}
