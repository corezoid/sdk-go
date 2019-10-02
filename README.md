# Go-lang SDK for Corezoid

The library aim is to simplify communication with [Corezoid](https://corezoid.com/) API for Golang developers.
The client implements Corezoid V2 API. For more details, please check the [API documentation](https://doc.corezoid.com/en/api/v2/) out.  

Create task example:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go"
)

func main() {
    secret := "theSecret"
    login := 111222
    
    client := corezoid.New(login, secret)

	ops := corezoid.Ops{}
	ops.Add(corezoid.MapOp{
		"type": "create",
		"conv_id": 123456,
		"obj": "task",
		"data": map[string]interface{}{
			"key": "taskPayload",
		},
	})

	res := client.Call(ops).Decode()
	if res.Err != nil {
		panic(res.Err)
	}
	defer res.Close()

	log.Printf("%+v", res)
    // &{RequestProc:ok Ops:[map[id: obj:task obj_id:5d94707560e327394302ceec proc:ok ref:<nil>]] Response:0xc00028c000 Err:<nil>}
}
```

Upload process schema:

```go
package main

import (
    "log"
    "github.com/corezoid/sdk-go"
)

func main() {
    // you can schema by selecting a process on Corezoid and downloading it as JSON file. 
    schema := `[{}]`

    secret := "theSecret"
    login := 111222
    
    client := corezoid.New(login, secret)

	res := client.Upload(corezoid.MapOp{
		"type":      "create",
		"obj":       "obj_scheme",
		"folder_id": "289791",
		"scheme":    scheme.Scheme0,
		"async": "false",
	}).Decode()

	if res.Err != nil {
		panic(res.Err)
	}
	defer res.Close()

	log.Printf("%+v", res)
    // &{RequestProc:ok Ops:[map[obj:obj_scheme proc:ok scheme:[map[description: hash:6ae60c09d82ca51fd1c33b7a3cae5c800a5b1a17 obj_id:614098 obj_type:conv old_obj_id:610371 old_parent_id:0 title:in]]]] Response:0xc00037a000 Err:<nil>}
}


