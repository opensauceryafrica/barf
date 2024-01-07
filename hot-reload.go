package barf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	logger "github.com/opensaucerer/barf/log"
	"github.com/opensaucerer/barf/typing"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Sentinel struct {
	config  *typing.HotReload
	startCh chan bool
	stopCh  chan bool
	watcher *fsnotify.Watcher
}

func NewSentinel(config *typing.HotReload) *Sentinel {
	watcher, _ := fsnotify.NewWatcher()
	return &Sentinel{
		config:  config,
		startCh: make(chan bool, 1),
		stopCh:  make(chan bool, 1),
		watcher: watcher,
	}
}

func (s *Sentinel) Start() error {
	if err := s.Patrol(); err != nil {
		return err
	}

	for {
		select {
		case <-restartCh:
			logger.Info("stop channel is blocked")
		case event, ok := <-s.watcher.Events:
			if !ok {
				return nil
			}
			if event.Op == fsnotify.Chmod {
				logger.Info(fmt.Sprintf("received event [%v], skipping...", event.Op))
			} else {
				logger.Info(fmt.Sprintf("received event [%v], restarting server...", event.Op))
				restartCh <- true
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return nil
			}
			return err

		}
	}
}

func (s *Sentinel) Patrol() error {
	root := s.config.GetRoot()
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			if s.config.IsIncludeFile(s.getFilename(path)) {
				logger.Info(fmt.Sprintf("watching; %v", path))
				return s.watcher.Add(path)
			}
			return nil
		}

		if s.config.IsExcludeDir(path) {
			logger.Warn(fmt.Sprintf("excluding; %v", path))
			return filepath.SkipDir
		}

		if s.isHiddenFile(s.getFilename(path)) {
			logger.Warn(fmt.Sprintf("excluding; %v", path))
			return filepath.SkipDir
		}

		if s.config.IsIncludeDir(path) {
			logger.Info(fmt.Sprintf("watching; %v", path))
			return s.watcher.Add(path)
		}
		return nil
	})
}

func (s *Sentinel) getFilename(filepath string) string {
	return path.Base(filepath)
}

func (s *Sentinel) isHiddenFile(filename string) bool {
	return len(filename) > 1 && strings.HasPrefix(filepath.Base(filename), ".")
}
