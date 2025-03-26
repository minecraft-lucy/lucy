/*
Copyright 2024 4rcadia

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

// Package datatypes are types to store data from APIs and files.
//
// TODO: For web APIs, their types will be moved to where they are used.
package datatypes

import (
	"time"
)

// ModrinthProject is a struct that represents a Modrinth project, the basic
// form of any item on Modrinth.
//
// API Example:
//   - https://api.modrinth.com/v2/project/P7dR8mSH
//     (Fabric API)
//   - https://api.modrinth.com/v2/project/1IjD5062
//     (Continuity)
type ModrinthProject struct {
	ClientSide       string      `json:"client_side"`
	ServerSide       string      `json:"server_side"`
	GameVersions     []string    `json:"game_versions"`
	Id               string      `json:"id"`
	Slug             string      `json:"slug"`
	ProjectType      string      `json:"project_type"`
	Team             string      `json:"team"`
	Organization     interface{} `json:"organization"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Body             string      `json:"body"`
	BodyUrl          interface{} `json:"body_url"`
	Published        time.Time   `json:"published"`
	Updated          time.Time   `json:"updated"`
	Approved         time.Time   `json:"approved"`
	Queued           interface{} `json:"queued"`
	Status           string      `json:"status"`
	RequestedStatus  interface{} `json:"requested_status"`
	ModeratorMessage interface{} `json:"moderator_message"`
	License          struct {
		Id   string      `json:"id"`
		Name string      `json:"name"`
		Url  interface{} `json:"url"`
	} `json:"license"`
	Downloads            int           `json:"downloads"`
	Followers            int           `json:"followers"`
	Categories           []string      `json:"categories"`
	AdditionalCategories []interface{} `json:"additional_categories"`
	Loaders              []string      `json:"loaders"`
	Versions             []string      `json:"versions"`
	IconUrl              string        `json:"icon_url"`
	IssuesUrl            string        `json:"issues_url"`
	SourceUrl            string        `json:"source_url"`
	WikiUrl              string        `json:"wiki_url"`
	DiscordUrl           string        `json:"discord_url"`
	DonationUrls         []struct {
		Id       string `json:"id"`
		Platform string `json:"platform"`
		Url      string `json:"url"`
	} `json:"donation_urls"`
	Gallery            []interface{} `json:"gallery"`
	Color              int           `json:"color"`
	ThreadId           string        `json:"thread_id"`
	MonetizationStatus string        `json:"monetization_status"`
}

// ModrinthSearchResults
//
// Docs
// https://docs.modrinth.com/api/operations/searchprojects/
//
// Example
// https://api.modrinth.com/v2/search?query=carpet&limit=100&index=relevance&facets=%5B%5B%22server_side:required%22,%22server_side:optional%22%5D%5D
type ModrinthSearchResults struct {
	Hits []struct {
		ProjectId         string    `json:"project_id"`
		ProjectType       string    `json:"project_type"`
		Slug              string    `json:"slug"`
		Author            string    `json:"author"`
		Title             string    `json:"title"`
		Description       string    `json:"description"`
		Categories        []string  `json:"categories"`
		DisplayCategories []string  `json:"display_categories"`
		Versions          []string  `json:"versions"`
		Downloads         int       `json:"downloads"`
		Follows           int       `json:"follows"`
		IconUrl           string    `json:"icon_url"`
		DateCreated       time.Time `json:"date_created"`
		DateModified      time.Time `json:"date_modified"`
		LatestVersion     string    `json:"latest_version"`
		License           string    `json:"license"`
		ClientSide        string    `json:"client_side"`
		ServerSide        string    `json:"server_side"`
		Gallery           []string  `json:"gallery"`
		FeaturedGallery   *string   `json:"featured_gallery"`
		Color             *int      `json:"color"`
	} `json:"hits"`
	Offset    int `json:"offset"`
	Limit     int `json:"limit"`
	TotalHits int `json:"total_hits"`
}

type ModrinthVersionFile struct {
	Hashes struct {
		Sha1   string `json:"sha1"`
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Primary  bool   `json:"primary"`
	Size     int    `json:"size"`
	FileType string `json:"file_type"`
}

// ModrinthVersion
//
// Docs
// https://docs.modrinth.com/api/operations/getversion/
//
// Example
// https://api.modrinth.com/v2/version/F7LVluUL
type ModrinthVersion struct {
	GameVersions    []string                      `json:"game_versions"`
	Loaders         []string                      `json:"loaders"`
	Id              string                        `json:"id"`
	ProjectId       string                        `json:"project_id"`
	AuthorId        string                        `json:"author_id"`
	Featured        bool                          `json:"featured"`
	Name            string                        `json:"name"`
	VersionNumber   string                        `json:"version_number"`
	Changelog       string                        `json:"changelog"`
	ChangelogUrl    interface{}                   `json:"changelog_url"`
	DatePublished   time.Time                     `json:"date_published"`
	Downloads       int                           `json:"downloads"`
	VersionType     string                        `json:"version_type"`
	Status          string                        `json:"status"`
	RequestedStatus interface{}                   `json:"requested_status"`
	Files           []ModrinthVersionFile         `json:"files"`
	Dependencies    []ModrinthVersionDependencies `json:"dependencies"`
}

type ModrinthVersionDependencyType string

const (
	ModrinthVersionDependencyTypeRequired     ModrinthVersionDependencyType = "required"
	ModrinthVersionDependencyTypeOptional     ModrinthVersionDependencyType = "optional"
	ModrinthVersionDependencyTypeIncompatible ModrinthVersionDependencyType = "incompatible"
	ModrinthVersionDependencyTypeEmbedded     ModrinthVersionDependencyType = "embedded"
)

type ModrinthVersionDependencies struct {
	VersionId      string                        `json:"version_id"`
	ProjectId      string                        `json:"project_id"`
	FileName       string                        `json:"file_name"`
	DependencyType ModrinthVersionDependencyType `json:"dependency_type"`
}

// ModrinthMember
//
// Docs
// https://docs.modrinth.com/api/operations/getprojectteammembers/
//
// Example
// https://api.modrinth.com/v2/project/carpet/members
type ModrinthMember struct {
	Role         string       `json:"role"`
	TeamId       string       `json:"team_id"`
	User         ModrinthUser `json:"user"`
	Permissions  interface{}  `json:"permissions"`
	Accepted     bool         `json:"accepted"`
	PayoutsSplit interface{}  `json:"payouts_split"`
	Ordering     int          `json:"ordering"`
}

// ModrinthUser
//
// # The url can either be a id or username
//
// Example
// https://modrinth.com/user/gnembon
type ModrinthUser struct {
	Id                  string    `json:"id"`
	Username            string    `json:"username"`
	AvatarUrl           string    `json:"avatar_url"`
	Bio                 string    `json:"bio"`
	Created             time.Time `json:"created"`
	Role                string    `json:"role"`
	Badges              int       `json:"badges"`
	AuthProviders       string    `json:"auth_providers"`
	Email               string    `json:"email"`
	EmailVerified       bool      `json:"email_verified"`
	HasPassword         bool      `json:"has_password"`
	HasTotp             bool      `json:"has_totp"`
	PayoutData          string    `json:"payout_data"`
	StripeCustomerId    string    `json:"stripe_customer_id"`
	AllowFriendRequests bool      `json:"allow_friend_requests"`
	GithubId            string    `json:"github_id"`
}

type ModrinthProjectDependencies struct {
	Projects []ModrinthProject `json:"projects"`
	Versions []ModrinthVersion `json:"versions"`
}
