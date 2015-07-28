package githubexample

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const redirectURI = "http://localhost:8080/oauth2callback"
const githubAPIURL = "https://api.github.com"

type GithubAPI struct {
	ctx         context.Context
	accessToken string
	username    string
}

type CommitStats struct {
	Additions, Deletions int
}

func NewGithubAPI(ctx context.Context) *GithubAPI {
	return &GithubAPI{
		ctx: ctx,
	}
}

func (api *GithubAPI) getUsername() (string, error) {
	client := urlfetch.Client(api.ctx)
	response, err := client.Get(
		githubAPIURL + "/user?access_token=" + api.accessToken)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var data struct {
		Login string
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Login, nil
}

func (api *GithubAPI) getAccessToken(state, code string) (string, error) {
	values := make(url.Values)
	values.Add("client_id", "0ccd33716940f347065e")
	values.Add("client_secret", "4c21ab338de0449ae13019de25629a7b85e08641")
	values.Add("code", code)
	values.Add("state", state)

	client := urlfetch.Client(api.ctx)
	response, err := client.PostForm("https://github.com/login/oauth/access_token", values)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	values, err = url.ParseQuery(string(bs))
	if err != nil {
		return "", err
	}
	return values.Get("access_token"), nil
}

func (api *GithubAPI) getCommitSummaryStats(since time.Time) (CommitStats, error) {
	var stats CommitStats
	// list organizations: GET /user/orgs
	organizations, err := api.getOrganizations()
	if err != nil {
		return stats, err
	}
	for _, organization := range organizations {
		// list repositories: GET /orgs/:org/repos
		repositories, err := api.getRepositories(organization)
		if err != nil {
			return stats, err
		}
		log.Infof(api.ctx, "REPOSITORIES:%v", repositories)
		for _, repository := range repositories {
			// list all commits for repository: GET /repos/:owner/:repo/commits
			shas, err := api.getUserCommitShas(organization, repository, since)
			if err != nil {
				return stats, err
			}
			log.Infof(api.ctx, "SHAS:%v", shas)
			for _, sha := range shas {
				// get a single commit: GET /repos/:owner/:repo/commits/:sha
				cs, err := api.getCommitStats(organization, repository, sha)
				if err != nil {
					return stats, err
				}
				stats.Additions += cs.Additions
				stats.Deletions += cs.Deletions
			}
		}
	}

	return stats, nil
}

func (api *GithubAPI) getOrganizations() ([]string, error) {
	endpoint := "/user/orgs"
	var data []struct {
		Login string `json:"login"`
	}
	err := api.makeAPIRequest(endpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(data))
	for i, v := range data {
		names[i] = v.Login
	}

	return names, nil
}

func (api *GithubAPI) getRepositories(organization string) ([]string, error) {
	// GET /orgs/:org/repos
	var data []struct {
		Name string `json:"name"`
	}
	err := api.makeAPIRequest("/orgs/"+organization+"/repos", nil, &data)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(data))
	for i, v := range data {
		names[i] = v.Name
	}
	return names, nil
}

func (api *GithubAPI) getUserCommitShas(organization, repository string, since time.Time) ([]string, error) {
	values := make(url.Values)
	values.Add("author", api.username)
	values.Add("since", since.Format(time.RFC3339))
	// GET /repos/:owner/:repo/commits
	endpoint := "/repos/" + organization + "/" + repository + "/commits"
	var data []struct {
		SHA string `json:"sha"`
	}
	err := api.makeAPIRequest(endpoint, values, &data)
	if err != nil {
		return nil, err
	}
	shas := make([]string, len(data))
	for i, v := range data {
		shas[i] = v.SHA
	}
	return shas, nil
}

func (api *GithubAPI) getCommitStats(organization, repository, sha string) (CommitStats, error) {
	var stats CommitStats
	endpoint := "/repos/" + organization + "/" + repository + "/commits/" + sha
	var data struct {
		Stats struct {
			Additions int `json:"additions"`
			Deletions int `json:"deletions"`
		} `json:"stats"`
	}
	err := api.makeAPIRequest(endpoint, nil, &data)
	if err != nil {
		return stats, err
	}
	stats.Additions = data.Stats.Additions
	stats.Deletions = data.Stats.Deletions
	return stats, nil
}

func (api *GithubAPI) makeAPIRequest(endpoint string, values url.Values, dst interface{}) error {
	client := urlfetch.Client(api.ctx)
	if values == nil {
		values = make(url.Values)
	}
	values.Add("access_token", api.accessToken)
	// GET /user/orgs
	response, err := client.Get(githubAPIURL + endpoint + "?" + values.Encode())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bs, _ := ioutil.ReadAll(response.Body)
	log.Infof(api.ctx, "GET %s RESPONSE: %s", endpoint, string(bs))
	return json.Unmarshal(bs, dst)

}
