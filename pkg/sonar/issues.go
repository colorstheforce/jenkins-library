package sonar

import (
	"net/http"
	"strings"

	sonargo "github.com/magicsong/sonargo/sonar"
	"github.com/pkg/errors"

	"github.com/SAP/jenkins-library/pkg/log"
)

type IssuesSearchOption struct {
	Branch       string `url:"branch,omitempty"`       // Description:"Branch key"
	Organization string `url:"organization,omitempty"` // Description:"Organization key"
	PullRequest  string `url:"pullRequest,omitempty"`  // Description:"Pull request id"
	// copied from https://github.com/magicsong/sonargo/blob/103eda7abc20bd192a064b6eb94ba26329e339f1/sonar/issues_service.go#L311
	AdditionalFields   string `url:"additionalFields,omitempty"`   // Description:"Comma-separated list of the optional fields to be returned in response. Action plans are dropped in 5.5, it is not returned in the response.",ExampleValue:""
	Asc                string `url:"asc,omitempty"`                // Description:"Ascending sort",ExampleValue:""
	Assigned           string `url:"assigned,omitempty"`           // Description:"To retrieve assigned or unassigned issues",ExampleValue:""
	Assignees          string `url:"assignees,omitempty"`          // Description:"Comma-separated list of assignee logins. The value '__me__' can be used as a placeholder for user who performs the request",ExampleValue:"admin,usera,__me__"
	Authors            string `url:"authors,omitempty"`            // Description:"Comma-separated list of SCM accounts",ExampleValue:"torvalds@linux-foundation.org"
	ComponentKeys      string `url:"componentKeys,omitempty"`      // Description:"Comma-separated list of component keys. Retrieve issues associated to a specific list of components (and all its descendants). A component can be a portfolio, project, module, directory or file.",ExampleValue:"my_project"
	ComponentRootUuids string `url:"componentRootUuids,omitempty"` // Description:"If used, will have the same meaning as componentUuids AND onComponentOnly=false.",ExampleValue:""
	ComponentRoots     string `url:"componentRoots,omitempty"`     // Description:"If used, will have the same meaning as componentKeys AND onComponentOnly=false.",ExampleValue:""
	ComponentUuids     string `url:"componentUuids,omitempty"`     // Description:"To retrieve issues associated to a specific list of components their sub-components (comma-separated list of component IDs). This parameter is mostly used by the Issues page, please prefer usage of the componentKeys parameter. A component can be a project, module, directory or file.",ExampleValue:"584a89f2-8037-4f7b-b82c-8b45d2d63fb2"
	Components         string `url:"components,omitempty"`         // Description:"If used, will have the same meaning as componentKeys AND onComponentOnly=true.",ExampleValue:""
	CreatedAfter       string `url:"createdAfter,omitempty"`       // Description:"To retrieve issues created after the given date (inclusive). <br>Either a date (server timezone) or datetime can be provided. <br>If this parameter is set, createdSince must not be set",ExampleValue:"2017-10-19 or 2017-10-19T13:00:00+0200"
	CreatedAt          string `url:"createdAt,omitempty"`          // Description:"Datetime to retrieve issues created during a specific analysis",ExampleValue:"2017-10-19T13:00:00+0200"
	CreatedBefore      string `url:"createdBefore,omitempty"`      // Description:"To retrieve issues created before the given date (inclusive). <br>Either a date (server timezone) or datetime can be provided.",ExampleValue:"2017-10-19 or 2017-10-19T13:00:00+0200"
	CreatedInLast      string `url:"createdInLast,omitempty"`      // Description:"To retrieve issues created during a time span before the current time (exclusive). Accepted units are 'y' for year, 'm' for month, 'w' for week and 'd' for day. If this parameter is set, createdAfter must not be set",ExampleValue:"1m2w (1 month 2 weeks)"
	Issues             string `url:"issues,omitempty"`             // Description:"Comma-separated list of issue keys",ExampleValue:"5bccd6e8-f525-43a2-8d76-fcb13dde79ef"
	Languages          string `url:"languages,omitempty"`          // Description:"Comma-separated list of languages. Available since 4.4",ExampleValue:"java,js"
	P                  string `url:"p,omitempty"`                  // Description:"1-based page number",ExampleValue:"42"
	Ps                 string `url:"ps,omitempty"`                 // Description:"Page size. Must be greater than 0 and less or equal than 500",ExampleValue:"20"
	Resolutions        string `url:"resolutions,omitempty"`        // Description:"Comma-separated list of resolutions",ExampleValue:"FIXED,REMOVED"
	Resolved           string `url:"resolved,omitempty"`           // Description:"To match resolved or unresolved issues",ExampleValue:""
	Rules              string `url:"rules,omitempty"`              // Description:"Comma-separated list of coding rule keys. Format is &lt;repository&gt;:&lt;rule&gt;",ExampleValue:"squid:AvoidCycles"
	S                  string `url:"s,omitempty"`                  // Description:"Sort field",ExampleValue:""
	Severities         string `url:"severities,omitempty"`         // Description:"Comma-separated list of severities",ExampleValue:"BLOCKER,CRITICAL"
	SinceLeakPeriod    string `url:"sinceLeakPeriod,omitempty"`    // Description:"To retrieve issues created since the leak period.<br>If this parameter is set to a truthy value, createdAfter must not be set and one component id or key must be provided.",ExampleValue:""
	Statuses           string `url:"statuses,omitempty"`           // Description:"Comma-separated list of statuses",ExampleValue:"OPEN,REOPENED"
	Tags               string `url:"tags,omitempty"`               // Description:"Comma-separated list of tags.",ExampleValue:"security,convention"
	Types              string `url:"types,omitempty"`              // Description:"Comma-separated list of types.",ExampleValue:"CODE_SMELL,BUG"
}

