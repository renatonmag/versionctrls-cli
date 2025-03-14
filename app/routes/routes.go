package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/init-repo", initRepo)
	app.Post("/create-branch", createBranch)
	app.Post("/delete-branch", deleteBranch)
	app.Post("/reproduce-with-hardlinks", reproduceWithHardlinks)
	// app.Post("/clean-directory", cleanDirectory)
	app.Post("/diff-dirs", diffDirs)
	app.Post("/get-ignore-service", getIgnoreService)
	app.Post("/file-map", fileMap)
	app.Post("/sync-dirs", syncDirs)
	app.Post("/watch-dir", watchDir)
}
