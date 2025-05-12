package routes

import (
	"log"
	"os"

	"github.com/helshabini/fsbroker"
	"github.com/renatonmag/version-ctrls-cli/pkg/fs"
	"github.com/renatonmag/version-ctrls-cli/pkg/git"

	"github.com/gofiber/fiber/v2"
)

var _repo *git.Repository

type InitRepoRequest struct {
	Path string `json:"path"`
}

type CreateBranchRequest struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
}

func initRepo(c *fiber.Ctx) error {
	var request InitRepoRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	err := git.InitRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Repository initialized")
}

func createBranch(c *fiber.Ctx) error {
	var request CreateBranchRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.CreateBranch(request.Branch, "master")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Branch created")
}

type DeleteBranchRequest struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
}

func deleteBranch(c *fiber.Ctx) error {
	var request DeleteBranchRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.DeleteBranch(request.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Branch deleted")
}

type ReproduceWithHardlinksRequest struct {
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Ignore string `json:"ignore"`
}

// func reproduceWithHardlinks(c *fiber.Ctx) error {
// 	var request ReproduceWithHardlinksRequest
// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
// 	}

// 	rep := fs.NewFsService(request.Ignore).Replicate
// 	srcPaths, err := rep.DiffDirs(request.Src, request.Dst)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}

// 	err = rep.CreateHardlinks(srcPaths[request.Src], request.Dst)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}
// 	return c.SendString("Hardlinks reproduced")
// }

type CleanDirectoryRequest struct {
	Path string `json:"path"`
}

func cleanDirectory(c *fiber.Ctx) error {
	var request CleanDirectoryRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	err := fs.NewFsService("").Replicate.CleanWorkingTree(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Directory cleaned")
}

type CommitFileOnChangeRequest struct {
	Path   string `json:"path"`
	Head   string `json:"head"`
	Watch  string `json:"watch"`
	Ignore string `json:"ignore"`
}

func commitFileOnChange(c *fiber.Ctx) error {
	var request CommitFileOnChangeRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	var err error
	_repo, err = git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	_repo.Open()
	fileContent, err := os.ReadFile(request.Head)
	if err != nil {
		log.Printf("error reading file %s: %v", request.Head, err)
	} else {
		_repo.SetMainRepoBranch(string(fileContent))
	}
	err = _repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	config := fsbroker.DefaultFSConfig()
	config.IgnorePath = request.Ignore
	broker, err := fsbroker.NewFSBroker(config)
	if err != nil {
		log.Fatalf("error creating FS Broker: %v", err)
	}
	defer broker.Stop()

	if err := broker.AddWatch(request.Head); err != nil {
		log.Printf("error adding watch: %v", err)
	}

	if err := broker.AddRecursiveWatch(request.Watch); err != nil {
		log.Printf("error adding watch: %v", err)
	}

	broker.Start()

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

				if event.Path == request.Head {
					if err := broker.AddWatch(request.Head); err != nil {
						log.Printf("error re-adding watch: %v", err)
					}
					// Open and print file contents
					fileContent, err := os.ReadFile(request.Head)
					if err != nil {
						log.Printf("error reading file %s: %v", request.Head, err)
					} else {
						_repo.SetMainRepoBranch(string(fileContent))
					}
				}
			}
			log.Printf("fs event has occurred: type=%s, path=%s, timestamp=%s, properties=%v", event.Type.String(), event.Path, event.Timestamp, event.Properties)
		case error := <-broker.Error():
			log.Printf("an error has occurred: %v", error)
		}
	}
}

type DiffDirsRequest struct {
	Dir1   string `json:"dir1"`
	Dir2   string `json:"dir2"`
	Ignore string `json:"ignore"`
}

