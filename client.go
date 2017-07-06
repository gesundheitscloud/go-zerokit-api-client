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
	"encoding/json"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	listTresorMembers        = "/api/v4/admin/tresor/list-members"
	initiateUserRegistration = "/api/v4/admin/user/init-user-registration"
)

type tresoritClient struct {
	requestSigner
	httpClient httpClient
	ServiceUrl string
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewTresoritClient(serviceUrl, adminUserId, adminKey string) (*tresoritClient, error) {
	if serviceUrl == "" || adminKey == "" || adminUserId == "" {
		return nil, errors.New("one or more arguments are empty")
	}

	return &tresoritClient{
		requestSigner: requestSigner{
			adminKey:    adminKey,
			adminUserId: adminUserId,
		},
		httpClient: http.DefaultClient,
		ServiceUrl: serviceUrl,
	}, nil
}

func (c *tresoritClient) SignAndDo(req *http.Request) (*http.Response, error) {
	err := c.sign(req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *tresoritClient) ListTresorMembers(tresorId string) ([]string, error) {
	u, err := url.Parse(c.ServiceUrl)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, listTresorMembers)
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Add("tresorid", tresorId)
	r.URL.RawQuery = q.Encode()

	resp, err := c.SignAndDo(r)
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

type UserRegistrationData struct {
	SessionVerifier string `json:"RegSessionVerifier"`
	SessionId       string `json:"RegSessionId"`
	UserId          string `json:"UserId"`
}

func (z *tresoritClient) InitUserRegistration() (*UserRegistrationData, error) {
	u, err := url.Parse(z.ServiceUrl)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, initiateUserRegistration)
	r, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := z.SignAndDo(r)
	if err != nil {
		return nil, err
	}

	var reg UserRegistrationData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&reg)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return &reg, nil
}
