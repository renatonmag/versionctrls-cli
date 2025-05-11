package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/helshabini/fsbroker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BranchExists(branchName string) bool {
	args := m.Called(branchName)
	return args.Bool(0)
}

func (m *MockRepository) CreateBranch(branchName, baseBranch string) error {
	args := m.Called(branchName, baseBranch)
	return args.Error(0)
}

func (m *MockRepository) CheckoutBranch(branchName string) error {
	args := m.Called(branchName)
	return args.Error(0)
}

func (m *MockRepository) CommitToBranch(branchName, filePath, message string) (string, error) {
	args := m.Called(branchName, filePath, message)
	return args.String(0), args.Error(1)
}

// MockFSService is a mock implementation of the FSService.
type MockFSService struct {
	mock.Mock
}

func (m *MockFSService) CreateHardlink(src, dest string) error {
	args := m.Called(src, dest)
	return args.Error(0)
}

func (m *MockFSService) NewFsService() FSService {
	args := m.Called()
	return args.Get(0).(FSService)
}

func (m *MockFSService) Replicate() Replicate {
	args := m.Called()
	return args.Get(0).(Replicate)
}

// Helper functions (can be in a separate test helper file)
func createTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

func cleanupTempDir(t *testing.T, dir string) {
	t.Helper()
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("failed to remove temp dir: %v", err)
	}
}

func TestOnCreate(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name              string
		branchExists      bool
		createBranchErr   error
		checkoutErr       error
		createHardlinkErr error
		commitErr         error
		expectedErr       error
		setupMocks        func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string)
		event             *fsbroker.FSEvent
	}{
		{
			name:         "Successful file creation and commit",
			branchExists: false,
			setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
				repo.On("BranchExists", "feature/file1").Return(false)
				repo.On("CreateBranch", "feature/file1", "master").Return(nil)
				repo.On("CheckoutBranch", "feature/file1").Return(nil)
				fsService.On("Replicate.CreateHardlink", event.Path, filepath.Join(repoPath, "base")).Return(nil)
				repo.On("CommitToBranch", "feature/file1", filepath.Join("base", event.Path), fmt.Sprintf("File created: %s", event.Path)).Return("commit_hash", nil)
			},
			event: &fsbroker.FSEvent{
				Path:  "file1",
				Event: fsbroker.Create,
			},
		},
		{
			name:         "Branch already exists",
			branchExists: true,
			setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
				repo.On("BranchExists", "feature/file1").Return(true)
				repo.On("CheckoutBranch", "feature/file1").Return(nil)
				fsService.On("Replicate.CreateHardlink", event.Path, filepath.Join(repoPath, "base")).Return(nil)
				repo.On("CommitToBranch", "feature/file1", filepath.Join("base", event.Path), fmt.Sprintf("File created: %s", event.Path)).Return("commit_hash", nil)
			},
			event: &fsbroker.FSEvent{
				Path:  "file1",
				Event: fsbroker.Create,
			},
		},
		{
			name:            "Create branch fails",
			branchExists:    false,
			createBranchErr: errors.New("failed to create branch"),
			expectedErr:     errors.New("failed to create branch"),
			setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
				repo.On("BranchExists", "feature/file1").Return(false)
				repo.On("CreateBranch", "feature/file1", "master").Return(errors.New("failed to create branch"))
			},
			event: &fsbroker.FSEvent{
				Path:  "file1",
				Event: fsbroker.Create,
			},
		},
		{
			name:            "Checkout branch fails",
			branchExists:    false,
			createBranchErr: nil,
			checkoutErr:     errors.New("failed to checkout branch"),
			expectedErr:     errors.New("failed to checkout branch"),
			setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
				repo.On("BranchExists", "feature/file1").Return(false)
				repo.On("CreateBranch", "feature/file1", "master").Return(nil)
				repo.On("CheckoutBranch", "feature/file1").Return(errors.New("failed to checkout branch"))
			},
			event: &fsbroker.FSEvent{
				Path:  "file1",
				Event: fsbroker.Create,
			},
			{
				name:              "Create hardlink fails",
				branchExists:      false,
				createBranchErr:   nil,
				checkoutErr:       nil,
				createHardlinkErr: errors.New("failed to create hardlink"),
				expectedErr:       errors.New("failed to create hardlink"),
				setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
					repo.On("BranchExists", "feature/file1").Return(false)
					repo.On("CreateBranch", "feature/file1", "master").Return(nil)
					repo.On("CheckoutBranch", "feature/file1").Return(nil)
					fsService.On("Replicate.CreateHardlink", event.Path, filepath.Join(repoPath, "base")).Return(errors.New("failed to create hardlink"))
				},
				event: &fsbroker.FSEvent{
					Path:  "file1",
					Event: fsbroker.Create,
				},
			},
			{
				name:              "Commit fails",
				branchExists:      false,
				createBranchErr:   nil,
				checkoutErr:       nil,
				createHardlinkErr: nil,
				commitErr:         errors.New("commit failed"),
				expectedErr:       errors.New("commit failed"),
				setupMocks: func(repo *MockRepository, fsService *MockFSService, event *fsbroker.FSEvent, repoPath string) {
					repo.On("BranchExists", "feature/file1").Return(false)
					repo.On("CreateBranch", "feature/file1", "master").Return(nil)
					repo.On("CheckoutBranch", "feature/file1").Return(nil)
					fsService.On("Replicate.CreateHardlink", event.Path, filepath.Join(repoPath, "base")).Return(nil)
					repo.On("CommitToBranch", "feature/file1", filepath.Join("base", event.Path), fmt.Sprintf("File created: %s", event.Path)).Return("", errors.New("commit failed"))
				},
				event: &fsbroker.FSEvent{
					Path:  "file1",
					Event: fsbroker.Create,
				},
			},
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := new(MockRepository)
			fsService := new(MockFSService)
			tempDir := createTempDir(t)
			defer cleanupTempDir(t, tempDir)

			// Mock the fs.NewFsService function
			// fs.NewFsService = func(config ...string) FSService {
			//  return fsService
			// }

			// Setup mock expectations
			if tc.setupMocks != nil {
				tc.setupMocks(repo, fsService, tc.event, tempDir)
			}

			// Define a mock FSService that returns the MockReplicate
			// mockFsService := new(MockFSService)
			mockReplicate := new(MockFSService)

			// Mock the behavior of Replicate() to return the mockReplicate object
			fsService.On("Replicate").Return(mockReplicate)

			// Run the test
			err := OnCreate(&Repository{Path: tempDir}, tc.event)

			// Assertions
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Assert that all expectations were met
			repo.AssertExpectations(t)
			fsService.AssertExpectations(t)
			mockReplicate.AssertExpectations(t)
		})
	}
}

// Mock implementations and helper functions

type Repository struct {
	Path string
}

func (r *Repository) BranchExists(branchName string) bool {
	return false
}

func (r *Repository) CreateBranch(branchName, baseBranch string) error {
	return nil
}

func (r *Repository) CheckoutBranch(branchName string) error {
	return nil
}

func (r *Repository) CommitToBranch(branchName, filePath, message string) (string, error) {
	return "hash", nil
}
