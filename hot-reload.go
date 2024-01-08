package barf

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/opensaucerer/barf/constant"
	logger "github.com/opensaucerer/barf/log"
	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Sentinel struct {
	config    *typing.HotReload
	restartCh chan bool
	watcher   *fsnotify.Watcher
}

func NewSentinel(config *typing.HotReload) *Sentinel {
	watcher, _ := fsnotify.NewWatcher()
	return &Sentinel{
		config:    config,
		restartCh: make(chan bool),
		watcher:   watcher,
	}
}

func (s *Sentinel) Start() error {
	if err := s.patrol(); err != nil {
		return err
	}

	go func() {
		<-s.restartCh
		s.restartProcess()
	}()

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return nil
			}
			if event.Op != fsnotify.Chmod {
				s.restartCh <- true
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return nil
			}
			return err
		}
	}
}

func (s *Sentinel) patrol() error {
	root := s.config.GetRoot()
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			if s.config.IsIncludeFile(s.getFilename(path)) {
				logger.Info(fmt.Sprintf("watching path; %v", path))
				return s.watcher.Add(path)
			}
			return nil
		}

		if s.config.IsExcludeDir(path) {
			logger.Warn(fmt.Sprintf("excluding path; %v", path))
			return filepath.SkipDir
		}

		if s.config.IsTmpDir(path) {
			logger.Warn(fmt.Sprintf("excluding path; %v", path))
			return filepath.SkipDir
		}

		if s.isHiddenFile(s.getFilename(path)) {
			logger.Warn(fmt.Sprintf("excluding path; %v", path))
			return filepath.SkipDir
		}

		if s.config.IsIncludeDir(path) {
			logger.Info(fmt.Sprintf("watching path; %v", path))
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

func (s *Sentinel) restartProcess() {
	logger.Info("Restarting BARF...")
	buildCmd := exec.Command("/bin/sh", "-c", s.config.GetBuildCmd())
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	err := buildCmd.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("Error building the application: %v", err))
		if s.config.StopOnError {
			constant.ShutdownChan <- syscall.SIGINT
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(server.Augment.ShutdownTimeout)*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.HTTP.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(fmt.Sprintf("error shutting down server; %v", err))
		return
	}

	cmd := exec.Command("/bin/sh", "-c", s.config.GetBin())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		logger.Error(fmt.Sprintf("Error restarting process: %v", err))
		return
	}

	err = cmd.Wait()
	if err != nil {
		logger.Error(fmt.Sprintf("Error waiting for restarted process: %v", err))
		return
	}
}
