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

import "time"

type Me Person

type Person struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	ID       int    `json:"id"`
	Email    string `json:"email"`
}

type TimeZone struct {
	Kind      string `json:"kind"`
	OlsonName string `json:"olson_name"`
	Offset    string `json:"offset"`
}

type Project struct {
	ID                           int       `json:"id"`
	Kind                         string    `json:"kind"`
	Name                         string    `json:"name"`
	Version                      int       `json:"version"`
	IterationLength              int       `json:"iteration_length"`
	WeekStartDay                 string    `json:"week_start_day"`
	PointScale                   string    `json:"point_scale"`
	PointScaleIsCustom           bool      `json:"point_scale_is_custom"`
	BugsAndChoresAreEstimatable  bool      `json:"bugs_and_chores_are_estimatable"`
	AutomaticPlanning            bool      `json:"automatic_planning"`
	EnableTasks                  bool      `json:"enable_tasks"`
	TimeZone                     TimeZone  `json:"time_zone"`
	VelocityAveragedOver         int       `json:"velocity_averaged_over"`
	NumberOfDoneIterationsToShow int       `json:"number_of_done_iterations_to_show"`
	HasGoogleDomain              bool      `json:"has_google_domain"`
	EnableIncomingEmails         bool      `json:"enable_incoming_emails"`
	InitialVelocity              int       `json:"initial_velocity"`
	Public                       bool      `json:"public"`
	AtomEnabled                  bool      `json:"atom_enabled"`
	ProjectType                  string    `json:"project_type"`
	StartDate                    string    `json:"start_date"`
	StartTime                    time.Time `json:"start_time"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
	AccountID                    int       `json:"account_id"`
	CurrentIterationNumber       int       `json:"current_iteration_number"`
	EnableFollowing              bool      `json:"enable_following"`
}

type Story struct {
	ID        int `json:"id,omitempty"`
	ProjectID int `json:"project_id,omitempty"`

	URL string `json:"url,omitempty"`

	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        StoryType  `json:"story_type,omitempty"`
	State       StoryState `json:"current_state,omitempty"`
	Estimate    int        `json:"estimate,omitempty"`

	Labels []Label `json:"labels,omitempty"`

	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	Blockers   []Blocker  `json:"blockers,omitempty"`
}

type NewStory struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        StoryType  `json:"story_type,omitempty"`
	State       StoryState `json:"current_state,omitempty"`
	Labels      []Label    `json:"labels,omitempty"`
	Tasks       []Task     `json:"tasks,omitempty"`
	StoryIDs    []int      `json:"story_ids,omitempty"`
	OwnerIDs    []int      `json:"owner_ids,omitempty"`
}

type Task struct {
	ID      int `json:"id,omitempty"`
	StoryID int `json:"story_id,omitempty"`

	Description string `json:"description,omitempty"`
	IsComplete  bool   `json:"complete,omitempty"`
	Position    int    `json:"position,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Comment struct {
	ID   int    `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
}

type Blocker struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

type Label struct {
	Kind      string     `json:"kind,omitempty"`
	ID        int        `json:"id,omitempty"`
	ProjectID int        `json:"project_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Name string `json:"name"`
}

type StoryType string

const (
	StoryTypeFeature = "feature"
	StoryTypeBug     = "bug"
	StoryTypeChore   = "chore"
	StoryTypeRelease = "release"
)

type StoryState string

const (
	StoryStateUnscheduled = "unscheduled"
	StoryStatePlanned     = "planned"
	StoryStateStarted     = "started"
	StoryStateFinished    = "finished"
	StoryStateDelivered   = "delivered"
	StoryStateAccepted    = "accepted"
	StoryStateRejected    = "rejected"
)

type Activity struct {
	Kind             string        `json:"kind"`
	GUID             string        `json:"guid"`
	ProjectVersion   int           `json:"project_version"`
	Message          string        `json:"message"`
	Highlight        string        `json:"highlight"`
	Changes          []interface{} `json:"changes"`
	PrimaryResources []interface{} `json:"primary_resources"`
	Project          interface{}   `json:"project"`
	PerformedBy      interface{}   `json:"performed_by"`
	OccurredAt       time.Time     `json:"occurred_at"`
}

type ProjectMembership struct {
	ID     int `json:"id"`
	Person Person
}
