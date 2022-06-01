## Account Golang SDK

The library is a Golang package for [Single Account](http://account.corezoid.com) service API.  

Get redirect authorization URL:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go/account/authurl"
    "github.com/corezoid/sdk-go/account/oauth"
)

func main() {
    api := authurl.New(oauth.Client{
        ClientId: "anId",
        ClientSecret: "aSecret",
        HttpHost: "https://account.corezoid.com",
        RedirectUri: "https://your.service/back/url",
    })
    
    url := api.AuthorizeUrl([]string{"single_account:account.read"})
    
    // redirect user to the url to processed with authorization.
}
```

Get access token by authorize code:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go/account/authcode"
    "github.com/corezoid/sdk-go/account/oauth"
)

func main() {
    authCode := oauth.AuthCode("anAuthorizationCode")

    api := authcode.New(oauth.Client{
        ClientId: "anId",
        ClientSecret: "aSecret",
        HttpHost: "https://account.corezoid.com",
        RedirectUri: "https://your.service/back/url",
    }, nil)
    
    r := api.Request(authCode).Decode()
    if r.Err != nil {
    	panic(r.Err)
    }
    defer r.Close()

    log.Printf("%+v", r)
}
```

Refresh access token:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go/account/authrefreshtoken"
    "github.com/corezoid/sdk-go/account/oauth"
)

func main() {
    token := oauth.RefreshToken("aRefreshToken")

    api := authrefreshtoken.New(oauth.Client{
        ClientId: "anId",
        ClientSecret: "aSecret",
        HttpHost: "https://account.corezoid.com",
        RedirectUri: "https://your.service/back/url",
    }, nil)
    
    r := api.Request(token).Decode()
    if r.Err != nil {
    	panic(r.Err)
    }
    defer r.Close()

    log.Printf("%+v", r)
}
```

Get user info:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go/account/userinfo"
    "github.com/corezoid/sdk-go/account/oauth"
)

func main() {
    token := oauth.AccessToken("anAccessToken")

    api := userinfo.New(oauth.Client{
        ClientId: "anId",
        ClientSecret: "aSecret",
        HttpHost: "https://account.corezoid.com",
        RedirectUri: "https://your.service/back/url",
    }, nil)
    
    r := api.Request(token).Decode()
    if r.Err != nil {
    	panic(r.Err)
    }
    defer r.Close()

    log.Printf("%+v", r)
}
```

Logout:

```go
package main

import (
	"log"
	"github.com/corezoid/sdk-go/account/logout"
	"github.com/corezoid/sdk-go/account/oauth"
)

func main() {
	token := oauth.AccessToken("anAccessToken")

	api := logout.New(oauth.Client{
		ClientId:     "anId",
		ClientSecret: "aSecret",
		HttpHost:     "https://account.corezoid.com",
		RedirectUri:  "https://your.service/back/url",
	}, nil)

	r := api.Request(token).Decode()
	if r.Err != nil {
		panic(r.Err)
	}
	defer r.Close()

	log.Printf("%+v", r)
}
```
