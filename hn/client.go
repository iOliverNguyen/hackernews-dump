// https://github.com/HackerNews/API
package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiBasePath = "https://hacker-news.firebaseio.com/v0/"

type HNClient struct {
	httpClient *http.Client
}

func NewHNClient() *HNClient {
	return &HNClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HNClient) Query(path string, out any) error {
	apiURL := apiBasePath + path
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %v", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func (c *HNClient) MaxItem() (out int, _ error) {
	err := c.Query("maxitem.json", &out)
	return out, wrapf(err, "query max item")
}

// Item can be nil in case of not found.
func (c *HNClient) GetItem(id int) (out *Item, _ error) {
	err := c.Query(spr("item/%v.json", id), &out)
	return out, wrapf(err, "query item %v", id)
}

func (c *HNClient) GetMultiple(startID int, count int) (out []*Item, _ error) {
	maxID := startID + count
	for id := startID; id < maxID; id++ {
		item, err := c.GetItem(id)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
