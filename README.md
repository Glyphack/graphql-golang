---
title: ‌Introduction To GraphQL Server With Golang
published: false
description: Introduction to GraphQL Server with Golang and Gqlgen.
tags: graphql, go, api, gqlgen
---
## Table Of Contents
- [Table Of Contents](#table-of-contents)
  - [Motivation ](#motivation)
      - [What is a GraphQL server?](#what-is-a-graphql-server)
      - [Schema-Driven Development](#schema-driven-development)
  - [Getting started ](#getting-started)
      - [Project Setup](#project-setup)
      - [Defining Our Schema](#defining-out-schema)
  - [Queries](#queries)
      - [What Is A Query](#what-is-a-query)
      - [Simple Query](#simple-query)
  - [Mutations](#mutations)
      - [What Is A Mutation](#what-is-a-mutation)
      - [A Simple Mutation](#a-simple-mutation)
  - [Database](#database)
      - [Setup MySQL](#setup-mysql)
      - [Models and migrations](#models-and-migrations)
  - [Create and Retrieve Links](#create-and-retrieve-links)
      - [CreateLinks](#createlinks)
      - [links Query](#links-query)
  - [Authentication](#authentication)
      - [JWT](#jwt)
      - [Setup](#setup)
        - [Generating and Parsing JWT Tokens](#generating-and-parsing-jwt-tokens)
        - [User Signup and Login Functionality](#user-signup-and-login-functionality)
        - [Authentication Middleware](#authentication-middleware)
  - [Continue Implementing schema](#continue-implementing-schema)
    - [CreateUser](#createuser)
    - [Login](#login)
	- [RefreshToken](#refresh-token)
    - [Completing Our App](#completing-our-app)

### Motivation <a name="motivation"></a>
[**Go**](https://golang.org/) is a modern general purpose programming language designed by google; best known for it's simplicity, concurrency and fast performance. It's being used by big players in the industry like Google, Docker, Lyft and Uber. If you are new to golang you can start from [golang tour](https://tour.golang.org/) to learn fundamentals.

[**gqlgen**](https://gqlgen.com/) is a library for creating GraphQL applications in Go.


In this tutorial we Implement a Hackernews GraphQL API clone  with *golang* and *gqlgen* and learn about GraphQL fundamentals along the way.
Source code and also this tutorial are available on Github at: https://github.com/Glyphack/go-graphql-hackernews

#### What is a GraphQL server? <a name="what-is-a-graphql-server"></a>
A GraphQL server is able to receive requests in GraphQL Query Language format and return response in desired form.
GraphQL is a query language for API so you can send queries and ask for what you need and exactly get that piece of data.
In this sample query we are looking for address, title of the links and name of the user who add it:
```
query {
	links{
    	title
    	address,
    	user{
      		name
    	}
  	}
}
```
response:
```
{
  "data": {
    "links": [
      {
        "title": "our dummy link",
        "address": "https://address.org",
        "user": {
          "name": "admin"
        }
      }
    ]
  }
}
```

#### Schema-Driven Development <a name="schema-driven-development"></a>
In GraphQL your API starts with a schema that defines all your types, queries and mutations, It helps others to understand your API. So it's like a contract between server and the client.
Whenever you need to add a new capability to a GraphQL API you must redefine schema file and then implement that part in your code. GraphQL has it's [Schema Definition Language](http://graphql.org/learn/schema/) for this purpose.
gqlgen is a Go library for building GraphQL servers and has a nice feature that generates code based on your schema definition.

### Getting started <a name="getting-started"></a>
In this tutorial we are going to create a Hackernews clone with Go and gqlgen, So our API will be able to handle registration, authentication, submitting links and getting list of links.

#### Project Setup <a name="project-setup"></a>
Create a directory for project and initialize go modules file:
```bash
go mod init github.com/[username]/hackernews
```

after that use ‍‍gqlgen `init` command to setup a gqlgen project.
```
go run github.com/99designs/gqlgen init
```
Here is a description from gqlgen about the generated files:
* `gqlgen.yml` — The gqlgen config file, knobs for controlling the generated code.
* `generated.go` — The GraphQL execution runtime, the bulk of the generated code.
* `models_gen.go` — Generated models required to build the graph. Often you will override these with your own models. Still very useful for input types.
* `resolver.go` — This is where your application code lives. generated.go will call into this to get the data the user has requested.
* `server/server.go` — This is a minimal entry point that sets up an http.Handler to the generated GraphQL server.
start the server with `go run server.go` and open your browser and you should see the graphql playground, So setup is right!

#### Defining Our Schema <a name="defining-out-schema"></a>
Now let's start with defining schema we need for our API. 
We have two types Link and User each of them for representing Link and User to client, a `links` Query to return list of Links. an input for creating new links and mutation for creating link. we also need mutations to for auth system which includes Login, createUser, refreshToken(I'll explain them later) then run the command below to regenerate graphql models.
```js
type Link {
  id: ID!
  title: String!
  address: String!
  user: User!
}

type User {
  id: ID!
  name: String!
}

type Query {
  links: [Link!]!
}

input NewLink {
  title: String!
  address: String!
}

input RefreshTokenInput{
  token: String!
}

input NewUser {
  username: String!
  password: String!
}

input Login {
  username: String!
  password: String!
}

type Mutation {
  createLink(input: NewLink!): Link!
  createUser(input: NewUser!): String!
  login(input: Login!): String!
  # we'll talk about this in authentication section
  refreshToken(input: RefreshTokenInput!): String!
}
```
Now remove resolver.go and  re-run the command to regenerate files;
```
rm resolver.go
go run github.com/99designs/gqlgen
```
After gqlgen generated code for us with have to implement our schema, we do that in ‍‍‍‍`resolver.go`, as you see there is functions for Queries and Mutations we defined in our schema.


### Queries <a name="queries"></a>
In the previous section we setup up the server, Now we try to implement a Query that we defined in `schema.grpahql`.

#### What Is A Query <a name="what-is-a-query"></a>
a query in graphql is asking for data, you use a query and specify what you want and graphql will return it back to you.

#### Simple Query <a name="simple-query"></a>
 open `resolver.go` file and take a look at Links function,
```go
func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
```
Notice that this function takes a Context and returns slice of Links and an error(is there is any).
ctx argument contains the data from the person who sends request like which user is working with app(we'll see how later), etc.

Let's make a dummy response for this function, for now.

`resolver.go`:
```go
func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	var links []*Link
	links = append(links, &Link{Title: "our dummy link", Address: "https://address.org", User: &User{Username: "admin"}})
	return links, nil
}
```
now run the server with `go run server/server.go` and send this query in graphiql:
```
query {
	links{
    title
    address,
    user{
      name
    }
  }
}
```
And you will get:
```
{
  "data": {
    "links": [
      {
        "title": "our dummy link",
        "address": "https://address.org",
        "user": {
          "name": "admin"
        }
      }
    ]
  }
}
```
Now you know how we generate response for our graphql server. But this response is just a dummy response we want be able to query all other users links, In the next section we setup database for our app to be able to save links and retrieve them from database.


### Mutations <a name="mutations"></a>
#### What Is A Mutation <a name="what-is-a-mutation"></a>
Simply mutations are just like queries but they can cause a data write, Technically Queries can be used to write data too however it's not suggested to use it.
So mutations are like queries, they have names, parameters and they can return data.
#### A Simple Mutation <a name="a-simple-mutation"></a>
Let's try to implement the createLink mutation, since we do not have a database set up yet(we'll get it done in the next section) we just receive the link data and construct a link object and send it back for response!
Open `resolver.go` and Look at `CreateLink` function:
```go
func (r *mutationResolver) CreateLink(ctx context.Context, input NewLink) (*Link, error) {
```
This function receives a `NewLink` with type of `input` we defined NewLink structure in our `schema.graphql` try to look at the structure and try Construct a `Link` object that be defined in our `schema.ghraphql`:
```go
func (r *mutationResolver) CreateLink(ctx context.Context, input NewLink) (*Link, error) {
	var link Link
	var user User
	link.Address = input.Address
	link.Title = input.Title
	user.Username = "test"
	link.User = &user
	return &link, nil
}
```
now run server and use the mutation to create a new link:
```
mutation {
  createLink(input: {title: "new link", address:"http://address.org"}){
    title,
    user{
      name
    }
    address
  }
}
```
and you will get:
```
{
  "data": {
    "createLink": {
      "title": "new link",
      "user": {
        "name": "test"
      },
      "address": "http://address.org"
    }
  }
}
```
Nice now we know what are mutations and queries we can setup our database and make these implementations more practical.


### Database <a name="database"></a>
Before we jump into implementing GraphQL schema we need to setup database to save users and links, This is not supposed to be tutorial about databases in go but here is what we are going to do:
* setup MySQL
* define our models and create migrations

##### Setup MySQL <a name="setup-mysql"></a>
If you have docker you can run [Mysql image]((https://hub.docker.com/_/mysql)) from docker and use it.
`docker run --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -d mysql:latest`
now run `docker ps` and you should see our mysql image is running:
```
CONTAINER ID        IMAGE                                                               COMMAND                  CREATED             STATUS              PORTS                  NAMES
8fea71529bb2        mysql:latest                                                        "docker-entrypoint.s…"   2 hours ago         Up 2 hours          3306/tcp, 33060/tcp    mysql

```

Now create a database for our application:
```bash
docker exec -it mysql bash
mysql -u root -p
CREATE DATABASE hackernews;
```

##### Models and migrations <a name="models-and-migrations"></a>
We need to create migrations for our app so every time our app runs it creates tables it needs to work properly, we are going to use [golang-migrate](https://github.com/golang-migrate/migrate) package.
create a folder structure for our database files:
```
go-graphql-hackernews
--internal
----pkg
------db
--------migrations
----------mysql
```
Install go mysql driver and golang-migrate packages then create migrations:
```
go get -u github.com/go-sql-driver/mysql
go build -tags 'mysql' -ldflags="-X main.Version=$(git describe --tags)" -o $GOPATH/bin/migrate github.com/golang-migrate/migrate/cmd/migrate
cd internal/pkg/db/migrations/
migrate create -ext sql -dir mysql -seq create_users_table
migrate create -ext sql -dir mysql -seq create_links_table
```
migrate command will create two files for each migration ending with .up and .down; up is responsible for applying migration and down is responsible for reversing it.
open `create_users_table.up.sql` and add table for our users:
```sql
CREATE TABLE IF NOT EXISTS Users(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Username VARCHAR (127) NOT NULL UNIQUE,
    Password VARCHAR (127) NOT NULL,
    PRIMARY KEY (ID)
)
```
in `create_links_table.up.sql`:
```sql
CREATE TABLE IF NOT EXISTS Links(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Title VARCHAR (255) ,
    Address VARCHAR (255) ,
    UserID INT ,
    FOREIGN KEY (UserID) REFERENCES Users(ID) ,
    PRIMARY KEY (ID)
)
```

We need one table for saving links and one table for saving users, Then we apply these to our database using migrate command.

```bash
  migrate -database mysql://root:dbpass@(172.17.0.2:3306)/hackernews -path internal/pkg/db/migrations/mysql up
```

Last thing is that we need a connection to our database, for this we create a mysql.go under mysql folder(We name this file after mysql since we are now using mysql and if we want to have multiple databases we can add other folders) with a function to initialize connection to database for later use.

`internal/pkg/db/mysql/mysql.go`:
```go
package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
)

var Db *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", "root:dbpass@(172.17.0.2:3306)/hackernews")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
 		log.Panic(err)
	}
	Db = db
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}
```
`InitDB` Function creates a connection to our database and `Migrate` function runs migrations file for us.
In `Migrate function we apply migrations just like we did with command line but with this function your app will always apply the latest migrations before start.

Then call `InitDB` and `Migrate`(Optional) In main func to create database connection at the start of the app:
```go
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	database.Migration()

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(hackernews.NewExecutableSchema(hackernews.Config{Resolvers: &hackernews.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

```


### Create and Retrieve Links <a name="create-and-retrieve-links"></a>
Now we have our database ready we can start implementing our schema!

#### CreateLinks <a name="createlinks"></a>
Lets implement CreateLink mutation; first we need a function to let us write a link to database.
Create a folders links and users inside internal folder, these packages are layers between database and our app.

`internal/users/users.go`:
```go
package users

type User struct {
	ID       string `json:"id"`
	Username     string `json:"name"`
	Password string `json:"password"`
}
```
`internal/links/links.go`:
```go
package links

import (
	database "github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"github.com/glyphack/go-graphql-hackernews/internal/users"
	"log"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

//#2
func (link Link) Save() int64 {
	//#3
	statement, err := database.Db.Prepare("INSERT INTO Links(Title,Address) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := statement.Exec(link.Title, link.Address)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

```
In users.go we just defined a `struct` that represent users we get from database, But let me explain links.go part by part:
* 1: definition of struct that represent a link.
* 2: function that insert a Link object into database and returns it's ID.
* 3: our sql query to insert link into Links table. you see we used prepare here before db.Exec, the prepared statements helps you with security and also performance improvement in some cases. you can read more about it [here](https://www.postgresql.org/docs/9.3/sql-prepare.html).
* 4: execution of our sql statement.
* 5: retrieving Id of inserted Link.

Now we use this function in our CreateLink resolver:

`resolver.go`:
```go
func (r *mutationResolver) CreateLink(ctx context.Context, input NewLink) (*Link, error) {
	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	linkId := link.Save()
	return &Link{ID: strconv.FormatInt(linkId, 10), Title:link.Title, Address:link.Address}, nil
}
```
Hopefully you understand this piece of code, we create a link object from input and save it to database then return newly created link(notice that we convert the ID to string with `strconv.FormatInt`).
note that here we have 2 structs for Link in our project, one is use for our graphql server and one is for our database.
run the server and open graphiql page to test what we just wrote:
```
mutation create{
  createLink(input: {title: "something", address: "somewhere"}){
    title,
    address,
    id,
  }
}
```
```
{
  "data": {
    "createLink": {
      "title": "something",
      "address": "somewhere",
      "id": "1"
    }
  }
}
```
Grate job!

#### links Query <a name="links-query"></a>
Just like how we implemented CreateLink mutation we implement links query, we need a function to retrieve links from database and pass it to graphql server in our resolver.
Create a function named GetAll

`internal/links/links.go`:
```go
func GetAll() []Link {
	stmt, err := database.Db.Prepare("select id, title, address from Links")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address)
		if err != nil{
			log.Fatal(err)
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
```

Return links from GetAll in Links query.

`resolver.go`:
```go
func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	var resultLinks []*Link
	var dbLinks []links.Link
	dbLinks = links.GetAll()
	for _, link := range dbLinks{
		resultLinks = append(resultLinks, &Link{ID:link.ID, Title:link.Title, Address:link.Address})
	}
	return resultLinks, nil
}
```
Now query Links at graphiql:
```
query {
  links {
    title
    address
    id
  }
}

```
result:
```
{
  "data": {
    "links": [
      {
        "title": "something",
        "address": "somewhere",
        "id": "1"
      }
    ]
  }
}
```

### Authentication <a name="authentication"></a>
One of most common layers in web applications is authentication system, our app is no exception. For authentication we are going to use jwt tokens as our way to authentication users, lets see how it works.

##### JWT <a name="jwt"></a>
[JWT](https://jwt.io/) or Json Web Token is a string containing a hash that helps us verify who is using application. Every token is constructed of 3 parts like `xxxxx.yyyyy.zzzzz` and name of these parts are: Header, Payload and Signature. Explanation about these parts are more about JWT than our application you can read more about them [here](https://jwt.io/introduction/).
whenever a user login to an app server generates a token for user, Usually server saves some information like username about the user in token to be able to recognize the user later using that token.This tokens get signed by a key so only the issuer app can reopen the token.
We are going to implement this behavior in our app. 

##### Setup <a name="setup"></a>
In our app we need to be able to generate a token for users when they sign up or login and a middleware to authenticate users by the given token, then in our views we can know the user interacting with app. We will be using `github.com/dgrijalva/jwt-go` library to generate and prase JWT tokens.
###### Generating and Parsing JWT Tokens <a name="generating-and-parsing-jwt-tokens"></a>
We create a new directory pkg in the root of our application, you have seen that we used internal for what we want to only be internally used withing our app, pkg directory is for files that we don't mind if some outer code imports it into itself and generation and validation jwt tokens are this kinds of code.
There is a concept named claims it's not only limited to JWT We'll see more about it in rest of the section.

`pkg/jwt/jwt.go`:
```go
package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// secret key being used to sign tokens
var (
	SecretKey = []byte("secret")
)

//GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

//ParseToken parses a jwt token and returns the username it it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
```
Let's talk about what above code does:
* GenerateToken function is going to be used whenever we want to generate a token for user, we save username in token claims and set token expire time to 5 minutes later also in claims.
* ParseToken function is going to be used whenever we receive a token and want to know who sent this token.

###### User SignUp and Login Functionality <a name="user-signup-and-login-functionality"></a>
Til now we can generate a token for each user but before generating token for every user, we need to assure user exists in our database. Simply we need to query database to match the user with given username and password.
Another thing is when a user tries to register we insert username and password in our database.

`internal/users/users.go`:
```go
package users

import (
	"database/sql"
	"github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"

	"log"
)

type User struct {
	ID       string `json:"id"`
	Username     string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```
The Create function is much like the [CreateLink](#createlinks) function we saw earlier but let's break down the Authenticate code:
* first we have a query to select password from users table where username is equal to the username we got from resolver.
* We use QueryRow instead of Exec we used earlier; the difference is `QueryRow()` will return a pointer to a `sql.Row`.
* Using `.Scan` method we fill the hashedPassword variable with the hashed password from database. Obviously you don't want to [save raw passwords](https://security.blogoverflow.com/2011/11/why-passwords-should-be-hashed/) in your database.
* then we check if any user with given username exists or not, if there is not any we return `false`, and if we found any we check the user hashedPassword with the raw password given.(Notice that we save hashed passwords not raw passwords in database in line 23)

In the next part we set the tools we have together to detect the user that is using the app.
###### Authentication Middleware <a name="authentication-middleware"></a>
Every time a request comes to our resolver before sending it to resolver we want to recognize the user sending request, for this purpose we have to write a code before every resolver, but using middleware we can have a auth middleware that executes before request send to resolver and does the authentication process. to read more about middlewares visit.


`internal/users/users.go`:
```go
//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}
```
We use this function to get user object with username in authentication middeware.

And now let's create our auth middleware, for more information visit [gql authentication docs](https://github.com/99designs/gqlgen/blob/master/docs/content/recipes/authentication.md).

`internal/auth/middleware.go`:
```go
package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/glyphack/go-graphql-hackernews/internal/users"
	"github.com/glyphack/go-graphql-hackernews/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("token")

			// Allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := c.Value
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = strconv.Itoa(id)

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
```

Now we use the middleware we declared in our server:

`server/server.go`:
```go
package main

import (
	"github.com/glyphack/go-graphql-hackernews/internal/auth"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	hackernews "github.com/glyphack/go-graphql-hackernews"
	"github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	database.InitDB()
	database.Migrate()
	server := handler.GraphQL(hackernews.NewExecutableSchema(hackernews.Config{Resolvers: &hackernews.Resolver{}}))
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
```

### Continue Implementing schema <a name="continue-implementing-schema"></a>
Now that we have working authentication system we can get back to implementing our schema.
#### CreateUser <a name="createuser"></a>
We continue our implementation of CreateUser mutation with functions we have written in auth section.

`resolver.go`:
```go
func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	err := user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil{
		return "", err
	}
	return token, nil
}
```
In our mutation first we create a user using given username and password and then generate a token for the user so we can recognize the user in requests.
Start the server and try it in graphiql:
query:
```
mutation {
  createUser(input: {username: "user1", password: "123"})
}
```
results:
```
{
  "data": {
    "createUser": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODE0NjAwODUsImlhdCI6MTU4MTQ1OTc4NX0.rYLOM123kSulGjvK5VP8c7S0kgk03WweS2VJUUbAgNA"
  }
}
```
So we successfully created our first user!

#### Login <a name="login"></a>
For this mutation, first we have to check if user exists in database and given password is correct, then we generate a token for user and give it bach to user.

`internal/users.go`:
```go
func (user *User) Authenticate() bool {
	statement, err := database.Db.Prepare("select Password from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```
Explanation:
* we select the user with the given username and then check if hash of the given password is equal to hashed password that we saved in database.

`resolver.go`
```go
func (r *mutationResolver) Login(ctx context.Context, input Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil{
		return "", err
	}
	return token, nil
}
```
We used the Authenticate function declared above and after that if the username and password are correct we return a new token for user and if not we return error, `&users.WrongUsernameOrPasswordError`, here is implementation for this error:

`internal/users/errors.go`:
```go
package users

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}
```
To define a custom error in go you need a struct with Error method implemented, here is our error for wrong username or password with it's Error() method.
Again you can try login with username and password from the user we created and get a token.

#### Refresh Token <a name="refresh-token"></a>
This is the last endpoint we need to complete our authentication system, imagine a user has loggedIn in our app and it's token is going to get expired after minutes we set(when generated the token), now we need a solution to keep our user loggedIn. One solution is to have a endpoint to get tokens that are going to expire and regenerate a new token for that user so that app uses new token.
So our endpoint should take a token, Parse the username and generate a token for that username.

`resolver.go`:
```go
func (r *mutationResolver) RefreshToken(ctx context.Context, input RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
```
Implementation is pretty straightforward so we skip the explanation for this.


#### Completing Our app <a name="completing-our-app"></a>
Our CreateLink mutation left incomplete because we could not authorize users back then, so let's get back to it and complete the implementation.
With what we did in [authentication middleware](#authentication-middleware) we can retrieve user in resolvers using ctx argument. so in CreateLink function add these lines:

`resolver.go`:
```go
func (r *mutationResolver) CreateLink(ctx context.Context, input NewLink) (*Link, error) {
	// 1
	user := auth.ForContext(ctx)
	if user == nil {
		return &Link{}, fmt.Errorf("access denied")
	}
	.
	.
	.
	// 2
	link.User = user
	linkId := link.Save()
	grahpqlUser := &User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &Link{ID: strconv.FormatInt(linkId, 10), Title:link.Title, Address:link.Address, User:grahpqlUser}, nil
}
```
Explanation:
* 1: we get user object from ctx and if user is not set we return error with message access denied.
* 2: then we set user of that link equal to the user is requesting to create the link.

And edit the links query to get user from db too.

`resolver.go`:
```go
func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	var resultLinks []*Link
	var dbLinks []links.Link
	dbLinks = links.GetAll()
	for _, link := range dbLinks{
		grahpqlUser := &User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &Link{ID:link.ID, Title:link.Title, Address:link.Address, User:grahpqlUser})
	}
	return resultLinks, nil
}
```


The part that is left here is our database operation for creating link, We need to create foreign key from the link we inserting to that user.

`internal/links/links.go`:
In our Save method from links changed the query statement to:
```go
statement, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?, ?)")
```
and the line that we execute query to:
```go
res, err := statement.Exec(link.Title, link.Address, link.User.ID)
```
Then when we query for users we also fill the `User` field for Link, so we need to join Links and Users table in our `GetAll` functions to fill the User field.
If you are not familiar with join checkout [this link](https://www.w3schools.com/sql/sql_join_inner.asp).

`internal/links/links.go`:
```go
func GetAll() []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID") // changed
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	var username string
	var id string
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username) // changed
		if err != nil{
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		} // changed
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
```

and Our app is finally complete.

## Summary
Congratulations on make it to here! You've learned about gqlgen library and some Graphql fundamentals. By implementing a HackerNews clone you've learned about queries, mutations, authentication and GraphQL query language. 
### Further Steps
If you want to create more complex GrahpQL APIs there are few things you can check out:
* Implement voting feature for hackernews clone app.
* Read about (Relays)[https://facebook.github.io/relay/docs/en/graphql-server-specification.html] and Pagination and implement them in app.
* don't forget to visit (Learn GraphQL page)[https://graphql.org/learn/].