func NewIssuesService(host, token, project, organization, branch, pullRequest string, client Sender) *IssueService {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	// Make sure the given URL end with a api/
	if !strings.HasSuffix(host, "api/") {
		host += "api/"
	}
	return &IssueService{
		Host:         host,
		Token:        token,
		Organization: organization,
		Project:      project,
		Branch:       branch,
		PullRequest:  pullRequest,
		apiClient:    NewBasicAuthClient(token, "", host, client),
	}
}

type IssueService struct {
	Host         string
	Token        string
	Organization string
	Project      string
	Branch       string
	PullRequest  string
	apiClient    *Requester
}

type issueSeverity string

func (s issueSeverity) ToString() string {
	return string(s)
}

const (
	blocker  issueSeverity = "BLOCKER"
	critical issueSeverity = "CRITICAL"
	major    issueSeverity = "MAJOR"
	minor    issueSeverity = "MINOR"
	info     issueSeverity = "INFO"
)

func (service *IssueService) GetNumberOfBlockerIssues() (int, error) {
	return service.getIssueCount(blocker)
}

func (service *IssueService) GetNumberOfCriticalIssues() (int, error) {
	return service.getIssueCount(critical)
}

func (service *IssueService) GetNumberOfMajorIssues() (int, error) {
	return service.getIssueCount(major)
}

func (service *IssueService) GetNumberOfMinorIssues() (int, error) {
	return service.getIssueCount(minor)
}

func (service *IssueService) GetNumberOfInfoIssues() (int, error) {
	return service.getIssueCount(info)
}

func (service *IssueService) getIssueCount(severity issueSeverity) (int, error) {
	log.Entry().Debugf("using api client for '%s'", service.Host)
	options := &IssuesSearchOption{
		ComponentKeys: service.Project,
		Severities:    severity.ToString(),
		Resolved:      "false",
		Ps:            "1",
	}
	if len(service.Branch) > 0 {
		options.Branch = service.Branch
	}
	if len(service.Organization) > 0 {
		options.Organization = service.Organization
	}
	if len(service.PullRequest) > 0 {
		options.PullRequest = service.PullRequest
	}
	result, _, err := service.SearchIssues(options)
	if err != nil {
		return -1, errors.Wrapf(err, "failed to fetch the numer of '%s' issues", severity)
	}
	return result.Total, nil
}

// SearchIssues Search for issues.<br>At most one of the following parameters can be provided at the same time: componentKeys, componentUuids, components, componentRootUuids, componentRoots.<br>Requires the 'Browse' permission on the specified project(s).
func (service *IssueService) SearchIssues(options *IssuesSearchOption) (result *sonargo.IssuesSearchObject, response *http.Response, err error) {
	request, err := service.apiClient.create("GET", "issues/search", options)
	if err != nil {
		return
	}
	// use custom HTTP client to send request
	response, err = service.apiClient.send(request)
	if err != nil {
		return
	}
	// reuse response verrification from sonargo
	err = sonargo.CheckResponse(response)
	if err != nil {
		return
	}
	// decode JSON response
	result = new(sonargo.IssuesSearchObject)
	err = service.apiClient.decode(response, result)
	if err != nil {
		return nil, response, err
	}
	return
}
