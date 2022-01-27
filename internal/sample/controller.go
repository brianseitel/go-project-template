package sample

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Controller is the basic struct that holds any dependencies that
// this set of resources may need. For this basic case, it only requires
// a Logger. However, if one were to utilize a data store of some kind,
// you would add it here. See the commented out examples for more information.
type Controller struct {
	Logger *zap.Logger

	// Examples of injected dependencies:

	// DB *pgxpool.Pool
	// Redis *redis.Client
	// HTTPClient *http.Client
}

// Register accepts a mux.Router, and then we define each of the routes and
// the handlers that should be associated with this controller. It's better
// to do this in this package as a method on the Controller, because later
// we can effectively set up the routes for testing without having to go
// all the way back to the bootstrapping function.
func (c *Controller) Register(r *mux.Router) {
	r.Handle("/v1/hello", c.Hello()).Methods(http.MethodGet)
	// r.Handle("/v1/hello/db", c.HelloDatabase()).Methods(http.MethodGet)
	// r.Handle("/v1/hello/redis", c.HelloRedis()).Methods(http.MethodGet)
}

// Hello is a sample Hello World! endpoint. It accepts no parameters, but
// returns a http.HandlerFunc. There are other ways to do this, but over
// time and experience, we've found that this is one of the ideal ways
// to set up controller methods. Everything inside that handler func can
// be replaced with business logic!
func (c *Controller) Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := SampleResponse{
			Greeting: "Hello World",
		}

		output, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Oops, something went wrong."))
			return
		}

		w.Write(output)
	}
}

// HelloDatabase is a sample endpoint that returns a greeting from the database
// func (c *Controller) HelloDatabase() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var greeting string
// 		err := c.DB.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 			os.Exit(1)
// 		}

// 		fmt.Println(greeting)
// 	}
// }

// HelloRedis is a sample endpoint that returns a greeting from Redis
// func (c *Controller) HelloRedis() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		greeting, err := c.Redis.Get(r.Context(), "greeting").Result()
// 		if err != nil {
// 			panic(err)
// 		}
// 		w.Write([]byte(greeting))
// 	}
// }
