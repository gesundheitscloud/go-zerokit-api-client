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
	"github.com/go-errors/errors"
	"net/http"
)

type tresoritClient struct {
	requestSigner
	httpClient httpClient
	serviceUrl string
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
		serviceUrl: serviceUrl,
	}, nil
}

func (c *tresoritClient) SignAndDo(req *http.Request) (*http.Response, error) {
	err := c.sign(req)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
