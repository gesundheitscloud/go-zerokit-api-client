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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

type ZeroKitAdminAPIClient struct {
	http.Client
	ServiceUrl  string
	AdminUserId string
	AdminKey    string
}

func (z *ZeroKitAdminAPIClient) SignAndDo(r *http.Request,
	b []byte) (*http.Response, error) {
	err := z.sign(r, b)
	if err != nil {
		return nil, err
	}
	return z.Do(r)
}

// As all server-side APIs of the Tresorit platform are stateless HTTP services,
// all requests made against it have to be authenticated individually. This is
// done by signing all request with the Tenant-specific admin key with the
// special scheme described below:
//
// 1. Assemble the request headers:
//
//     - Content-Type: application/json (Post only)
//     - Content-SHA256: <sha256hex of the request body> (Post only)
//     - TresoritDate: <timestamp in ISO-8601>
//     - UserId: <tenant admin user id>
//     - HMACHeaders: <comma separated headers>
//
// 2. Assemble the canonicalized request string, which will be the subject of
// the signing. The canonical string complies to the following format:
//
// 		<request verb>\n +
// 		<path>[?key=value [& ...] ]\n +
// 		header:value[\n ...]
//
// The headers key-value pair must be listed in the same order as in the
// HMACHeaders header.
//
// 3. Sing the request. The algorithm is the following:
//
//     BASE64ENCODE(
//     		HMACSHA256(
//     			key=HEXTOBIN(AdminKey),
//     			data=UTF8TOBINARY(Canonical Request String)
//     		)
//     )
//
// 4. Add authorization header:
//
// 		Authorization: AdminKey <SignatureBase64>
//
// where SignatureBase64 is the previously computed signature of the canonical
// request string.
func (z *ZeroKitAdminAPIClient) sign(req *http.Request, content []byte) error {
	if req.Method == "POST" {
		req.Header["Content-Type"] = []string{"application/json"}
		req.Header["Content-SHA256"] = []string{sha256hex(content)}
	}

	// The Add and Set methods of http.Header canonicalize header names when
	// adding values to the header map, which result in the incorrect signature
	// validation by the tresorit API. Therefore, we bypass the behavior of the
	// Set and Get by setting the headers using map operation.
	//req.Header["Content-Length"] = []string{strconv.Itoa(len(content))}
	req.Header["TresoritDate"] = []string{time.Now().UTC().Format(time.RFC3339)}
	req.Header["UserId"] = []string{z.AdminUserId}
	req.Header["HMACHeaders"] = []string{}

	var headers []string
	for k := range req.Header {
		headers = append(headers, k)
	}
	req.Header["HMACHeaders"] = []string{strings.Join(headers, ",")}

	// assemble a canonicalized string of the requests
	var buffer bytes.Buffer
	buffer.WriteString(req.Method + "\n")
	buffer.WriteString(strings.TrimPrefix(req.URL.Path, "/"))
	if req.URL.RawQuery != "" {
		buffer.WriteString("?" + req.URL.RawQuery)
	}
	for _, key := range headers {
		buffer.WriteString("\n")
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(req.Header[key][0])

	}

	// sign the canonicalized string of the requests
	sig, err := computeHmacSHA256([]byte(buffer.String()), z.AdminKey)
	if err != nil {
		return err
	}

	// Set authorization header.
	req.Header["Authorization"] = []string{
		"AdminKey " + base64.StdEncoding.EncodeToString(sig),
	}
	return nil
}

// The function to compute the keyed-hash message authentication code (HMAC)
// of the input data, keyed with the binary input key and using the SHA256
// algorithm as hash function.
func computeHmacSHA256(data []byte, secret string) ([]byte, error) {
	key := []byte(secret)
	dst := make([]byte, hex.DecodedLen(len(key)))
	_, err := hex.Decode(dst, key)
	if err != nil {
		return nil, err
	}
	h := hmac.New(sha256.New, dst)
	h.Write(data)
	return h.Sum(nil), nil
}

func sha256hex(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}
