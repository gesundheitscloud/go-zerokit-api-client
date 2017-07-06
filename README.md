[![Build Status](https://travis-ci.com/hpihc/go-tresorit-api-client.svg?token=phqzPcrEVgWxvsyLbP6L&branch=master)](https://travis-ci.com/hpihc/go-tresorit-api-client)

# ZeroKit Admin API client for Go

ZeroKit tenant's admin API client library in [Golang](https://golang.org/).

For further information please see:

- [ZeroKit encryption platform](https://tresorit.com/zerokit)
- [ZeroKit management portal](https://manage.tresorit.io)
- [Tresorit Admin API Reference](https://tresorit.com/zerokit/docs/admin_api/API_reference.html)

## Tresorit API

Implemented Tresorit API methods:

 - ListMembers

## Examples

Initiate a user registration process:
 
```go
    package main

    import (
        "net/url"
        "io/ioutil"
        "github.com/hpihc/go-tresorit-api-client"
        "path"
        "net/http"
        "fmt"
    )

    func main() {
        client, err := zerokit.NewTresoritClient(
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

Lists all members of the given tresor using implemented Tresorit API methods:

```go
    package main

    import (
        "github.com/hpihc/go-tresorit-api-client"
        "fmt"
    )

    func main() {
        client, err := zerokit.NewTresoritClient(
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
