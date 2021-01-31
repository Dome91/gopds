package main

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"gopds/configuration"
	"gopds/database"
	"gopds/services"
	"gopds/web"
)

func main() {
	// Database
	db := database.New(configuration.GetDatabasePath())
	defer db.Close()

	database.Migrate(db.DB.DB, configuration.GetMigrationsPath())
	userRepository := database.NewUserRepository(db)
	sourceRepository := database.NewSourceRepository(db)
	catalogRepository := database.NewCatalogRepository(db)

	// Services
	createUser := services.CreateUserProvider(userRepository)
	fetchAllUsers := services.FetchAllUsersProvider(userRepository)
	userExistsByRole := services.UserExistsByRoleProvider(userRepository)
	deleteUser := services.DeleteUserProvider(userRepository)
	fetchUserByUsername := services.FetchUserByUsernameProvider(userRepository)
	checkCredentials := services.CheckCredentialsProvider(userRepository)
	createSource := services.CreateSourceProvider(sourceRepository)
	fetchAllSources := services.FetchAllSourcesProvider(sourceRepository)
	deleteSource := services.DeleteSourceProvider(sourceRepository)
	fetchCatalogEntryByID := services.FetchCatalogEntryByIDProvider(catalogRepository)
	fetchCatalogEntriesByParentCatalogEntryIDInPage := services.FetchCatalogEntriesByParentCatalogEntryIDInPageProvider(catalogRepository)
	fetchCatalogRootDirectories := services.FetchCatalogRootDirectoriesProvider(catalogRepository)
	fetchAllBooksInPage := services.FetchAllBooksInPageProvider(catalogRepository)
	countBooks := services.CountBooksProvider(catalogRepository)
	countCatalogEntriesByParentCatalogEntryID := services.CountCatalogEntriesByParentCatalogEntryIDProvider(catalogRepository)
	synchronizeCatalog := services.SynchronizeCatalogProvider(sourceRepository, catalogRepository)
	generateOPDSRootFeed := services.GenerateOPDSRootFeedProvider()
	generateOPDSAllFeed := services.GenerateOPDSAllFeedProvider(catalogRepository)
	generateOPDSDirectoriesFeed := services.GenerateOPDSDirectoriesFeedProvider(catalogRepository)
	generateOPDSFeedByID := services.GenerateOPDSFeedByIDProvider(catalogRepository)

	// Initialization
	adminInitializer := NewAdminInitializer(createUser, userExistsByRole)
	Initialize(adminInitializer)

	// Web
	store := session.New()
	loginHandler := web.NewLoginHandler(store, checkCredentials, fetchUserByUsername)
	userHandler := web.NewUserHandler(createUser, fetchAllUsers, deleteUser)
	sourceHandler := web.NewSourceHandler(createSource, fetchAllSources, deleteSource, synchronizeCatalog)
	opdsHandler := web.NewOPDSHandler(generateOPDSRootFeed, generateOPDSAllFeed, generateOPDSDirectoriesFeed, generateOPDSFeedByID, fetchCatalogEntryByID)
	catalogHandler := web.NewCatalogHandler(fetchAllBooksInPage, countBooks, fetchCatalogRootDirectories, fetchCatalogEntriesByParentCatalogEntryIDInPage, countCatalogEntriesByParentCatalogEntryID, fetchCatalogEntryByID)

	withRoles := web.WithRolesProvider(store)
	basicAuth := web.BasicAuthProvider(checkCredentials)
	authorization := web.NewAuthorization(withRoles, basicAuth)
	server := web.NewServer(authorization, userHandler, loginHandler, sourceHandler, opdsHandler, catalogHandler)
	server.Start()
}
