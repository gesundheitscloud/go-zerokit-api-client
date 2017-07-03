# ZeroKit Admin API client for Go

ZeroKit tenant's admin API client library in [Golang](https://golang.org/).

For further information please see:

- [ZeroKit encryption platform](https://tresorit.com/zerokit)
- [ZeroKit management portal](https://manage.tresorit.io)


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
    )

    func main() {
        zk := zerokit.ZeroKitAdminAPIClient{
            ServiceUrl: "https://example.api.tresorit.io",
            AdminKey: "",
            AdminUserId: "admin@example.tresorit.io",
        }

        u, err := url.Parse(zk.ServiceUrl)
        if err != nil {
            return err
        }
        u.Path = path.Join(u.Path, "/api/v4/admin/user/init-user-registration")
        r, err := http.NewRequest("POST", u.String(), nil)
        if err != nil {
            return err
        }

        resp, err := zk.SignAndDo(r, nil)
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

Lists all members of the given tresor:

```go
    package main

    import (
        "net/url"
        "io/ioutil"
        "github.com/hpihc/go-tresorit-api-client"
        "path"
        "net/http"
    )

    func main() {
        zk := zerokit.ZeroKitAdminAPIClient{
            ServiceUrl: "https://example.api.tresorit.io",
            AdminKey: "",
            AdminUserId: "admin@example.tresorit.io",
        }

        u, err := url.Parse(zk.ServiceUrl)
        if err != nil {
            return err
        }
        u.Path = path.Join(u.Path, "/api/v4/admin/tresor/list-members")
        r, err := http.NewRequest("GET", u.String(), nil)
        if err != nil {
            return err
        }
        q := r.URL.Query()
        q.Add("tresorid", "0000v6c5wl03ms87ldqf9p8r")
        r.URL.RawQuery = q.Encode()

        resp, err := zk.SignAndDo(r, nil)
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
