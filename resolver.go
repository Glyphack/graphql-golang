package hackernews_graphql_go

import (
	"context"
	"fmt"
	"github.com/glyphack/go-graphql-hackernews/internal/auth"
	"github.com/glyphack/go-graphql-hackernews/internal/links"
	"github.com/glyphack/go-graphql-hackernews/internal/users"
	"github.com/glyphack/go-graphql-hackernews/pkg/jwt"
	"strconv"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateLink(ctx context.Context, input NewLink) (*Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &Link{}, fmt.Errorf("access denied")
	}
	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	link.User = user
	linkId := link.Save()
	grahpqlUser := &User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &Link{ID: strconv.FormatInt(linkId, 10), Title:link.Title, Address:link.Address, User:grahpqlUser}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil{
		return "", err
	}
	return token, nil
}

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

type queryResolver struct{ *Resolver }

func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	var resultLinks []*Link
	var dbLinks []links.Link
	dbLinks = links.GetAll()
	for _, link := range dbLinks{
		grahpqlUser := &User{
			Name: link.User.Password,
		}
		resultLinks = append(resultLinks, &Link{ID:link.ID, Title:link.Title, Address:link.Address, User:grahpqlUser})
	}
	return resultLinks, nil
}
