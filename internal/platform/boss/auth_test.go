package boss

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yahao333/get_jobs/internal/storage"
)

func TestCheckLoginStatus_NilPage(t *testing.T) {
	// Create a temp cookie file
	tmpFile, err := os.CreateTemp("", "cookies.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write dummy cookies
	cookies := []storage.Cookie{
		{
			Name:    "test",
			Value:   "value",
			Domain:  "example.com",
			Path:    "/",
			Expires: time.Now().Add(time.Hour),
		},
	}
	data, _ := json.Marshal(cookies)
	if err := os.WriteFile(tmpFile.Name(), data, 0644); err != nil {
		t.Fatal(err)
	}

	// Create BossClient with nil page
	client := NewBossClient(tmpFile.Name())
	// Ensure page is nil (it is by default from NewBossClient)
	assert.Nil(t, client.page)

	// Call CheckLoginStatus
	// This should not panic and return false
	status, err := client.CheckLoginStatus()

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.IsLoggedIn)
}
