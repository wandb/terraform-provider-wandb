package main

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type (
	Client struct {
	 host       string
	 httpClient *http.Client
	 apiKey     string
	}
)

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

func (c *Client) doQuery(method, endpoint string, api_key string, query *bytes.Buffer) (*http.Response, error) {
	request, err := http.NewRequest(method, endpoint, query)
	if err != nil {
	 return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", api_key)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Test Function
func (c *Client) queryProject(method string, endpoint string, api_key string) (err error) {
	jsonData := map[string]string{
        "query": `
            { 
                projects (entityName: "ibindlish"){
                    pageInfo{
						hasNextPage
					}
                }
            }
        `,
    }
    jsonValue, _ := json.Marshal(jsonData)
	resp, err := c.doQuery(method, endpoint, api_key, bytes.NewBuffer(jsonValue))

	if err != nil{
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
			Projects struct        {
				PageInfo struct     {
					HasNextPage bool	`json:"hasNextPage"`
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
	jsonDataOrgID := map[string]string{
        "query":fmt.Sprintf(`
            { 
                organization (name: "%s"){
                    id
					available
                }
            }
        `, 
		name,
	),
    }
    jsonValue, _ := json.Marshal(jsonDataOrgID)
	resp, err := c.doQuery(method, endpoint, api_key, bytes.NewBuffer(jsonValue))

	if err != nil{
		return err
	}
	defer resp.Body.Close()

   	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var orgResult struct {
		OrgData struct {
			Org struct {
			Available bool `json:"available"`
			ID string `json:"id"`
		}	`json:"organization"`
	} `json:"data"`
	}
	
	err = json.Unmarshal(body, &orgResult)
	if err != nil{
		fmt.Println(err)
		return err
	}

	if orgResult.OrgData.Org.Available == false {
		fmt.Println("organization doesn't have any teams left")
	}

	// Create Team
	// TODO: input arguments for mutation
	teamName := "tmp-team2"
	organizationId := orgResult.OrgData.Org.ID
	jsonData := map[string]string{
        "query": fmt.Sprintf(` 
			mutation {
				createTeam (
					input: {
						teamName: "%s"
						organizationId: "%s"
					}
				){
					entity{
						id
						name
					}
				}
			}
        `,
		teamName,
		organizationId,
	),
    }
	jsonValue, _ = json.Marshal(jsonData)
	resp, err = c.doQuery(method, endpoint, api_key, bytes.NewBuffer(jsonValue))

	if err != nil{
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
				Id string	`json:"id"`
				Name string	`json:"name"`
			}
		}
	}

	err = json.Unmarshal(body, &createTeamResult)
	if err != nil{
		return err
	}
	fmt.Println(createTeamResult.CreateTeamData.Entity.Name)

	return nil

}


func main() {
	defaultTimeout := time.Second * 10
	client := NewClient("https://api.wandb.ai", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913", defaultTimeout)

	host := "https://api.wandb.ai"
	err := client.CreateTeam("POST", host + "/graphql", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913")
	if err != nil{
		fmt.Println(err)
	}
}


