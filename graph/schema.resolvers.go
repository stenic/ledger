package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stenic/ledger/graph/generated"
	"github.com/stenic/ledger/graph/model"
	"github.com/stenic/ledger/internal/auth"
	"github.com/stenic/ledger/internal/pkg/applications"
	"github.com/stenic/ledger/internal/pkg/environments"
	"github.com/stenic/ledger/internal/pkg/query"
	"github.com/stenic/ledger/internal/pkg/versions"
)

// CreateVersion is the resolver for the createVersion field.
func (r *mutationResolver) CreateVersion(ctx context.Context, input model.NewVersion) (*model.Version, error) {
	if user := auth.TokenFromContext(ctx); user == nil {
		return nil, fmt.Errorf("access denield")
	}
	version := versions.Version{
		Application: input.Application,
		Environment: input.Environment,
		Version:     input.Version,
	}
	versionID := version.Save()
	return &model.Version{
		ID:          strconv.FormatInt(versionID, 10),
		Application: &model.Application{Name: version.Application},
		Environment: &model.Environment{Name: version.Environment},
		Version:     version.Version,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.AuthPayload, error) {
	if username == "username" && password == "password" {
		tkn, err := auth.GenerateToken(username)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return &model.AuthPayload{
			Token: tkn,
		}, nil
	}
	return nil, fmt.Errorf("login failure")
}

// Versions is the resolver for the versions field.
func (r *queryResolver) Versions(ctx context.Context, orderBy *model.VersionOrderByInput) ([]*model.Version, error) {
	if user := auth.TokenFromContext(ctx); user == nil {
		return nil, fmt.Errorf("access denield")
	}
	order := versions.VersionOrderBy{}
	if orderBy != nil {
		if orderBy.Timestamp != nil {
			order.Timestamp = query.SortFromString(orderBy.Timestamp.String())
		}
		if orderBy.Application != nil {
			order.Application = query.SortFromString(orderBy.Application.String())
		}
		if orderBy.Environment != nil {
			order.Environment = query.SortFromString(orderBy.Environment.String())
		}
	}

	var resultLinks []*model.Version
	for _, item := range versions.GetAll(order) {
		resultLinks = append(resultLinks, &model.Version{
			ID:          strconv.FormatInt(item.ID, 10),
			Application: &model.Application{Name: item.Application},
			Environment: &model.Environment{Name: item.Environment},
			Version:     item.Version,
			Timestamp:   item.Timestamp,
		})
	}
	return resultLinks, nil
}

// Environments is the resolver for the environments field.
func (r *queryResolver) Environments(ctx context.Context) ([]*model.Environment, error) {
	if user := auth.TokenFromContext(ctx); user == nil {
		return nil, fmt.Errorf("access denield")
	}
	var resultLinks []*model.Environment
	for _, item := range environments.GetAll() {
		resultLinks = append(resultLinks, &model.Environment{
			Name: item.Name,
		})
	}
	return resultLinks, nil
}

// Applications is the resolver for the applications field.
func (r *queryResolver) Applications(ctx context.Context) ([]*model.Application, error) {
	if user := auth.TokenFromContext(ctx); user == nil {
		return nil, fmt.Errorf("access denield")
	}
	var resultLinks []*model.Application
	for _, item := range applications.GetAll() {
		resultLinks = append(resultLinks, &model.Application{
			Name: item.Name,
		})
	}
	return resultLinks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
