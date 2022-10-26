package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	host       string
	httpClient *http.Client
	apiKey     string
}

func NewClient(host string, apiKey string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		host:       host,
		httpClient: client,
		apiKey:     apiKey,
	}
}

type QueryParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func (c *Client) doQuery(query QueryParams) (*http.Response, error) {
	jsonBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/graphql", c.host)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	authHeader := fmt.Sprintf("api:%s", c.apiKey)
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64Encode(authHeader)))

	return c.httpClient.Do(request)
}

func base64Encode(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

// Test Function
func (c *Client) queryProject(method string, endpoint string, api_key string) (err error) {
	params := QueryParams{
		Query: `query:
            { 
                projects (entityName: $entityName){
                    pageInfo{
						hasNextPage
					}
                }
            }
        `,
		Variables: map[string]interface{}{
			"entityName": "ibindlish",
		},
	}
	resp, err := c.doQuery(params)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	var projectsResult struct {
		ProjectsData struct {
			Projects struct {
				PageInfo struct {
					HasNextPage bool `json:"hasNextPage"`
				}
			}
		}
	}

	err = json.Unmarshal(body, &projectsResult)
	fmt.Println(projectsResult.ProjectsData.Projects.PageInfo.HasNextPage)

	if err = json.Unmarshal(body, &projectsResult); err != nil {
		fmt.Println(projectsResult.ProjectsData.Projects.PageInfo.HasNextPage)
		return err
	}

	return nil
}

func (c *Client) CreateTeam(method string, endpoint string, api_key string) (err error) {

	// Organization ID from Organization Name
	name := "xyzw"
	params := QueryParams{
		Query: `query: 
            { 
                organization (name: $name){
                    id
					available
                }
            }
        `,
		Variables: map[string]interface{}{
			"name": name,
		},
	}
	resp, err := c.doQuery(params)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var orgResult struct {
		OrgData struct {
			ID        string `json:"id"`
			Available bool   `json:"available"`
		} //`json:"organization (name: "%s") "`
	}
	err = json.Unmarshal(body, &orgResult)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if orgResult.OrgData.Available == false {
		fmt.Println("organization doesn't have any teams left")
	}
	fmt.Println(string(body))
	fmt.Println(orgResult.OrgData)

	// Create Team
	// TODO: input arguments for mutation
	teamName := "tmp-team"
	organizationId := orgResult.OrgData.ID
	params = QueryParams{
		Query: `mutation:
            { 
                createTeam (teamName: $teamName, organizationId: $organizationId){
                    entity{
						id
						name
					}
                }
            }
        `,
		Variables: map[string]interface{}{
			"teamName":       teamName,
			"organizationId": organizationId,
		},
	}
	resp, err = c.doQuery(params)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	var createTeamResult struct {
		CreateTeamData struct {
			Entity struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}
		}
	}

	err = json.Unmarshal(body, &createTeamResult)
	if err != nil {
		return err
	}
	fmt.Println(createTeamResult.CreateTeamData.Entity.Name)

	return nil

}

func (c *Client) DeleteTeam(name string) (err error) {
	params := QueryParams{
		Query: `mutation { deleteTeam(input:{teamName:$teamName}){success}}`,
		Variables: map[string]interface{}{
			"teamName": name,
		},
	}
	resp, err := c.doQuery(params)
	if err != nil {
		return fmt.Errorf("Error deleting team: %s", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting team")
	}

	return nil
}

func main() {
	defaultTimeout := time.Second * 10
	client := NewClient("https://api.wandb.ai", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913", defaultTimeout)

	host := "https://api.wandb.ai"
	err := client.CreateTeam("POST", host+"/graphql", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913")
	if err != nil {
		fmt.Println(err)
	}
}
