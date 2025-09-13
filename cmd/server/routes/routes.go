package routes 



func RegisterRoutes(router *fiber.App) {
    grImport := router.Group("/import")

	grImport.Post("/upload", handlers.ImportHandler)
}