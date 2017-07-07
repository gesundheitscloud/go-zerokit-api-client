//BSD 3-Clause License
//
//Copyright (c) 2017, Hasso-Plattner-Institut für Softwaresystemtechnik GmbH
//All rights reserved.
//
//Redistribution and use in source and binary forms, with or without
//modification, are permitted provided that the following conditions are met:
//
//* Redistributions of source code must retain the above copyright notice, this
//list of conditions and the following disclaimer.
//
//* Redistributions in binary form must reproduce the above copyright notice,
//this list of conditions and the following disclaimer in the documentation
//and/or other materials provided with the distribution.
//
//* Neither the name of the copyright holder nor the names of its
//contributors may be used to endorse or promote products derived from
//this software without specific prior written permission.
//
//THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
//AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
//IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
//FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
//DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
//SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
//CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
//OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package zerokit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	ListTresorMembersPath = "/api/v4/admin/tresor/list-members"
	InitiateUserRegistrationPath = "/api/v4/admin/user/init-user-registration"
	ApproveTresorCreationPath = "/api/v4/admin/tresor/approve-tresor-creation"
	ValidateUserRegistrationPath = "/api/v4/admin/user/validate-user-registration"
)

type tresoritClient struct {
	requestSigner
	httpClient httpClient
	ServiceUrl url.URL
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewTresoritClient(serviceUrl, adminUserId, adminKey string) (*tresoritClient, error) {
	if serviceUrl == "" || adminKey == "" || adminUserId == "" {
		return nil, errors.New("one or more arguments are empty")
	}

	u, err := url.Parse(serviceUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid service url: %s", serviceUrl))
	}

	return &tresoritClient{
		requestSigner: requestSigner{
			adminKey:    adminKey,
			adminUserId: adminUserId,
		},
		httpClient: http.DefaultClient,
		ServiceUrl: *u,
	}, nil
}

func (c *tresoritClient) doSignedPost(urlPath string,
	body []byte) (*http.Response, error) {
	endpoint := c.ServiceUrl
	endpoint.Path = path.Join(endpoint.Path, urlPath)
	r, err := http.NewRequest("POST", endpoint.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return c.SignAndDo(r)
}

func (c *tresoritClient) doSignedGet(urlPath string,
	query url.Values) (*http.Response, error) {
	endpoint := c.ServiceUrl
	endpoint.Path = path.Join(endpoint.Path, urlPath)
	r, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = query.Encode()
	return c.SignAndDo(r)
}

func (c *tresoritClient) SignAndDo(req *http.Request) (*http.Response, error) {
	err := c.sign(req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *tresoritClient) ListTresorMembers(tresorId string) ([]string, error) {
	q := url.Values{}
	q.Add("tresorid", tresorId)

	resp, err := c.doSignedGet(ListTresorMembersPath, q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := map[string][]string{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return m["Members"], nil
}

func (c *tresoritClient) InitUserRegistration() (*UserRegistrationData, error) {
	resp, err := c.doSignedPost(InitiateUserRegistrationPath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reg UserRegistrationData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

type UserRegistrationData struct {
	SessionId       string `json:"RegSessionId"`
	SessionVerifier string `json:"RegSessionVerifier"`
	UserId          string `json:"UserId"`
}

func (c *tresoritClient) ApproveTresorCreation(tresorId string) error {
	m := map[string]string{"TresorId": tresorId}
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := c.doSignedPost(ApproveTresorCreationPath, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *tresoritClient) ValidateUserRegistration(zeroKitId, sessionId,
	sessionVerifier, validationVerifier string) error {
	m := orderedMap{
		{"RegSessionId", sessionId},
		{"RegSessionVerifier", sessionVerifier},
		{"RegValidationVerifier", validationVerifier},
		{"UserId", zeroKitId},
	}
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := c.doSignedPost(ValidateUserRegistrationPath, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type orderedMap []mapItem

type mapItem struct {
	key   string
	value interface{}
}

func (m orderedMap) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, item := range m {
		if i != 0 {
			buffer.WriteString(",")
		}
		key, err := json.Marshal(item.key)
		if err != nil {
			return nil, err
		}
		buffer.Write(key)
		buffer.WriteString(":")
		val, err := json.Marshal(item.value)
		if err != nil {
			return nil, err
		}
		buffer.Write(val)
	}

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
