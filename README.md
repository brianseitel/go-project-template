# Go Project Template

This is a template that is intended to be a starter Go template for new services. If you start with this, you can jump straight into business logic with a working webserver, data store dependencies, standard middleware, and config variable instantiation. There are examples out the wazoo, and the code is very well documented. If you have any questions, feel free to reach out in the #eng-golang channel.

## Getting Started

### Installing Dependencies

There's only one external dependency that you'll need for this project, and that's [cobra](https://github.com/spf13/cobra), the CLI tool we use to generate actions. There's more about this later, but for now, simply run: 

```
go get -u github.com/spf13/cobra
```

### Modules

The first thing you'll need to do is ensure that your modules are installed. Run the following command:

`go mod tidy`

This will install any module dependencies into your system. It should only take a few seconds.

If you get an error, you may need to run `go mod init` in order to let Go know that this project's dependencies are managed by `go mod`. This _should_ already be taken care of for you, but if it's not, go ahead and run that command.

### Starting the Webserver

This repo uses [cobra](https://github.com/spf13/cobra) as the CLI tool that generates actions. The most basic one has already been generated for you. You can view the available commands by running:

`go run main.go`

This should return something that looks like this:

```
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  go-project-template [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  serve       Start our HTTP server

Flags:
  -h, --help     help for go-project-template
  -t, --toggle   Help message for toggle

Use "go-project-template [command] --help" for more information about a command.
```

As you can see, there's already an available command called `serve` that will start the HTTP server. Go ahead and run that:

```
go run main.go serve
```

You should see output like the following:

```
$ go run main.go serve
Listening on port 8080...
```

You can now visit [http://localhost:8080/v1/hello](http://localhost:8080/v1/hello) to view the sample endpoint. Tada! You've got a web server!

## Project Structure

The project is laid out as follows:

```
cmd/                        -- directory contains all of the executable commands
    root.go                 -- the root command, don't modify
    serve.go                -- the starting command to serve up a HTTP server
internal/                   -- the internal directory, see below about package structure
    application/            -- the application package
        application.go      -- bootstraps the application
    middleware           
        timing.go           -- timing middleware
    sample/
        controller.go       -- a sample controller for your sample package
        response.go         -- a sample response struct
go.mod 
go.sum
LICENSE
main.go
README.md
```

### Packages in Go

In Go, code is organized in "packages" (directories where the files contain `package $name` as the first line). These packages are imported similarly to ther imports in other languages. In many languages, code is organized around structures like Model-View-Controller. It's fine if you want to do it that way, but in our experience, we've found that a superior way to organize Go code is through domain-level packages.

In other words, instead of a `controllers` package that contains all of the controllers, a `models` package that contains all of the models, and `views` package that contains all of the views, it is better to have a `users` package that contains the models, views, and controllers for everything relating to Users. There are many reasons for preferring this method, but the most immediate benefit is that you avoid cyclical dependencies.

In short, trust us: this is the right way to do it for now. Once you're an advanced Go user, you can play around with new strategies, but for now, stick with this.

#### The Internal Directory

Go is designed, in part, to allow any package to be imported by any other project. If you look at our code, we are importing things like `github.com/jackc/pgxpool` or `go.uber.org/zap`. These are packages in _other_ projects that we do not own. Similarly, Uber could import our `sample` app, if they wanted to!

Except they can't, because we've placed our packages in the `internal/` directory. By doing this, we've explicitly told Go that this is _ours_ and is not a public package. If you come from an OOP background, you can think of the `internal/` directory as a private set of packages that only this repository can access, while packages _not_ in the `internal/` directory are public packages.

Typically, if you're creating packages for external use or that are intended to be shared, you might add another `pkg/` directory at the top-level.

For our purposes, 98% of our code should be private to our services, so the recommendation is to organize within our `internal/` directory. Again, trust us for now.

## Adding New Commands

Let's say you want to add a `migrate` command to execute migrations. You can simply run

```
cobra add migrate
```

This will create a new file in your `cmd/` directory called `migrate.go`. To ensure that the command works, run:

```
go run main.go migrate
```

You should see something along the lines of `migrate called` appear in your terminal. Done!
