package routes

import (
	"github.com/agnarbriantama/DesaNgembruk-Backend/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	//user
	app.Post("/api-blog-ngebruk/register", controller.Register)
	app.Post("/api-blog-ngebruk/login", controller.Login)
	app.Patch("/api-blog-ngebruk/user/:id_user", controller.UpdateUser)
	app.Patch("/api-blog-ngebruk/users/:id_user", controller.UpdateUserWithoutPassword)
	app.Get("/api-blog-ngebruk/user", controller.GetAllUsers)
	app.Get("/api-blog-ngebruk/user-id", controller.GetUserById)
	app.Get("/api-blog-ngebruk/user/:id_user", controller.GetUserId)
	app.Delete("/api-blog-ngebruk/user/:id_user", controller.DeleteUser)
	app.Get("/api-blog-ngebruk/user/:id_user/blogger", controller.GetBloggersByUser)
	app.Post("/api-blog-ngebruk/change-password", controller.ChangePassword)

	//Blogger
	app.Post("/api-blog-ngebruk/CreateBlogger", controller.CreateBlogger)
	app.Patch("/api-blog-ngebruk/blogger/:id_blogger", controller.EditBlogger)
	app.Get("/api-blog-ngebruk/blogger/:id_blogger", controller.GetBloggerByID)
	app.Get("/api-blog-ngebruk/blogger", controller.GetAllBloggers)
	app.Get("/api-blog-ngebruk/blogger-byuser", controller.UserBloggersHandler)
	app.Delete("/api-blog-ngebruk/blogger/:id_blogger", controller.DeleteBlogger)

	//Kategori
	app.Get("/api-blog-ngebruk/kategori", controller.GetAllKategori)
	app.Post("/api-blog-ngebruk/kategori", controller.CreateKategori)
	app.Patch("/api-blog-ngebruk/kategori/:id", controller.UpdateKategori)
	app.Delete("/api-blog-ngebruk/kategori/:id", controller.DeleteKategori)
	app.Get("/api-blog-ngebruk/kategori/:id", controller.GetCategoryById)
	app.Get("/kategori/:id_kategori/blogger", controller.BloggerByCategoryHandler)

	app.Post("/api-blog-ngebruk/upload_image", controller.UploadImages)

}
