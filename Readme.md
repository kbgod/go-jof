# JSON Optional field
This is a simple library to handle optional fields in JSON.

#### Installation
```bash
go get github.com/kbgod/go-jof
```

#### Usage
The most popular case is to handle optional fields in data updating handlers.
```go
type UpdateUser struct {
	Username    jof.Field[string]  `json:"name"`
	CallbackURL jof.Field[*string] `json:"age"`
}

func UpdateUserHandler(rw http.ResponseWriter, req *http.Request) {
	var u UpdateUser
	if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	mustUpdate := map[string]any{}

	if !u.Username.Defined {
		http.Error(rw, "name is required", http.StatusBadRequest)
		return
	}
	mustUpdate["name"] = u.Username.Value // username can be empty string, but not null.

	if u.CallbackURL.Defined {
		// update only if defined.
		// null value means to remove the callback URL.
		mustUpdate["callback_url"] = u.CallbackURL.Value
	}
	
	db.Model(&User{}).Where("id = ?", req.URL.Query().Get("id")).Updates(mustUpdate)
}
```

You can also use `Field` to marshal JSON. But undefined field will be marshaled as `null`. If you want to omit undefined fields, you should use native type.
```go
type User struct {
	ID    int
	Item0 jof.Field[*int]
	Item1 jof.Field[*int]
	Item2 jof.Field[*int]
}

func main() {
	idx := 1
	u := User{
		Item0: jof.NewField[*int](),    // undefined
		Item1: jof.NewField(&idx),      // defined with non-nil value
		Item2: jof.NewField[*int](nil), // defined with nil value
	}
	bytes, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes)) // {"ID":0,"Item0":null,"Item1":1,"Item2":null}
}
```
