package chatops

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dictyBase-docker/github-actions/internal/logger"
	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/urfave/cli"
)

type Payload struct {
	Event github.RepositoryDispatchEvent `json:"event"`
}

type PullRequest struct {
	Links struct {
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Issue struct {
			Href string `json:"href"`
		} `json:"issue"`
		ReviewComment struct {
			Href string `json:"href"`
		} `json:"review_comment"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"_links"`
	ActiveLockReason  interface{}   `json:"active_lock_reason"`
	Additions         int           `json:"additions"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	AuthorAssociation string        `json:"author_association"`
	Base              struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Repo  struct {
			ArchiveURL       string      `json:"archive_url"`
			Archived         bool        `json:"archived"`
			AssigneesURL     string      `json:"assignees_url"`
			BlobsURL         string      `json:"blobs_url"`
			BranchesURL      string      `json:"branches_url"`
			CloneURL         string      `json:"clone_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			CommentsURL      string      `json:"comments_url"`
			CommitsURL       string      `json:"commits_url"`
			CompareURL       string      `json:"compare_url"`
			ContentsURL      string      `json:"contents_url"`
			ContributorsURL  string      `json:"contributors_url"`
			CreatedAt        time.Time   `json:"created_at"`
			DefaultBranch    string      `json:"default_branch"`
			DeploymentsURL   string      `json:"deployments_url"`
			Description      string      `json:"description"`
			Disabled         bool        `json:"disabled"`
			DownloadsURL     string      `json:"downloads_url"`
			EventsURL        string      `json:"events_url"`
			Fork             bool        `json:"fork"`
			Forks            int         `json:"forks"`
			ForksCount       int         `json:"forks_count"`
			ForksURL         string      `json:"forks_url"`
			FullName         string      `json:"full_name"`
			GitCommitsURL    string      `json:"git_commits_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitURL           string      `json:"git_url"`
			HasDownloads     bool        `json:"has_downloads"`
			HasIssues        bool        `json:"has_issues"`
			HasPages         bool        `json:"has_pages"`
			HasProjects      bool        `json:"has_projects"`
			HasWiki          bool        `json:"has_wiki"`
			Homepage         interface{} `json:"homepage"`
			HooksURL         string      `json:"hooks_url"`
			HTMLURL          string      `json:"html_url"`
			ID               int         `json:"id"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			IssuesURL        string      `json:"issues_url"`
			KeysURL          string      `json:"keys_url"`
			LabelsURL        string      `json:"labels_url"`
			Language         string      `json:"language"`
			LanguagesURL     string      `json:"languages_url"`
			License          interface{} `json:"license"`
			MergesURL        string      `json:"merges_url"`
			MilestonesURL    string      `json:"milestones_url"`
			MirrorURL        interface{} `json:"mirror_url"`
			Name             string      `json:"name"`
			NodeID           string      `json:"node_id"`
			NotificationsURL string      `json:"notifications_url"`
			OpenIssues       int         `json:"open_issues"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			Owner            struct {
				AvatarURL         string `json:"avatar_url"`
				EventsURL         string `json:"events_url"`
				FollowersURL      string `json:"followers_url"`
				FollowingURL      string `json:"following_url"`
				GistsURL          string `json:"gists_url"`
				GravatarID        string `json:"gravatar_id"`
				HTMLURL           string `json:"html_url"`
				ID                int    `json:"id"`
				Login             string `json:"login"`
				NodeID            string `json:"node_id"`
				OrganizationsURL  string `json:"organizations_url"`
				ReceivedEventsURL string `json:"received_events_url"`
				ReposURL          string `json:"repos_url"`
				SiteAdmin         bool   `json:"site_admin"`
				StarredURL        string `json:"starred_url"`
				SubscriptionsURL  string `json:"subscriptions_url"`
				Type              string `json:"type"`
				URL               string `json:"url"`
			} `json:"owner"`
			Private         bool      `json:"private"`
			PullsURL        string    `json:"pulls_url"`
			PushedAt        time.Time `json:"pushed_at"`
			ReleasesURL     string    `json:"releases_url"`
			Size            int       `json:"size"`
			SSHURL          string    `json:"ssh_url"`
			StargazersCount int       `json:"stargazers_count"`
			StargazersURL   string    `json:"stargazers_url"`
			StatusesURL     string    `json:"statuses_url"`
			SubscribersURL  string    `json:"subscribers_url"`
			SubscriptionURL string    `json:"subscription_url"`
			SvnURL          string    `json:"svn_url"`
			TagsURL         string    `json:"tags_url"`
			TeamsURL        string    `json:"teams_url"`
			TreesURL        string    `json:"trees_url"`
			UpdatedAt       time.Time `json:"updated_at"`
			URL             string    `json:"url"`
			Watchers        int       `json:"watchers"`
			WatchersCount   int       `json:"watchers_count"`
		} `json:"repo"`
		Sha  string `json:"sha"`
		User struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int    `json:"id"`
			Login             string `json:"login"`
			NodeID            string `json:"node_id"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"user"`
	} `json:"base"`
	Body         string      `json:"body"`
	ChangedFiles int         `json:"changed_files"`
	ClosedAt     interface{} `json:"closed_at"`
	Comments     int         `json:"comments"`
	CommentsURL  string      `json:"comments_url"`
	Commits      int         `json:"commits"`
	CommitsURL   string      `json:"commits_url"`
	CreatedAt    time.Time   `json:"created_at"`
	Deletions    int         `json:"deletions"`
	DiffURL      string      `json:"diff_url"`
	Draft        bool        `json:"draft"`
	Head         struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Repo  struct {
			ArchiveURL       string      `json:"archive_url"`
			Archived         bool        `json:"archived"`
			AssigneesURL     string      `json:"assignees_url"`
			BlobsURL         string      `json:"blobs_url"`
			BranchesURL      string      `json:"branches_url"`
			CloneURL         string      `json:"clone_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			CommentsURL      string      `json:"comments_url"`
			CommitsURL       string      `json:"commits_url"`
			CompareURL       string      `json:"compare_url"`
			ContentsURL      string      `json:"contents_url"`
			ContributorsURL  string      `json:"contributors_url"`
			CreatedAt        time.Time   `json:"created_at"`
			DefaultBranch    string      `json:"default_branch"`
			DeploymentsURL   string      `json:"deployments_url"`
			Description      string      `json:"description"`
			Disabled         bool        `json:"disabled"`
			DownloadsURL     string      `json:"downloads_url"`
			EventsURL        string      `json:"events_url"`
			Fork             bool        `json:"fork"`
			Forks            int         `json:"forks"`
			ForksCount       int         `json:"forks_count"`
			ForksURL         string      `json:"forks_url"`
			FullName         string      `json:"full_name"`
			GitCommitsURL    string      `json:"git_commits_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitURL           string      `json:"git_url"`
			HasDownloads     bool        `json:"has_downloads"`
			HasIssues        bool        `json:"has_issues"`
			HasPages         bool        `json:"has_pages"`
			HasProjects      bool        `json:"has_projects"`
			HasWiki          bool        `json:"has_wiki"`
			Homepage         interface{} `json:"homepage"`
			HooksURL         string      `json:"hooks_url"`
			HTMLURL          string      `json:"html_url"`
			ID               int         `json:"id"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			IssuesURL        string      `json:"issues_url"`
			KeysURL          string      `json:"keys_url"`
			LabelsURL        string      `json:"labels_url"`
			Language         string      `json:"language"`
			LanguagesURL     string      `json:"languages_url"`
			License          interface{} `json:"license"`
			MergesURL        string      `json:"merges_url"`
			MilestonesURL    string      `json:"milestones_url"`
			MirrorURL        interface{} `json:"mirror_url"`
			Name             string      `json:"name"`
			NodeID           string      `json:"node_id"`
			NotificationsURL string      `json:"notifications_url"`
			OpenIssues       int         `json:"open_issues"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			Owner            struct {
				AvatarURL         string `json:"avatar_url"`
				EventsURL         string `json:"events_url"`
				FollowersURL      string `json:"followers_url"`
				FollowingURL      string `json:"following_url"`
				GistsURL          string `json:"gists_url"`
				GravatarID        string `json:"gravatar_id"`
				HTMLURL           string `json:"html_url"`
				ID                int    `json:"id"`
				Login             string `json:"login"`
				NodeID            string `json:"node_id"`
				OrganizationsURL  string `json:"organizations_url"`
				ReceivedEventsURL string `json:"received_events_url"`
				ReposURL          string `json:"repos_url"`
				SiteAdmin         bool   `json:"site_admin"`
				StarredURL        string `json:"starred_url"`
				SubscriptionsURL  string `json:"subscriptions_url"`
				Type              string `json:"type"`
				URL               string `json:"url"`
			} `json:"owner"`
			Private         bool      `json:"private"`
			PullsURL        string    `json:"pulls_url"`
			PushedAt        time.Time `json:"pushed_at"`
			ReleasesURL     string    `json:"releases_url"`
			Size            int       `json:"size"`
			SSHURL          string    `json:"ssh_url"`
			StargazersCount int       `json:"stargazers_count"`
			StargazersURL   string    `json:"stargazers_url"`
			StatusesURL     string    `json:"statuses_url"`
			SubscribersURL  string    `json:"subscribers_url"`
			SubscriptionURL string    `json:"subscription_url"`
			SvnURL          string    `json:"svn_url"`
			TagsURL         string    `json:"tags_url"`
			TeamsURL        string    `json:"teams_url"`
			TreesURL        string    `json:"trees_url"`
			UpdatedAt       time.Time `json:"updated_at"`
			URL             string    `json:"url"`
			Watchers        int       `json:"watchers"`
			WatchersCount   int       `json:"watchers_count"`
		} `json:"repo"`
		Sha  string `json:"sha"`
		User struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int    `json:"id"`
			Login             string `json:"login"`
			NodeID            string `json:"node_id"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"user"`
	} `json:"head"`
	HTMLURL             string        `json:"html_url"`
	ID                  int           `json:"id"`
	IssueURL            string        `json:"issue_url"`
	Labels              []interface{} `json:"labels"`
	Locked              bool          `json:"locked"`
	MaintainerCanModify bool          `json:"maintainer_can_modify"`
	MergeCommitSha      string        `json:"merge_commit_sha"`
	Mergeable           bool          `json:"mergeable"`
	MergeableState      string        `json:"mergeable_state"`
	Merged              bool          `json:"merged"`
	MergedAt            interface{}   `json:"merged_at"`
	MergedBy            interface{}   `json:"merged_by"`
	Milestone           interface{}   `json:"milestone"`
	NodeID              string        `json:"node_id"`
	Number              int           `json:"number"`
	PatchURL            string        `json:"patch_url"`
	Rebaseable          bool          `json:"rebaseable"`
	RequestedReviewers  []interface{} `json:"requested_reviewers"`
	RequestedTeams      []interface{} `json:"requested_teams"`
	ReviewCommentURL    string        `json:"review_comment_url"`
	ReviewComments      int           `json:"review_comments"`
	ReviewCommentsURL   string        `json:"review_comments_url"`
	State               string        `json:"state"`
	StatusesURL         string        `json:"statuses_url"`
	Title               string        `json:"title"`
	UpdatedAt           time.Time     `json:"updated_at"`
	URL                 string        `json:"url"`
	User                struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		NodeID            string `json:"node_id"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"user"`
}

type SlashCommand struct {
	Args Args `json:"args"`
}

type Args struct {
	All     string      `json:"all"`
	Named   NamedArgs   `json:"name"`
	Unnamed UnnamedArgs `json:"unnamed"`
}

type NamedArgs struct {
	Cluster string `json:"cluster"`
	Branch  string `json:"branch"`
	PR      string `json:"pr"`
	Commit  string `json:"commit"`
}

type UnnamedArgs struct {
	All string `json:"all"`
}

type Output struct {
	Cluster  string `json:"cluster"`
	Ref      string `json:"ref"`
	ImageTag string `json:"image_tag"`
}

func GetSlashCommand(data []byte) (*SlashCommand, error) {
	var s string
	o := new(SlashCommand)
	if err := json.Unmarshal(data, &s); err != nil {
		return o, fmt.Errorf("error in decoding json data to string %s", err)
	}
	if err := json.Unmarshal([]byte(s), o); err != nil {
		return o, fmt.Errorf("error in decoding string to structure %s", err)
	}
	return o, nil
}

func GetPullRequestPayload(data []byte) (*PullRequest, error) {
	var s string
	o := new(PullRequest)
	if err := json.Unmarshal(data, &s); err != nil {
		return o, fmt.Errorf("error in decoding json data to string %s", err)
	}
	if err := json.Unmarshal([]byte(s), o); err != nil {
		return o, fmt.Errorf("error in decoding string to structure %s", err)
	}
	return o, nil
}

func ShareChatOpsPayload(c *cli.Context) error {
	r, err := os.Open(c.String("payload-file"))
	if err != nil {
		return fmt.Errorf("error in reading content from file %s", err)
	}
	defer r.Close()
	d := &Payload{}
	if err := json.NewDecoder(r).Decode(d); err != nil {
		return fmt.Errorf("error in decoding json %s", err)
	}
	p, err := GetSlashCommand(d.Event.ClientPayload)
	if err != nil {
		return err
	}
	a := githubactions.New()
	log := logger.GetLogger(c)
	pr := containsPR(p.Args.Unnamed.All)
	// if command is for a PR
	if pr {
		u, err := parsePullRequestData(p.Args.Named, d)
		if err != nil {
			return err
		}
		a.SetOutput("cluster", c.String("cluster"))
		a.SetOutput("ref", u.Ref)
		a.SetOutput("image_tag", u.ImageTag)
	}
	log.Info("added all keys to the output")
	return nil
}

/**
Examples:

/deploy cluster=erickube commit=xyz
/deploy cluster=erickube pr=1 commit=xyz
/deploy cluster=erickube pr
/deploy cluster=erickube pr commit=xyz
/deploy cluster=erickube branch=develop
*/

func containsPR(s string) bool {
	return strings.Contains(s, "pr")
}

func parsePullRequestData(args NamedArgs, p *Payload) (*Output, error) {
	o := &Output{}
	pr, err := GetPullRequestPayload(p.Event.ClientPayload)
	if err != nil {
		return o, err
	}
	if args.Commit != "" {
		o.ImageTag = pr.Head.Sha[0:6]
		o.Ref = pr.Head.Sha
	} else {
		o.ImageTag = fmt.Sprintf("pr-%d-%s", pr.Number, args.Commit[0:6])
		o.Ref = args.Commit
	}
	return o, nil
}
