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
	"encoding/json"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/go-tracker"
)

var _ = Describe("Me", func() {
	It("has attributes", func() {
		var me tracker.Me
		reader := strings.NewReader(Fixture("me.json"))
		err := json.NewDecoder(reader).Decode(&me)
		Expect(err).NotTo(HaveOccurred())

		Expect(me.Username).To(Equal("vader"))
		Expect(me.Name).To(Equal("Darth Vader"))
		Expect(me.Initials).To(Equal("DV"))
		Expect(me.ID).To(Equal(101))
		Expect(me.Email).To(Equal("vader@deathstar.mil"))
	})
})

var _ = Describe("Story", func() {
	It("has attributes", func() {
		var stories []tracker.Story
		reader := strings.NewReader(Fixture("stories.json"))
		err := json.NewDecoder(reader).Decode(&stories)
		Expect(err).NotTo(HaveOccurred())
		story := stories[0]

		Expect(story.ID).To(Equal(560))
		Expect(story.Name).To(Equal("Tractor beam loses power intermittently"))
		Expect(story.Labels).To(Equal([]tracker.Label{
			{ID: 10, ProjectID: 99, Name: "some-label"},
			{ID: 11, ProjectID: 99, Name: "some-other-label"},
		}))
		Expect(story.Estimate).To(Equal(3))
		Expect(*story.CreatedAt).To(Equal(time.Date(2015, 07, 20, 22, 50, 50, 0, time.UTC)))
		Expect(*story.UpdatedAt).To(Equal(time.Date(2015, 07, 20, 22, 51, 50, 0, time.UTC)))
		Expect(*story.AcceptedAt).To(Equal(time.Date(2015, 07, 20, 22, 52, 50, 0, time.UTC)))
		Expect(story.Blockers).To(Equal([]tracker.Blocker{
			{ID: 12, Description: "some blocker"},
			{ID: 13, Description: "some other blocker"},
		}))
	})
})

var _ = Describe("Task", func() {
	It("has attributes", func() {
		var tasks []tracker.Task
		reader := strings.NewReader(Fixture("tasks.json"))
		err := json.NewDecoder(reader).Decode(&tasks)
		Expect(err).NotTo(HaveOccurred())
		task := tasks[0]

		Expect(task.ID).To(Equal(52167427))
		Expect(task.StoryID).To(Equal(137910061))
		Expect(task.Description).To(Equal("some-task-description"))
		Expect(task.IsComplete).To(BeTrue())
		Expect(task.Position).To(Equal(1))
	})
})

var _ = Describe("Activity", func() {
	It("has attributes", func() {
		var activities []tracker.Activity
		reader := strings.NewReader(Fixture("activities.json"))
		err := json.NewDecoder(reader).Decode(&activities)
		Expect(err).NotTo(HaveOccurred())
		activity := activities[0]

		Expect(activity.GUID).To(Equal("99_45"))
		Expect(activity.Message).To(Equal("Darth Vader started this feature"))
	})
})

var _ = Describe("Project Memberships", func() {
	It("has attributes", func() {
		var projectMemberships []tracker.ProjectMembership
		reader := strings.NewReader(Fixture("project_memberships.json"))
		err := json.NewDecoder(reader).Decode(&projectMemberships)
		Expect(err).NotTo(HaveOccurred())

		membership := projectMemberships[0]
		Expect(membership.ID).To(Equal(100))
		Expect(membership.Person.ID).To(Equal(100))
		Expect(membership.Person.Name).To(Equal("Emperor Palpatine"))
		Expect(membership.Person.Email).To(Equal("emperor@galacticrepublic.gov"))
		Expect(membership.Person.Initials).To(Equal("EP"))
		Expect(membership.Person.Username).To(Equal("palpatine"))
	})
})