func diffDirs(c *fiber.Ctx) error {
	var request DiffDirsRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	result, err := fs.NewFsService(request.Ignore).Replicate.DiffDirs(request.Dir1, request.Dir2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(result)
}

type DetectCheckoutRequest struct {
	Path string `json:"path"`
	Head string `json:"head"`
}

func detectCheckout(c *fiber.Ctx) error {
	var request DetectCheckoutRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	config := fsbroker.DefaultFSConfig()
	broker, err := fsbroker.NewFSBroker(config)
	if err != nil {
		log.Fatalf("error creating FS Broker: %v", err)
	}
	defer broker.Stop()

	if err := broker.AddWatch(request.Head); err != nil {
		log.Printf("error adding watch: %v", err)
	}

	broker.Start()

	for {
		select {
		case event := <-broker.Next():
			if event.Type == fsbroker.Remove {
				if err := broker.AddWatch(request.Head); err != nil {
					log.Printf("error re-adding watch: %v", err)
				}
				// Open and print file contents
				fileContent, err := os.ReadFile(request.Head)
				if err != nil {
					log.Printf("error reading file %s: %v", request.Head, err)
				} else {
					_repo.SetMainRepoBranch(string(fileContent))
				}
			}
			log.Printf("fs event has occurred: type=%s, path=%s, timestamp=%s, properties=%v", event.Type.String(), event.Path, event.Timestamp, event.Properties)
		case error := <-broker.Error():
			log.Printf("an error has occurred: %v", error)
		}
	}
}

type RenameBranchRequest struct {
	Path      string `json:"path"`
	Branch    string `json:"branch"`
	NewBranch string `json:"newBranch"`
}

func renameBranch(c *fiber.Ctx) error {
	var request RenameBranchRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.RenameBranch(request.Branch, request.NewBranch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Branch renamed")
}

type FileMapRequest struct {
	Path   string `json:"path"`
	Ignore string `json:"ignore"`
}

func fileMap(c *fiber.Ctx) error {
	var request FileMapRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	fileMap, err := fs.NewFsService(request.Ignore).Replicate.BuildFileMap(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fileMap)
}

type GetIgnoreServiceRequest struct {
	Path string `json:"path"`
	File string `json:"file"`
}

func getIgnoreService(c *fiber.Ctx) error {
	var request GetIgnoreServiceRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	ignoreService := fs.NewIgnoreService(request.Path)
	return c.JSON(ignoreService.MatchesPath(request.File))
}

type SyncDirsRequest struct {
	Src    string `json:"src"`
	Ignore string `json:"ignore"`
}

func syncDirs(c *fiber.Ctx) error {
	var request SyncDirsRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// fsService := fs.NewFsService(request.Ignore)
	return c.SendString("Directories synced")
}

type WatchDirRequest struct {
	Path   string `json:"path"`
	Ignore string `json:"ignore"`
}

func watchDir(c *fiber.Ctx) error {
	var request WatchDirRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	config := fsbroker.DefaultFSConfig()
	config.IgnorePath = request.Ignore
	broker, err := fsbroker.NewFSBroker(config)
	if err != nil {
		log.Fatalf("error creating FS Broker: %v", err)
	}
	defer broker.Stop()

	if err := broker.AddRecursiveWatch(request.Path); err != nil {
		log.Printf("error adding watch: %v", err)
	}

	broker.Start()

	for {
		select {
		case event := <-broker.Next():
			log.Printf("fs event has occurred: type=%s, path=%s, timestamp=%s, properties=%v", event.Type.String(), event.Path, event.Timestamp, event.Properties)
		case error := <-broker.Error():
			log.Printf("an error has occurred: %v", error)
		}
	}
}

type RemoveFileRequest struct {
	Path   string `json:"path"`
	Ignore string `json:"ignore"`
}

func removeFile(c *fiber.Ctx) error {
	var request RemoveFileRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	err := os.Remove(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("File removed")
}

type AddRemoteRequest struct {
	Path   string `json:"path"`
	Remote string `json:"remote"`
	URL    string `json:"url"`
}

func addRemote(c *fiber.Ctx) error {
	var request AddRemoteRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.AddRemote(request.Remote, request.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Remote added")
}

type UpdateRemoteRequest struct {
	Path   string `json:"path"`
	Remote string `json:"remote"`
	URL    string `json:"url"`
}

func updateRemote(c *fiber.Ctx) error {
	var request UpdateRemoteRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.UpdateRemote(request.Remote, request.URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Remote updated")
}

type PushRequest struct {
	Path   string `json:"path"`
	Remote string `json:"remote"`
}

func push(c *fiber.Ctx) error {
	var request PushRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	repo, err := git.NewRepository(request.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = repo.Push(request.Remote)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Pushed")
}
