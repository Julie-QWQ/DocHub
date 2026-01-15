package repository

import (
	"context"
	"testing"

	"github.com/study-upc/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.User{})
	require.NoError(t, err)

	// 清理表数据（确保每个测试从干净状态开始）
	db.Exec("DELETE FROM users")

	return db
}

func TestUserRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		RealName: "Test User",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
		Major:    "计算机科学",
		Class:    "2101",
	}

	err := repo.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 先创建用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	// 查找用户
	found, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	found, err := repo.FindByUsername(ctx, "testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)

	// 测试不存在的用户
	_, err = repo.FindByUsername(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	found, err := repo.FindByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	// 更新用户信息
	user.RealName = "Updated Name"
	user.Major = "软件工程"
	err = repo.UpdateUser(ctx, user)
	assert.NoError(t, err)

	// 验证更新
	found, _ := repo.FindByID(ctx, user.ID)
	assert.Equal(t, "Updated Name", found.RealName)
	assert.Equal(t, "软件工程", found.Major)
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	exists, err := repo.ExistsByUsername(ctx, "testuser")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUsername(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	exists, err := repo.ExistsByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}
