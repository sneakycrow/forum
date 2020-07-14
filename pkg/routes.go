package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type newUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

type user struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func CreateUser(c *gin.Context) {
	type graphQLRequest struct {
		Query     string  `json:"query"`
		Variables newUser `json:"variables"`
	}
	var submittedData newUser
	err := c.ShouldBindJSON(&submittedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	gqlMarshalled, err := json.Marshal(graphQLRequest{Query: InsertUsersQuery, Variables: submittedData})

	resp, err := http.Post("http://localhost:8080/v1/graphql", "application/json", bytes.NewBuffer(gqlMarshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Could not send data to database"})
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var respData hasuraResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Could not unmarshall DB response"})
	}
	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("New user %s successfully created", respData.Data.InsertUsers.ReturningData[0].Username)})
}
