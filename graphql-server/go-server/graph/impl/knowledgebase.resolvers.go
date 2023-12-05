package impl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/kubeagi/arcadia/api/base/v1alpha1"
	"github.com/kubeagi/arcadia/graphql-server/go-server/graph/generated"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/auth"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/client"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/knowledgebase"
	"github.com/kubeagi/arcadia/pkg/config"
)

// CreateKnowledgeBase is the resolver for the createKnowledgeBase field.
func (r *knowledgeBaseMutationResolver) CreateKnowledgeBase(ctx context.Context, obj *generated.KnowledgeBaseMutation, input generated.CreateKnowledgeBaseInput) (*generated.KnowledgeBase, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}

	var filegroups []v1alpha1.FileGroup
	var vectorstore v1alpha1.TypedObjectReference
	vector, _ := config.GetVectorStore(ctx, c)
	displayname, description, embedder := "", "", ""
	if input.DisplayName != nil {
		displayname = *input.DisplayName
	}
	if input.Description != nil {
		description = *input.Description
	}
	if input.VectorStore != nil {
		vectorstore = v1alpha1.TypedObjectReference(*input.VectorStore)
	} else {
		vectorstore = *vector
	}
	if input.Embedder != "" {
		embedder = input.Embedder
	}
	if input.FileGroups != nil {
		for _, f := range input.FileGroups {
			filegroup := v1alpha1.FileGroup{
				Source: (*v1alpha1.TypedObjectReference)(&f.Source),
				Paths:  f.Path,
			}
			filegroups = append(filegroups, filegroup)
		}
	}
	return knowledgebase.CreateKnowledgeBase(ctx, c, input.Name, input.Namespace, displayname, description, embedder, vectorstore, filegroups)
}

// UpdateKnowledgeBase is the resolver for the updateKnowledgeBase field.
func (r *knowledgeBaseMutationResolver) UpdateKnowledgeBase(ctx context.Context, obj *generated.KnowledgeBaseMutation, input *generated.UpdateKnowledgeBaseInput) (*generated.KnowledgeBase, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	return knowledgebase.UpdateKnowledgeBase(ctx, c, input)
}

// DeleteKnowledgeBase is the resolver for the deleteKnowledgeBase field.
func (r *knowledgeBaseMutationResolver) DeleteKnowledgeBase(ctx context.Context, obj *generated.KnowledgeBaseMutation, input *generated.DeleteCommonInput) (*string, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name := ""
	labelSelector, fieldSelector := "", ""
	if input.Name != nil {
		name = *input.Name
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	return knowledgebase.DeleteKnowledgeBase(ctx, c, name, input.Namespace, labelSelector, fieldSelector)
}

// GetKnowledgeBase is the resolver for the getKnowledgeBase field.
func (r *knowledgeBaseQueryResolver) GetKnowledgeBase(ctx context.Context, obj *generated.KnowledgeBaseQuery, name string, namespace string) (*generated.KnowledgeBase, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	return knowledgebase.ReadKnowledgeBase(ctx, c, name, namespace)
}

// ListKnowledgeBases is the resolver for the listKnowledgeBases field.
func (r *knowledgeBaseQueryResolver) ListKnowledgeBases(ctx context.Context, obj *generated.KnowledgeBaseQuery, input generated.ListKnowledgeBaseInput) (*generated.PaginatedResult, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}

	return knowledgebase.ListKnowledgeBases(ctx, c, input)
}

// Datasource is the resolver for the Datasource field.
func (r *mutationResolver) KnowledgeBase(ctx context.Context) (*generated.KnowledgeBaseMutation, error) {
	return &generated.KnowledgeBaseMutation{}, nil
}

// KnowledgeBase is the resolver for the KnowledgeBase field.
func (r *queryResolver) KnowledgeBase(ctx context.Context) (*generated.KnowledgeBaseQuery, error) {
	return &generated.KnowledgeBaseQuery{}, nil
}

// KnowledgeBaseMutation returns generated.KnowledgeBaseMutationResolver implementation.
func (r *Resolver) KnowledgeBaseMutation() generated.KnowledgeBaseMutationResolver {
	return &knowledgeBaseMutationResolver{r}
}

// KnowledgeBaseQuery returns generated.KnowledgeBaseQueryResolver implementation.
func (r *Resolver) KnowledgeBaseQuery() generated.KnowledgeBaseQueryResolver {
	return &knowledgeBaseQueryResolver{r}
}

type knowledgeBaseMutationResolver struct{ *Resolver }
type knowledgeBaseQueryResolver struct{ *Resolver }
