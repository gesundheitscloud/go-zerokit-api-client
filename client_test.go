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
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

const (
	ServiceUrl  = "https://exampletenant.tresorit.io"
	AdminUserId = "admin@exampletenant.tresorit.io"
)

type mockHttpClient struct {
	DoMock func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoMock != nil {
		return m.DoMock(req)
	}
	return &http.Response{}, nil
}

func TestListTresorMembers(t *testing.T) {
	tresorId := "xyz"
	tresorMembers := []string{"zk1", "zk2"}

	// mock list members response
	client := &mockHttpClient{
		DoMock: func(req *http.Request) (*http.Response, error) {
			tresorIdActual := req.URL.Query().Get("tresorid")
			if tresorIdActual != tresorId {
				t.Errorf(
					"tresorid query parameter = %s, want = %s",
					tresorIdActual, tresorId)
			}

			m := map[string][]string{}
			m["Members"] = append(m["Members"], tresorMembers...)
			body, _ := json.Marshal(m)
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}, nil
		},
	}
	c, err := NewTresoritClient(ServiceUrl, AdminUserId, AdminKey)
	if err != nil {
		t.Fatal("cannot initialize tresorit client")
	}
	// package protected
	c.httpClient = client

	members, err := c.ListTresorMembers(tresorId)
	if len(members) != len(tresorMembers) {
		t.Errorf(
			"number of tresor's members = %d, want = %d",
			len(members), len(tresorMembers),
		)
	}

	if !reflect.DeepEqual(members, tresorMembers) {
		t.Errorf(
			"tresor's members = %v, want = %v", members, tresorMembers)
	}
}
