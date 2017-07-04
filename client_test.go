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
	"encoding/base64"
	"net/http"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"
)

const (
	ServiceUrl     = "https://exampletenant.tresorit.io"
	AdminUserId    = "admin@exampletenant.tresorit.io"
	AdminKey       = "204bcf1b"
	Sha256HexEmpty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

func TestComputeHmac256(t *testing.T) {
	hmac, err := computeHmacSHA256([]byte(""), AdminKey)
	if err != nil {
		t.Error(err)
	}
	if base64.StdEncoding.EncodeToString(hmac) != "ip8xpaW+rCdNJgqSGAAeeIcUdn9waFlCAdcyj4GeRUc=" {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}

	hmac, err = computeHmacSHA256([]byte("{}"), AdminKey)
	if err != nil {
		t.Error(err)
	}
	if base64.StdEncoding.EncodeToString(hmac) != "wYn05Tvf7BB3OvnCNLAfrQ7/3/p5gCX4sGDyxgn6y2I=" {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}

	hmac, err = computeHmacSHA256(nil, AdminKey)
	if err != nil {
		t.Error(err)
	}
	if base64.StdEncoding.EncodeToString(hmac) != "ip8xpaW+rCdNJgqSGAAeeIcUdn9waFlCAdcyj4GeRUc=" {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}

	hmac, err = computeHmacSHA256(nil, "0g")
	if err == nil {
		t.Error("expected invalid byte error; got none")
	}
}

func TestSha256hex(t *testing.T) {
	if sha256hex([]byte("")) != Sha256HexEmpty {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}

	if sha256hex([]byte(nil)) != Sha256HexEmpty {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}

	if sha256hex([]byte("xyz")) != "3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282" {
		t.Error("incorrect base64 keyed-hash message of the input data")
	}
}

func TestPostRequestSigning(t *testing.T) {
	zk := ZeroKitAdminAPIClient{
		ServiceUrl:  ServiceUrl,
		AdminUserId: AdminUserId,
		AdminKey:    AdminKey,
	}
	u, _ := url.Parse(ServiceUrl)
	u.Path = path.Join(u.Path, "/api/v1/somepath")

	// post request without content
	r, _ := http.NewRequest("POST", u.String(), nil)
	err := zk.sign(r, nil)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "Content-Type") != "application/json" {
		t.Errorf("content type must be \"application/json\", was %s",
			header(r, "Content-Type"))
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId type must be %s, was %s",
			AdminUserId, header(r, "UserId"))
	}

	v, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", v)
	}

	if header(r, "Content-SHA256") != Sha256HexEmpty {
		t.Errorf("content sha-256 must be %s, was %s",
			Sha256HexEmpty, header(r, "Content-SHA256"))
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Error("authorization type must be \"AdminKey\"")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("authorization credentilas are not base64: %s", auth[1])
	}

	// post request with content
	content := []byte("{\"TresorId\":\"e32ve3ve\"}")
	r, _ = http.NewRequest("POST", u.String(), bytes.NewBuffer(content))
	err = zk.sign(r, content)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "Content-Type") != "application/json" {
		t.Errorf("content type must be \"application/json\", was %s",
			header(r, "Content-Type"))
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId type must be %s, was %s",
			AdminUserId, header(r, "UserId"))
	}

	v, err = time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", v)
	}

	if header(r, "Content-SHA256") != "c376cd2c3528f9766119a9e2cb5bc2df47c1f89727d454041bf288a44c23d866" {
		t.Errorf("content sha-256 must be %s, was %s",
			"c376cd2c3528f9766119a9e2cb5bc2df47c1f89727d454041bf288a44c23d866",
			header(r, "Content-SHA256"),
		)
	}

	auth = strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Error("authorization type must be \"AdminKey\"")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("authorization credentilas are not base64: %s", auth[1])
	}
}

func header(r *http.Request, key string) string {
	if v, ok := r.Header[key]; ok {
		return v[0]
	}
	return ""
}

func TestGetRequestSigning(t *testing.T) {
	zk := ZeroKitAdminAPIClient{
		ServiceUrl:  ServiceUrl,
		AdminUserId: AdminUserId,
		AdminKey:    AdminKey,
	}
	u, _ := url.Parse(ServiceUrl)
	u.Path = path.Join(u.Path, "/api/v1/somepath")

	// get request
	r, _ := http.NewRequest("GET", u.String(), nil)
	err := zk.sign(r, nil)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId type must be %s, was %s",
			AdminUserId, header(r, "UserId"))
	}

	if header(r, "Content-Type") != "" {
		t.Error("get request must not the \"Content-Type\" header")
	}

	if header(r, "Content-SHA256") != "" {
		t.Error("get request must not the \"Content-SHA256\" header")
	}

	v, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", v)
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Error("authorization type must be \"AdminKey\"")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("authorization credentilas are not base64: %s", auth[1])
	}

	// get request with query
	r, _ = http.NewRequest("GET", u.String(), nil)
	q := r.URL.Query()
	q.Add("tresorId", "0000v6c5wl03ms87ldqf9p8r")
	r.URL.RawQuery = q.Encode()

	err = zk.sign(r, nil)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId type must be %s, was %s",
			AdminUserId, header(r, "UserId"))
	}

	if header(r, "Content-Type") != "" {
		t.Error("get request must not the \"Content-Type\" header")
	}

	if header(r, "Content-SHA256") != "" {
		t.Error("get request must not the \"Content-SHA256\" header")
	}

	v, err = time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", v)
	}

	auth = strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Error("authorization type must be \"AdminKey\"")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("authorization credentilas are not base64: %s", auth[1])
	}

	zk.AdminKey = "0g"
	r, _ = http.NewRequest("GET", u.String(), nil)
	err = zk.sign(r, nil)
	if err == nil {
		t.Error("expected invalid byte error; got none")
	}
}
