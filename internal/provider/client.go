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

func (c *Client) doQuery(method, endpoint string, query *bytes.Buffer) (*http.Response, error) {
	request, err := http.NewRequest(method, endpoint, query)
	if err != nil {
	 return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Test Function
func (c *Client) queryProject(method string, endpoint string) (err error) {
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
	resp, err := c.doQuery(method, endpoint, bytes.NewBuffer(jsonValue))

	if err != nil{
		return err
	}
	defer resp.Body.Close()

   	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	type rowStruct []interface{}
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

func main() {
	defaultTimeout := time.Second * 10
	client := NewClient("https://api.wandb.ai", "19f7df3fa4db872d5e4cea31ed8076e6b1ff5913", defaultTimeout)

	host := "https://api.wandb.ai"
	err := client.queryProject("POST", host + "/graphql")
	if err != nil{
		fmt.Println(err)
	}
}

