package services

import (
	"context"
	"backend/internal/models"
	"backend/internal/repository"
)

type BlogpostService struct {
	repo repository.BlogpostRepository
}

func NewBlogpostService(repo repository.BlogpostRepository) *BlogpostService {
	return &BlogpostService{repo: repo}
}

func (b *BlogpostService) GetBlogposts(ctx context.Context) ([]*models.Blogpost, error) {
	return b.repo.GetBlogposts(ctx)
}

func (b *BlogpostService) GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error) {
	return b.repo.GetBlogpostBySlug(ctx, slug)
}

func (b *BlogpostService) CreateBlogpost(ctx context.Context, post *models.Blogpost) error {
	return b.repo.CreateBlogpost(ctx, post)
}
