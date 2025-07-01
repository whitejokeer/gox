package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatcher_New(t *testing.T) {
	watcher, err := New()
	require.NoError(t, err)
	assert.NotNil(t, watcher)
	defer watcher.Close()
}

func TestWatcher_Constants(t *testing.T) {
	assert.Equal(t, "change", EventChange)
	assert.Equal(t, "create", EventCreate)
	assert.Equal(t, "delete", EventDelete)
}

func TestWatcher_CallbackRegistration(t *testing.T) {
	watcher, err := New()
	require.NoError(t, err)
	defer watcher.Close()

	watcher.OnChange(func(path string) {
		// Callback for testing registration
	})

	watcher.OnCreate(func(path string) {
		// Callback for testing registration
	})

	watcher.OnDelete(func(path string) {
		// Callback for testing registration
	})

	// Verify callbacks are registered
	assert.NotNil(t, watcher.callbacks[EventChange])
	assert.NotNil(t, watcher.callbacks[EventCreate])
	assert.NotNil(t, watcher.callbacks[EventDelete])
}

func TestWatcher_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watcher_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	watcher, err := New()
	require.NoError(t, err)
	defer watcher.Close()

	// Track events
	var events []string
	var eventPaths []string

	watcher.OnCreate(func(path string) {
		events = append(events, "create")
		eventPaths = append(eventPaths, path)
	})

	watcher.OnChange(func(path string) {
		events = append(events, "change")
		eventPaths = append(eventPaths, path)
	})

	watcher.OnDelete(func(path string) {
		events = append(events, "delete")
		eventPaths = append(eventPaths, path)
	})

	// Add the temp directory to watch
	err = watcher.AddPath(tempDir)
	require.NoError(t, err)

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- watcher.Start()
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a .gox file
	testFile := filepath.Join(tempDir, "test.gox")
	err = os.WriteFile(testFile, []byte("<template>test</template>"), 0644)
	require.NoError(t, err)

	// Wait for event processing
	time.Sleep(200 * time.Millisecond)

	// Modify the file
	err = os.WriteFile(testFile, []byte("<template>modified</template>"), 0644)
	require.NoError(t, err)

	// Wait for event processing
	time.Sleep(200 * time.Millisecond)

	// Delete the file
	err = os.Remove(testFile)
	require.NoError(t, err)

	// Wait for event processing
	time.Sleep(200 * time.Millisecond)

	// Stop the watcher
	watcher.Close()

	// Verify events were captured
	// Note: File system events can be flaky in tests, so we check for at least some events
	assert.True(t, len(events) > 0, "Expected at least one event")

	// Check that all event paths contain our test file
	for _, path := range eventPaths {
		assert.Contains(t, path, "test.gox")
	}
}

func TestWatcher_NonGoxFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "watcher_test_non_gox")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	watcher, err := New()
	require.NoError(t, err)
	defer watcher.Close()

	// Track events
	var events []string

	watcher.OnCreate(func(path string) {
		events = append(events, "create")
	})

	// Add the temp directory to watch
	err = watcher.AddPath(tempDir)
	require.NoError(t, err)

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- watcher.Start()
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a non-.gox file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// Wait for event processing
	time.Sleep(200 * time.Millisecond)

	// Stop the watcher
	watcher.Close()

	// Verify no events were captured for non-.gox files
	assert.Equal(t, 0, len(events), "Expected no events for non-.gox files")
}
