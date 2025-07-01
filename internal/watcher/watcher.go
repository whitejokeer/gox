// Package watcher provides file watching functionality for .gox files
package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Event constants for watcher callbacks
const (
	EventChange = "change"
	EventCreate = "create"
	EventDelete = "delete"
)

// Watcher handles file system watching for .gox files
type Watcher struct {
	watcher   *fsnotify.Watcher
	callbacks map[string]func(string)
}

// New creates a new file watcher
func New() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	return &Watcher{
		watcher:   fsWatcher,
		callbacks: make(map[string]func(string)),
	}, nil
}

// AddPath adds a path to watch for .gox file changes
func (w *Watcher) AddPath(path string) error {
	return w.watcher.Add(path)
}

// OnChange sets a callback for when .gox files change
func (w *Watcher) OnChange(callback func(string)) {
	w.callbacks[EventChange] = callback
}

// OnCreate sets a callback for when .gox files are created
func (w *Watcher) OnCreate(callback func(string)) {
	w.callbacks[EventCreate] = callback
}

// OnDelete sets a callback for when .gox files are deleted
func (w *Watcher) OnDelete(callback func(string)) {
	w.callbacks[EventDelete] = callback
}

// Start starts the file watcher
func (w *Watcher) Start() error {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return fmt.Errorf("watcher events channel closed")
			}

			// Only process .gox files
			if filepath.Ext(event.Name) != ".gox" {
				continue
			}

			switch {
			case event.Has(fsnotify.Write):
				if callback, exists := w.callbacks[EventChange]; exists {
					callback(event.Name)
				}
			case event.Has(fsnotify.Create):
				if callback, exists := w.callbacks[EventCreate]; exists {
					callback(event.Name)
				}
			case event.Has(fsnotify.Remove):
				if callback, exists := w.callbacks[EventDelete]; exists {
					callback(event.Name)
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher errors channel closed")
			}
			return fmt.Errorf("watcher error: %w", err)
		}
	}
}

// Close closes the file watcher
func (w *Watcher) Close() error {
	return w.watcher.Close()
}
