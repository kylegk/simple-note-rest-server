# Simple Notes REST Server

## Description

This project provides the backend for a simple multi-user note application. The application allows users the ability to view, create, modify, and delete plain text notes. 

As this application is limited in scope, so is the functionality the users can perform. The ability to perform any of the actions described above is limited to the owner of the note (e.g. a user cannot view or act upon a note owned by another user).  
 
While the majority of this project is original code, it does make use of two third-party libraries: [go-membdb](https://github.com/hashicorp/go-memdb) an in-memory database solution created by HashiCorp, and [jwt-go](https://github.com/golang-jwt/jwt) a Golang implementation of JSON Web Tokens.

## Schema

Like the functionality the application provides, the schema for the in-memory data store is also very simple. There are only three tables employed at this time:
1. **USER** contains details about the user. It stores the user's id and username.
2. **NOTES** stores the content of the note, and when the note was created or last updated.
3. **USER_NOTES** is the relationship between the user and the notes they own.

The names of these tables and their associated indexes can be found in: `db/schema.go`

## Methods

All of the methods below, except for user creation require a valid token. This token is generated when creating a new user. At this time, the token is generated with claims, including an expiration value. However, the token validation currently ignores the expiration time. Future iterations of the project would add the expiration validation where noted in the code, provide a means to store tokens, and refresh tokens on demand. 

**Create A User**

```
/users
```

> Method: **POST**

> Adds a new user to the data store and returns a valid token. User accounts must be created before performing any other requests associated with notes. **NOTE**: Username must be unique.

> `Request:`

```
{
        "user": "test.account",
}
```

> `Response:`

```
{
    "userid":1,
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzIjoxNjI4NDYyMzMwLCJ1c2VyaWQiOjF9.LGte1UTgmzCg8L_FOdDPY7YsSlBeQQs3QyZDe2A7kNY"
}
```

**Create A Note**

```
/notes
```

> Method: **POST**

> Adds a new note and creates a relationship between the user account and the note. Requires a valid auth token for the user.

> `Request:`

```
{
        "content": "This is a note to be created",
}
```

> `Response:`

```
{
        "noteid": 1
}
```

**Update A Note**

```
/notes/{id}
```

> Method: **POST**

> Updates the content of the note specified. Requires a valid auth token for the user (i.e. the note must be owned by the user performing the update).

> `Request:`

```
{
        "content": "This is an updated note",
}
```

> `Response:`

```
{
        "Message": "Note updated"
}
```

**Delete A Note**

```
/notes/{id}
```

> Method: **DELETE**

> Deletes a note and the relationship between the user and the note. Requires a valid auth token for the user (i.e. the note must be owned by the user performing the deletion).

> `Response:`

```
{
        "Message": "Note deleted"
}
```

**Get A Note**

```
/notes/{id}
```

> Method: **GET**

> Retrieve the content of a note. Requires a valid auth token for the user (i.e. the note must be owned by the user).

> `Response:`

```
{
    "noteid": 1,
    "content": "This is the content of the note",
    "modified": "2009-11-10 23:00:00 +0000 UTC m=+0.000000001"
}
```

**Get All Notes**

```
/notes
```

> Method: **GET**

> Retrieves a list of note ids that are owned by the user. Requires a valid auth token.

> `Response:`

```
{
    "notes": [1,2,55,999]
}
```

## Getting Started

This project can either be built manually or run in a Docker container.

### Building manually

Building the project manually requires that you have a recent version of Golang installed on your system.

To build the project manually, perform the following steps:

```
cd /path/to/project/
go build -o bin/notes
cd bin/
./notes
```

### Docker

This project includes a Dockerfile for easy building and containerization of the application. 

To build the Docker container, perform the following steps:

```
cd /path/to/project
docker build -t notes-rest-server .
docker run -d -p 8080:8080 notes-rest-server
```

Once the container has started, you can connect to it and view the output of any logs in the application, such as requests and errors.