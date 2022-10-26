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

func (c *Client) CreateTeam(organization_name string, team_name string, bucket_name string, bucket_provider string) (err error) {

	// Organization ID from Organization Name
	// var params QueryParams
	var organization_id string
	if organization_name != ""{
		params := QueryParams{
			Query: `
				query availableOrg($name: String!) { 
					organization (name: $name){
						id
					}
				}
			`,
			Variables: map[string]interface{}{
				"name": organization_name,
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
				Org struct {
					ID        string `json:"id"`
				} `json:"organization"`
			} `json:"data"`
		}
	
		err = json.Unmarshal(body, &orgResult)
		if err != nil {
			fmt.Println(err)
			return err
		}
		organization_id = orgResult.OrgData.Org.ID
	}

	// Create Team
	params := QueryParams{
		Query: `
		mutation CreateTeam (
			$teamName: String!
			$organizationId: String
			$bucketName: String
			$bucketProvider: String
		){
                createTeam (
					input: {
						teamName: $teamName
						organizationId: $organizationId
						storageBucketInfo: {
							name: $bucketName
							provider: $bucketProvider
						}
					}
				){
                    entity{
						id
						name
					}
				}
			}
        `,
		Variables: map[string]interface{}{
			"teamName":       team_name,
		},
	}
	if organization_name != ""{
		params.Variables["organizationId"] = organization_id
	}
	if bucket_name != ""{
		params.Variables["bucketName"] = bucket_name
	}
	if bucket_provider != ""{
		params.Variables["bucketProvider"] = bucket_provider
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

	var createTeamResult struct {
		CreateTeamData struct {
			CreateTeam struct {
				Entity struct {
					Id   string `json:"id"`
					Name string `json:"name"`
				}
			}
		}
	}

	err = json.Unmarshal(body, &createTeamResult)
	if err != nil {
		return err
	}
	fmt.Println(createTeamResult.CreateTeamData.CreateTeam.Entity.Name)
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

func (c *Client) ReadTeam(name string) (err error) {
	params := QueryParams{
		Query: `query:
            {
                entity (name: $name){
                    id
					name
					createdAt
      				updatedAt
                }
            }
        `,
		Variables: map[string]interface{}{
			"name": name,
		},
	}
	resp, err := c.doQuery(params)

	fmt.Printf("Response: %+v\n", resp)

	if err != nil {
		return err
	}

	return nil
}

// func main() {
// //Testing
// 	defaultTimeout := time.Second * 10
// 	client := NewClient("https://api.wandb.ai", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913", defaultTimeout)

// 	err := client.CreateTeam("xyzw", "team-tmp", "", "")
// 	if err != nil{
// 		fmt.Println(err)
// 	}
// }
