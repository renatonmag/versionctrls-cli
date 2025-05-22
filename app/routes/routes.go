package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/init-repo", initRepo)
	app.Post("/create-branch", createBranch)
	app.Post("/delete-branch", deleteBranch)
	app.Post("/create-directory-structure", createDirectoryStructure)
	app.Post("/diff-dirs", diffDirs)
	app.Post("/get-ignore-service", getIgnoreService)
	app.Post("/file-map", fileMap)
	// app.Post("/reproduce-with-hardlinks", reproduceWithHardlinks)
	// app.Post("/sync-dirs", syncDirs)
	app.Post("/watch-dir", watchDir)
	app.Post("/detect-checkout", detectCheckout)
	app.Post("/commit-file-on-change", commitFileOnChange)
	app.Post("/remove-file", removeFile)
	app.Post("/rename-branch", renameBranch)
	app.Post("/add-remote", addRemote)
	app.Post("/update-remote", updateRemote)
	app.Post("/push", push)
	app.Post("/get-config", getConfig)
}
