package envloader_test

import (
	"os"
	"testing"

	"github.com/cleopatrio/mockado/pkg/envloader"
	"github.com/cleopatrio/mockado/pkg/redactor"
)

func Test_parseTemplateVariable(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		variable string
		want     string
	}{
		{
			name:     "0",
			variable: "{{ env.REPOSITORY }}",
			want:     "/home/root/myrepo",
		},
		{
			name:     "1",
			variable: "{{ env.ABC }}",
			want:     "/var/abc",
		},
		{
			name:     "2",
			variable: "{{ env.XYZ }}",
			want:     "/tmp/random",
		},
		{
			name:     "3",
			variable: "{{ env.DB_PASSWORD }}", // Part of the default block list
			want:     "****",
		},
		{
			name:     "4",
			variable: "{{ env.GIT_CREDS }}", // Part of the default block list
			want:     "****",
		},
		{
			name:     "5",
			variable: "{{ env.ACCESS_KEY_ID }}", // Part of the default block list
			want:     "****",
		},
	}

	for _, tt := range tests {
		os.Setenv("ABC", "/var/abc")
		os.Setenv("XYZ", "/tmp/random")
		os.Setenv("REPOSITORY", "/home/root/myrepo")

		loader := envloader.WithRedactor(
			redactor.WithBlockList([]byte("CRED"), []byte("DB")),
		)

		t.Run(tt.name, func(t *testing.T) {
			got := loader.LoadVar(tt.variable)
			if got != tt.want {
				t.Errorf("parseTemplateVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestUserOperations(t *testing.T) {
	// Setup code for all tests
	db := setupTestDatabase()
	defer cleanupTestDatabase(db)

	// Subtest 1
	t.Run("CreateUser", func(t *testing.T) {
		user := createTestUser(db, "alice")
		if user.ID == 0 {
			t.Errorf("Expected user ID to be set, got 0")
		}
	})

	// Subtest 2
	t.Run("GetUser", func(t *testing.T) {
		user := createTestUser(db, "bob")
		fetchedUser, err := getUser(db, user.ID)
		if err != nil {
			t.Errorf("Failed to get user: %v", err)
		}
		if fetchedUser.Name != "bob" {
			t.Errorf("Expected name 'bob', got '%s'", fetchedUser.Name)
		}
	})

	// Subtest 3
	t.Run("DeleteUser", func(t *testing.T) {
		user := createTestUser(db, "charlie")
		err := deleteUser(db, user.ID)
		if err != nil {
			t.Errorf("Failed to delete user: %v", err)
		}

		_, err = getUser(db, user.ID)
		if err == nil {
			t.Errorf("Expected error when getting deleted user, got nil")
		}
	})
}
*/
