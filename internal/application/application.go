package application

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brianseitel/go-project-template/internal/middleware"
	"github.com/brianseitel/go-project-template/internal/sample"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Start() {
	// Instantiate the logger!
	logger, err := zap.NewDevelopment()

	// If something went wrong, let's bail immediately.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to instantiate logger: %v\n", err)
		os.Exit(1)
	}

	/*
		Config / Environment Variables
		------------------------------
		The tool we use to get config vars is called viper. For more
		information, check it out here: https://github.com/spf13/viper.

		The usage is pretty simple:
		dbHost := viper.GetString("DATABASE_URL")

		Viper will attempt to do two things:
		1) read from a config file, typically .env or .{project}.json
		in the current directory or in your $HOME directory

		2) if that fails, it will attempt to read $DATABASE_URL from
		your environment
	*/

	/*
		# Declaring Dependencies
		-------------------------
		In order to properly utilize dependency injection and ensure that
		all of our code is testable, the best thing to do is treat this
		appligation.go file as a bootstrapping file, instantiate all
		dependencies, and then inject them into the appropriate controllers
		or structs that we will use later. Some examples to follow.

		Below we have some examples of dependencies. Feel free to un-comment
		them and use them as-is. If you are unfamiliar with Go, it's best to
		just use them as-is, or contact one of the Go experts on the team
		if you feel customizations are needed.
	*/

	/*
		DEPENDENCY EXAMPLE: Database (Postgres)
		---------------------------------------
		If you need a Postgres database, you can use something like pgx.
		https://github.com/jackc/pgx. Here's a commented out example:

		databaseHost := viper.GetString("DATABASE_URL")

		// Instantiates a pool that connects to the database.e
		dbpool, err := pgxpool.Connect(context.Background(), databaseHost))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		// The "defer" function will ensure that the dbpool connection is always
		// closed at the end, preventing memory leaks.
		defer dbpool.Close()
	*/

	/*
		REDIS EXAMPLE
		-------------
		If you need to use Redis, you can use the go-redis client
		(https://github.com/go-redis/redis) and the following code to inject:

		rdb := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("REDIS_URL"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	*/

	// Define your controllers here. We do this here so that we can inject
	// dependencies into it from the bootstrapping part of the app. You'll
	// see this make more sense when we write tests.
	sampleController := sample.Controller{
		Logger: logger,
		// DB: dbpool,  // Uncomment if you're using the dbpool from above
		// Redis: rdb, // Uncomment if you're using the Redis client from above
	}

	// This is your basic router from gorilla/mux. (https://github.com/gorilla/mux)
	// We define the router here, then pass it into each controller's Register()
	// func so that each controller can add the routes that they need to add.
	router := mux.NewRouter()

	// The Use() method adds middleware. If you don't know what middleware is,
	// here's a good resource: https://drstearns.github.io/tutorials/gomiddleware/
	router.Use(middleware.TimingMiddleware)

	// Register the routes for our sample controller
	sampleController.Register(router)

	// Start our HTTP server and listen for the router endpoints!
	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
