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

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type connection struct {
	token  string
	client *http.Client
}

func newConnection(token string) connection {
	return connection{
		token:  token,
		client: &http.Client{},
	}
}

type Pagination struct {
	Total    int
	Offset   int
	Limit    int
	Returned int
}

const paginationTotalHeader = "X-Tracker-Pagination-Total"
const paginationOffsetHeader = "X-Tracker-Pagination-Offset"
const paginationLimitHeader = "X-Tracker-Pagination-Limit"
const paginationReturnedHeader = "X-Tracker-Pagination-Returned"

func (c connection) Do(request *http.Request, response interface{}) (Pagination, error) {
	resp, err := c.sendRequest(request)
	if err != nil {
		return Pagination{}, err
	}

	defer resp.Body.Close()

	pagination := Pagination{}

	if val := resp.Header.Get(paginationTotalHeader); len(val) > 0 {
		pagination.Total, err = strconv.Atoi(val)
		if err != nil {
			return Pagination{}, err
		}
	}

	if val := resp.Header.Get(paginationOffsetHeader); len(val) > 0 {
		pagination.Offset, err = strconv.Atoi(val)
		if err != nil {
			return Pagination{}, err
		}
	}

	if val := resp.Header.Get(paginationLimitHeader); len(val) > 0 {
		pagination.Limit, err = strconv.Atoi(val)
		if err != nil {
			return Pagination{}, err
		}
	}

	if val := resp.Header.Get(paginationReturnedHeader); len(val) > 0 {
		pagination.Returned, err = strconv.Atoi(val)
		if err != nil {
			return Pagination{}, err
		}
	}

	if response != nil {
		return pagination, c.decodeResponse(resp, response)
	}

	return pagination, nil
}

func (c connection) CreateRequest(method string, path string, params url.Values) (*http.Request, error) {
	url := DefaultURL + "/services/v5" + path
	query := params.Encode()
	if query != "" {
		url += "?" + query
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	request.Header.Add("X-TrackerToken", c.token)

	return request, nil
}

func (c connection) sendRequest(request *http.Request) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %s", err)
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("invalid token")
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		d, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("request failed (%s): %s", response.Status, string(d))
	}

	return response, nil
}

func (c connection) decodeResponse(response *http.Response, object interface{}) error {
	if err := json.NewDecoder(response.Body).Decode(object); err != nil {
		return fmt.Errorf("invalid json response: %s", err)
	}

	return response.Body.Close()
}
