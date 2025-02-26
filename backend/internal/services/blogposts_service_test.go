package services_test

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"errors"
	"testing"
	"time"
)

// MockBlogpostRepository implements the repository.BlogpostRepository interface
type MockBlogpostRepository struct {
	blogposts       map[int64]*models.Blogpost
	blogpostsBySlug map[string]*models.Blogpost
	nextID          int64
	getBlogpostsErr error
	getBySlugErr    error
	createErr       error
}

func NewMockBlogpostRepository() repository.BlogpostRepository {
	return &MockBlogpostRepository{
		blogposts:       make(map[int64]*models.Blogpost),
		blogpostsBySlug: make(map[string]*models.Blogpost),
		nextID:          1,
	}
}

func (m *MockBlogpostRepository) GetBlogposts(ctx context.Context) ([]*models.Blogpost, error) {
	if m.getBlogpostsErr != nil {
		return nil, m.getBlogpostsErr
	}
	var result []*models.Blogpost
	for _, bp := range m.blogposts {
		result = append(result, bp)
	}
	return result, nil
}

func (m *MockBlogpostRepository) GetBlogpostBySlug(ctx context.Context, slug string) (*models.Blogpost, error) {
	if m.getBySlugErr != nil {
		return nil, m.getBySlugErr
	}
	bp, exists := m.blogpostsBySlug[slug]
	if !exists {
		return nil, errors.New("blogpost not found")
	}
	return bp, nil
}

func (m *MockBlogpostRepository) CreateBlogpost(ctx context.Context, post *models.Blogpost) error {
	if m.createErr != nil {
		return m.createErr
	}
	post.ID = m.nextID
	m.nextID++
	post.CreatedAt = time.Now()
	m.blogposts[post.ID] = post
	m.blogpostsBySlug[post.Slug] = post
	return nil
}

// Helper function to access the underlying mock
func getMockBlogpostRepo(repo repository.BlogpostRepository) *MockBlogpostRepository {
	return repo.(*MockBlogpostRepository)
}

// Test for GetBlogposts
func TestGetBlogposts(t *testing.T) {
	// Setup
	mockRepo := NewMockBlogpostRepository()
	service := services.NewBlogpostService(mockRepo)
	ctx := context.Background()

	// Create test blogposts
	post1 := &models.Blogpost{
		Title:   "Test Post 1",
		Post:    "Content of test post 1",
		Summary: "Summary 1",
		Slug:    "test-post-1",
	}
	post2 := &models.Blogpost{
		Title:   "Test Post 2",
		Post:    "Content of test post 2",
		Summary: "Summary 2",
		Slug:    "test-post-2",
	}
	getMockBlogpostRepo(mockRepo).CreateBlogpost(ctx, post1)
	getMockBlogpostRepo(mockRepo).CreateBlogpost(ctx, post2)

	// Test successful retrieval
	posts, err := service.GetBlogposts(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	// Test error case
	getMockBlogpostRepo(mockRepo).getBlogpostsErr = errors.New("database error")
	_, err = service.GetBlogposts(ctx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for GetBlogpostBySlug
func TestGetBlogpostBySlug(t *testing.T) {
	// Setup
	mockRepo := NewMockBlogpostRepository()
	service := services.NewBlogpostService(mockRepo)
	ctx := context.Background()

	// Create a test blogpost
	testPost := &models.Blogpost{
		Title:   "Test Post",
		Post:    "Content of test post",
		Summary: "Summary",
		Slug:    "test-post",
	}
	getMockBlogpostRepo(mockRepo).CreateBlogpost(ctx, testPost)

	// Test successful retrieval
	post, err := service.GetBlogpostBySlug(ctx, testPost.Slug)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if post.ID != testPost.ID {
		t.Errorf("Expected post ID %d, got %d", testPost.ID, post.ID)
	}
	if post.Title != testPost.Title {
		t.Errorf("Expected title %s, got %s", testPost.Title, post.Title)
	}
	if post.Slug != testPost.Slug {
		t.Errorf("Expected slug %s, got %s", testPost.Slug, post.Slug)
	}

	// Test post not found
	_, err = service.GetBlogpostBySlug(ctx, "nonexistent-slug")
	if err == nil {
		t.Error("Expected error for nonexistent slug, got nil")
	}

	// Test error case
	getMockBlogpostRepo(mockRepo).getBySlugErr = errors.New("database error")
	_, err = service.GetBlogpostBySlug(ctx, testPost.Slug)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// Test for CreateBlogpost
func TestCreateBlogpost(t *testing.T) {
	// Setup
	mockRepo := NewMockBlogpostRepository()
	service := services.NewBlogpostService(mockRepo)
	ctx := context.Background()

	// Test successful creation
	testPost := &models.Blogpost{
		Title:   "New Post",
		Post:    "Content of new post",
		Summary: "Summary",
		Slug:    "new-post",
	}
	err := service.CreateBlogpost(ctx, testPost)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if testPost.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", testPost.ID)
	}

	// Verify post was stored correctly
	storedPost, err := getMockBlogpostRepo(mockRepo).GetBlogpostBySlug(ctx, testPost.Slug)
	if err != nil {
		t.Errorf("Expected no error retrieving stored post, got %v", err)
	}
	if storedPost.Title != testPost.Title {
		t.Errorf("Expected title %s, got %s", testPost.Title, storedPost.Title)
	}

	// Test error case
	getMockBlogpostRepo(mockRepo).createErr = errors.New("database error")
	err = service.CreateBlogpost(ctx, testPost)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}