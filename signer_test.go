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
	"strings"
	"testing"
	"time"
)

const (
	AdminUserId    = "admin@exampletenant.tresorit.io"
	AdminKey       = "204bcf1b"
	Sha256HexEmpty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

type hmac256 struct {
	data   []byte
	secret string
	hmac   string
}

var testDataHmac256 = []hmac256{
	{[]byte(""), AdminKey, "ip8xpaW+rCdNJgqSGAAeeIcUdn9waFlCAdcyj4GeRUc="},
	{[]byte(nil), AdminKey, "ip8xpaW+rCdNJgqSGAAeeIcUdn9waFlCAdcyj4GeRUc="},
	{[]byte("{}"), AdminKey, "wYn05Tvf7BB3OvnCNLAfrQ7/3/p5gCX4sGDyxgn6y2I="},
}

func TestComputeHmac256(t *testing.T) {
	for _, test := range testDataHmac256 {
		r, err := computeHmacSHA256(test.data, test.secret)
		if err != nil {
			t.Errorf(
				"cannot compute hmac256 for data: %q and secret: %s",
				test.data, test.secret,
			)
		}
		hmac256hex := base64.StdEncoding.EncodeToString(r)
		if hmac256hex != test.hmac {
			t.Errorf(
				"computeHmacSHA256(%q, %q) = %q, want %q",
				test.data, test.secret, hmac256hex, test.hmac,
			)
		}
	}
}

func TestComputeHmac256InvalidSecret(t *testing.T) {
	_, err := computeHmacSHA256(nil, "0g")
	if err == nil {
		t.Errorf(
			"expected invalid byte error for data: %q and secret: %s",
			nil, "0g",
		)
	}
}

type bytesHex struct {
	data []byte
	hex  string
}

var testDataSha256hex = []bytesHex{
	{[]byte(""), Sha256HexEmpty},
	{[]byte(nil), Sha256HexEmpty},
	{[]byte("xyz"), "3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282"},
}

func TestSha256hex(t *testing.T) {
	for _, test := range testDataSha256hex {
		hex := sha256hex(test.data)
		if hex != test.hex {
			t.Errorf("sha256hex(%q) = %q, want %q", test.data, hex, test.hex)
		}
	}
}

func TestPostRequestSigningWithContent(t *testing.T) {
	s := requestSigner{adminUserId: AdminUserId, adminKey: AdminKey}

	content := []byte("{\"TresorId\":\"e32ve3ve\"}")
	r, _ := http.NewRequest("POST", "", bytes.NewBuffer(content))
	err := s.sign(r)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "Content-Type") != "application/json" {
		t.Errorf(
			"Content-Type = %s; want %s",
			header(r, "Content-Type"), "application/json",
		)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId = %s; want %s", header(r, "UserId"), AdminUserId)
	}

	ts, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", ts)
	}

	contentSha256hex := sha256hex(content)
	if header(r, "Content-SHA256") != contentSha256hex {
		t.Errorf(
			"Content-SHA256 = %s; want %s",
			header(r, "Content-SHA256"), contentSha256hex,
		)
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Errorf("auth type = %s; want %s", auth[0], "AdminKey")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("auth credentials are not in base64 encoding: %s", auth[1])
	}
}

func TestPostRequestSigningWithoutContent(t *testing.T) {
	s := requestSigner{adminUserId: AdminUserId, adminKey: AdminKey}

	r, _ := http.NewRequest("POST", "", nil)
	err := s.sign(r)
	if err != nil {
		t.Errorf("cannot sign the request: %v", r)
	}

	if header(r, "Content-Type") != "application/json" {
		t.Errorf(
			"Content-Type = %s; want %s",
			header(r, "Content-Type"), "application/json",
		)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId = %s; want %s", header(r, "UserId"), AdminUserId)
	}

	ts, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %v", ts)
	}

	if header(r, "Content-SHA256") != Sha256HexEmpty {
		t.Errorf(
			"Content-SHA256 = %s; want %s",
			header(r, "Content-SHA256"), Sha256HexEmpty,
		)
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Errorf("UserId = %s; want %s", auth[0], "AdminKey")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("auth credentials are not in base64 encoding: %s", auth[1])
	}
}

func header(r *http.Request, key string) string {
	if v, ok := r.Header[key]; ok {
		return v[0]
	}
	return ""
}

func TestGetRequestSigning(t *testing.T) {
	s := requestSigner{adminUserId: AdminUserId, adminKey: AdminKey}

	r, _ := http.NewRequest("GET", "", nil)
	err := s.sign(r)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId = %s; want %s", header(r, "UserId"), AdminUserId)
	}

	if header(r, "Content-Type") != "" {
		t.Errorf("Content-Type = %s; want %s", header(r, "Content-Type"), "")
	}

	if header(r, "Content-SHA256") != "" {
		t.Errorf("Content-SHA256 = %s; want %s", header(r, "Content-SHA256"), "")
	}

	ts, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %v", ts)
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Errorf("auth type = %s; want %s", auth[0], "AdminKey")

	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("auth credentials are not in base64 encoding: %s", auth[1])
	}
}

func TestGetRequestSigningWithQueryParameters(t *testing.T) {
	s := requestSigner{adminUserId: AdminUserId, adminKey: AdminKey}

	r, _ := http.NewRequest("GET", "", nil)
	q := r.URL.Query()
	q.Add("tresorId", "0000v6c5wl03ms87ldqf9p8r")
	r.URL.RawQuery = q.Encode()

	err := s.sign(r)
	if err != nil {
		t.Errorf("cannot sign request: %v", r)
	}

	if header(r, "UserId") != AdminUserId {
		t.Errorf("UserId = %s; want %s", header(r, "UserId"), AdminUserId)
	}

	if header(r, "Content-Type") != "" {
		t.Errorf("Content-Type = %s; want %s", header(r, "Content-Type"), "")
	}

	if header(r, "Content-SHA256") != "" {
		t.Errorf("Content-SHA256 = %s; want %s", header(r, "Content-SHA256"), "")
	}

	ts, err := time.Parse(time.RFC3339, header(r, "TresoritDate"))
	if err != nil {
		t.Errorf("timestamp is not ISO-8601: %s", ts)
	}

	auth := strings.Split(header(r, "Authorization"), " ")
	if auth[0] != "AdminKey" {
		t.Errorf("auth type = %s; want %s", auth[0], "AdminKey")
	}
	_, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		t.Errorf("auth credentials are not in base64 encoding: %s", auth[1])
	}

	s.adminKey = "0g"
	r, _ = http.NewRequest("GET", "", nil)
	err = s.sign(r)
	if err == nil {
		t.Error("expected invalid byte error")
	}
}
