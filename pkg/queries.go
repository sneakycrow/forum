package main

const InsertUsersQuery = `
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
