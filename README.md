# SurrealDB Record ID library

SurrealDB Record ID lib in pure Go

## Usage

Start a SurrealDB instance via

```sh
$ surreal start --auth -u root -p root memory
```

```go
package main

import (
	"fmt"

	thing "github.com/dbaudisch/go-surreal-id"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/conn/gorilla"
	"github.com/surrealdb/surrealdb.go/pkg/marshal"
)

type User struct {
	ID      *thing.Thing `json:"id,omitempty"`
	Name    string       `json:"name"`
	Surname string       `json:"surname"`
}

func main() {
	db, err := surrealdb.New("ws://localhost:8000/rpc", gorilla.Create())
	if err != nil {
		panic(err)
	}

	auth := &surrealdb.Auth{
		Database:  "test",
		Namespace: "test",
		Username:  "root",
		Password:  "root",
	}
	if _, err = db.Signin(auth); err != nil {
		panic(err)
	}

	user := User{Name: "Tobie", Surname: "Morgan Hitchcock"}

	u, err := db.Create("user", user)
	res, err := marshal.SmartUnmarshal[User](u, err)
	if err != nil {
		panic(err)
	}

	fmt.Printf("CREATE %s: %v\n\n", res[0].ID, res)

	if res, err = marshal.SmartUnmarshal[User](db.Select(res[0].ID.String())); err != nil {
		panic(err)
	}

	fmt.Printf("SELECT * FROM %s:\n%v\n%#v\n", res[0].ID, res, res)
}
```
