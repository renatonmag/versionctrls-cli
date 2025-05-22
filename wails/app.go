package main

import (
	"context"
	"os"

	"github.com/renatonmag/fsbroker"
	appconfig "github.com/renatonmag/version-ctrls/pkg/config"
	"github.com/renatonmag/version-ctrls/pkg/fs"
	"github.com/renatonmag/version-ctrls/pkg/git"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	repo *git.LocalRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) GetAppConfig() *appconfig.ApplicationConfig {
	return a.repo.AppConfig
}

func (a *App) startWatching() {
	_repo, err := git.NewRepository()
	if err != nil {
		runtime.LogErrorf(a.ctx, "error creating repository: %v", err)
		return
	}
	head := git.GetHeadPath(_repo.AppConfig.Repository.Path)
	_repo.Open()
	fileContent, err := os.ReadFile(head)
	if err != nil {
		runtime.LogErrorf(a.ctx, "error reading file %s: %v", head, err)
		return
	} else {
		// runtime.EventsEmit(a.ctx, "branch-checkout", string(fileContent))
		_repo.SetMainRepoBranch(string(fileContent))
	}
	a.repo = _repo

	config := fsbroker.DefaultFSConfig()
	config.IgnorePath = _repo.AppConfig.Repository.Path
	broker, err := fsbroker.NewFSBroker(config)
	if err != nil {
		runtime.LogErrorf(a.ctx, "error creating FS Broker: %v", err)
		return
	}
	defer broker.Stop()

	if err := broker.AddWatch(head); err != nil {
		runtime.LogErrorf(a.ctx, "error adding watch: %v", err)
		return
	}

	if err := broker.AddRecursiveWatch(_repo.Path); err != nil {
		runtime.LogErrorf(a.ctx, "error adding watch: %v", err)
		return
	}

	broker.Start()

	runtime.LogDebug(a.ctx, "WATCHING ==========================================")
	for {
		select {
		case event := <-broker.Next():
			if event.Type == fsbroker.Create {
				git.OnCreate(_repo, event)
			}

			if event.Type == fsbroker.Modify {
				git.OnModify(_repo, event)
			}

			if event.Type == fsbroker.Move {
				git.OnMove(_repo, event)
			}
			if event.Type == fsbroker.Remove {
				// git.OnRemove(_repo, event)

				if event.Path == head {
					if err := broker.AddWatch(head); err != nil {
						runtime.LogErrorf(a.ctx, "error re-adding watch: %v", err)
					}
					// Open and print file contents
					fileContent, err := os.ReadFile(head)
					if err != nil {
						runtime.LogErrorf(a.ctx, "error reading file %s: %v", head, err)
					} else {
						runtime.EventsEmit(a.ctx, "branch-checkout", string(fileContent))
						runtime.LogDebug(a.ctx, string(fileContent))
						_repo.SetMainRepoBranch(string(fileContent))
					}
				}
			}
			runtime.LogDebugf(a.ctx, "fs event has occurred: type=%s, path=%s, timestamp=%s, properties=%v", event.Type.String(), event.Path, event.Timestamp, event.Properties)
		case error := <-broker.Error():
			runtime.LogErrorf(a.ctx, "an error has occurred: %v", error)
		}
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startWatching()
	runtime.LogDebug(a.ctx, "STARTUP ==========================================")
}

func (a *App) ReadFolder(path string) *fs.Folder {
	folder := fs.CreateDirectoryStructure(path)
	return folder
}
