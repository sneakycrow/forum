package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type graphQLRequest struct {
	Query     string                  `json:"query"`
	Variables newUserGraphQLVariables `json:"variables"`
}

type newUserGraphQLVariables struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type user struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type hasuraResponse struct {
	Data hasuraInsertUsers `json:"data"`
}

type hasuraInsertUsers struct {
	InsertUsers hasuraInsertUsersData `json:"insert_users"`
}

type hasuraInsertUsersData struct {
	AffectedRows  int    `json:"affected_rows"`
	ReturningData []user `json:"returning"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/createcrow", createCrow)
	r.Run(":3000")
}

func createCrow(c *gin.Context) {
	query := `
	mutation CreateCrow($email: String = "", $password: String = "", $username: String = "") {
	  insert_users(objects: {email: $email, password: $password, username: $username}) {
	    affected_rows
	    returning {
	      uuid
	      username
	      email
	    }
	  }
	}`

	variables := newUserGraphQLVariables{
		"zach@sneakycrow.dev",
		"something",
		"sneakycrow",
	}

	gqlMarshalled, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://localhost:8080/v1/graphql", "application/json", bytes.NewBuffer(gqlMarshalled))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var data hasuraResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("raw body %s", string(body))
	// fmt.Printf("%+v\n", data)
	c.JSON(http.StatusOK, data)
}
