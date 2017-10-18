[![Build Status](https://travis-ci.org/gesundheitscloud/go-zerokit-api-client.svg?branch=master)](https://travis-ci.org/gesundheitscloud/go-zerokit-api-client)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c33f5a21aeaf497a85ebf9acfb797939)](https://www.codacy.com/app/theintz/go-zerokit-api-client?utm_source=github.com&utm_medium=referral&utm_content=gesundheitscloud/go-zerokit-api-client&utm_campaign=badger)
[![codecov](https://codecov.io/gh/gesundheitscloud/go-zerokit-api-client/branch/master/graph/badge.svg)](https://codecov.io/gh/gesundheitscloud/go-zerokit-api-client)

### File updated by Toby Moncaster on 15th October 2017. ([toby@moncaster.com](mailto:toby@moncaster.com))
-------
# ZeroKit Admin API client for Go

A [Golang](https://golang.org/) implementation of ZeroKit's tenant admin API client library.

For further information please see the following:

- [ZeroKit encryption platform](https://tresorit.com/zerokit)
- [ZeroKit management portal](https://manage.tresorit.io)
- [ZeroKit Admin API Reference](https://tresorit.com/zerokit/docs/admin_api/API_reference.html)

## Introduction

This API uses ZeroKit, a stateless authentication and encryption protocol. Because ZeroKit is stateless, every HTTP request has to be authenticated seperately. This makes it a little different to some other authentication and encryption protocols. ZeroKit consists of two parts, the Tresorit ZeroKit server side that managed the creation and storage of tresors and a per-application client side that manages local authentication of users and tresors.

## Tresors

The ZeroKit API uses *tresors* to store keys in a secure fashion. Tresors can be created by any user, however the tresor can only be used once it has been authorised by the client backend. Creation of a tresor generates a unique ID, but as the server can't provide a list of IDs after creation these must be stored on the client side. However data that has been encrypted with a given tresor includes the tresor ID, so it is robust to loss of the ID.

## ZeroKit Admin Client API

Requests to this API provide the following functions:

 - *InitUserRegistration* – this initiates the new user creation process.
 - *ValidateUser* - validate a newly created user. Prior to calling this the user can be added to Tresors but cannot login.
 - *ApproveTresorCreation* – approve the creation of a new tresor. 
 - *ListMembers* – this provides a list of all users who are members of a given tresor.

Since user registration is a two step process it allows the use of out-of-band authentication methods (e.g. email address verification). 
 
## Examples

Initiate a new user registration using **InitUserRegistration**:
 
```go
package main

import (
    "net/url"
    "io/ioutil"
    "github.com/gesundheitscloud/go-zerokit-api-client"
    "path"
    "net/http"
    "fmt"
)

func main() {
    client, err := zerokit.NewZeroKitAdminApiClient(
        "https://example.api.tresorit.io",
        "admin@example.tresorit.io",
        "fsdfq34r2efe",
    )
    if err != nil {
        return err
    }
    u, err := url.Parse(client.ServiceUrl)
    if err != nil {
        return err
    }
    u.Path = path.Join(u.Path, "/api/v4/admin/user/init-user-registration")
    r, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return err
    }

    resp, err := client.SignAndDo(r)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // do something with response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(string(body))
}
```

Validate a new user using **ValidateUser**:
 
```go
package main

import (
    "net/url"
    "github.com/gesundheitscloud/go-zerokit-api-client"
    "path"
    "net/http"
    "fmt"
)

func main() {
    client, err := zerokit.NewZeroKitAdminApiClient(
        "https://example.api.tresorit.io",
        "admin@example.tresorit.io",
        "fsdfq34r2efe",
    )
    if err != nil {
        return err
    }
    u, err := url.Parse(client.ServiceUrl)
    if err != nil {
        return err
    }
    u.Path = path.Join(u.Path, "/api/v4/admin/user/validate-user-registration")
    r, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return err
    }

}
```

Approve the creation of a new tresor with **ApproveTresorCreation**:

```go
package main

import (
    "net/url"
    "github.com/gesundheitscloud/go-zerokit-api-client"
    "path"
    "net/http"
    "fmt"
)

func main() {
    client, err := zerokit.NewZeroKitAdminApiClient(
        "https://example.api.tresorit.io",
        "admin@example.tresorit.io",
        "fsdfq34r2efe",
    )
    if err != nil {
        return err
    }
    u, err := url.Parse(client.ServiceUrl)
    if err != nil {
        return err
    }
    u.Path = path.Join(u.Path, "/api/v4/admin/user/approve-tresor-creation")
    r, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return err
    }

}
```

Lists all members of a given tresor using **ListMembers**:

```go
package main

import (
    "github.com/gesundheitscloud/go-zerokit-api-client"
    "fmt"
)

func main() {
    client, err := zerokit.NewZeroKitAdminApiClient(
        "https://example.api.tresorit.io",
        "admin@example.tresorit.io",
        "fsdfq34r2efe",
    )
    if err != nil {
        return err
    }

    members, err := client.ListTresorMembers("0000slpj4r86xbqlg9wmjhug")
    if err != nil {
        return err
    }
    fmt.Println(members)
}
```
