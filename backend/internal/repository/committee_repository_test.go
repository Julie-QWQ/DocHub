package repository

import (
	"context"
	"testing"

	"github.com/study-upc/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CommitteeRepositoryTestSuite 学委申请仓库测试套件
type CommitteeRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo CommitteeRepository
}

// SetupTest 设置测试环境
func (s *CommitteeRepositoryTestSuite) SetupTest() {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(s.T(), err)

	// 自动迁移
	err = db.AutoMigrate(&model.User{}, &model.CommitteeApplication{})
	assert.NoError(s.T(), err)

	s.db = db
	s.repo = NewCommitteeRepository(db)
}

// TearDownTest 清理测试环境
func (s *CommitteeRepositoryTestSuite) TearDownTest() {
	// 清理数据库
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

// TestCreateApplication 测试创建学委申请
func (s *CommitteeRepositoryTestSuite) TestCreateApplication() {
	ctx := context.Background()

	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user)

	// 创建申请
	application := &model.CommitteeApplication{
		UserID: user.ID,
		Status: model.ApplicationPending,
		Reason: "我想成为学委，为同学们服务",
	}

	err := s.repo.CreateApplication(ctx, application)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), application.ID)
}

// TestFindApplicationByID 测试根据ID查找申请
func (s *CommitteeRepositoryTestSuite) TestFindApplicationByID() {
	ctx := context.Background()

	// 创建测试数据
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user)

	application := &model.CommitteeApplication{
		UserID: user.ID,
		Status: model.ApplicationPending,
		Reason: "我想成为学委",
	}
	s.db.Create(application)

	// 测试查找
	found, err := s.repo.FindApplicationByID(ctx, application.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), application.ID, found.ID)
	assert.Equal(s.T(), application.Reason, found.Reason)
}

// TestListApplications 测试获取申请列表
func (s *CommitteeRepositoryTestSuite) TestListApplications() {
	ctx := context.Background()

	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user)

	// 创建多个申请
	for i := 0; i < 5; i++ {
		application := &model.CommitteeApplication{
			UserID: user.ID,
			Status: model.ApplicationPending,
			Reason: "申请理由",
		}
		s.db.Create(application)
	}

	// 测试列表查询
	applications, total, err := s.repo.ListApplications(ctx, 1, 10, nil)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(5), total)
	assert.Len(s.T(), applications, 5)
}

// TestUpdateApplicationStatus 测试更新申请状态
func (s *CommitteeRepositoryTestSuite) TestUpdateApplicationStatus() {
	ctx := context.Background()

	// 创建测试数据
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user)

	reviewer := &model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "hashed_password",
		Role:     model.RoleAdmin,
		Status:   model.StatusActive,
	}
	s.db.Create(reviewer)

	application := &model.CommitteeApplication{
		UserID: user.ID,
		Status: model.ApplicationPending,
		Reason: "我想成为学委",
	}
	s.db.Create(application)

	// 更新状态
	err := s.repo.UpdateApplicationStatus(ctx, application.ID, model.ApplicationApproved, &reviewer.ID, "审核通过")
	assert.NoError(s.T(), err)

	// 验证更新
	updated, _ := s.repo.FindApplicationByID(ctx, application.ID)
	assert.Equal(s.T(), model.ApplicationApproved, updated.Status)
	assert.Equal(s.T(), reviewer.ID, *updated.ReviewerID)
	assert.Equal(s.T(), "审核通过", updated.ReviewComment)
}

// TestExistsPendingApplication 测试检查待审核申请
func (s *CommitteeRepositoryTestSuite) TestExistsPendingApplication() {
	ctx := context.Background()

	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user)

	// 创建待审核申请
	application := &model.CommitteeApplication{
		UserID: user.ID,
		Status: model.ApplicationPending,
		Reason: "我想成为学委",
	}
	s.db.Create(application)

	// 测试存在性检查
	exists, err := s.repo.ExistsPendingApplication(ctx, user.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)

	// 创建另一个用户（无待审核申请）
	user2 := &model.User{
		Username: "testuser2",
		Email:    "test2@example.com",
		Password: "hashed_password",
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	s.db.Create(user2)

	exists2, err := s.repo.ExistsPendingApplication(ctx, user2.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists2)
}

// TestCountPendingApplications 测试统计待审核申请
func (s *CommitteeRepositoryTestSuite) TestCountPendingApplications() {
	ctx := context.Background()

	// 创建测试用户
	for i := 0; i < 3; i++ {
		user := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     model.RoleStudent,
			Status:   model.StatusActive,
		}
		s.db.Create(user)

		application := &model.CommitteeApplication{
			UserID: user.ID,
			Status: model.ApplicationPending,
			Reason: "我想成为学委",
		}
		s.db.Create(application)
	}

	count, err := s.repo.CountPendingApplications(ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(3), count)
}

// 运行测试套件
func TestCommitteeRepositorySuite(t *testing.T) {
	suite.Run(t, new(CommitteeRepositoryTestSuite))
}
