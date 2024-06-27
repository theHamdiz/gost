# GoSt - Go for the lazy

Gost is a go web projects bootstrapper & automation tool!

## Features

- Work with multiple backends.
- Customize the frontend however you like.
- Use any db system you prefer.
- Resource management concept, similar to rails but through commands.
- Everything is a plugin, you can develop your own.
- White label code, less dependencies.
- NaturalOrm Plugin.
- Detailed yet easy dir structure.
- Easy to use, easy to extend.

## Installation

To install and run this project, follow these steps:

1. Clone the repository:

```sh
git clone https://github.com/theHamdiz/gost.git
cd gost
```

2. Install dependencies:

```sh
go mod tidy
```

3. Set up environment variables (if any):

```sh
cp .env.dev .env
# Edit the .env file with your configuration
```

4. Run the application:

```sh
go run cmd/app/main.go
```

## Usage

### Creating a Project

To create a project, use:

```sh
gost create appName
```

or

```sh
gost c appName
```

> Anywhere in gost, if you can use create as a command you can also use init & new to do the same thing, they're all synonyms for each other.

### Running the Project

To start the project, use:

```sh
gost run
```

or

```sh
gost r
```

### Project Structure

By default gost creates the following structure for you:

```
"```
			
    ├── app
│   ├── api
│   │   └── v1
│   ├── assets
│   │   └── static
│   │       ├── css
│   │       ├── img
│   │       └── js
│   ├── cfg
│   │   └── cfg.go
│   ├── db
│   │   ├── migrations
│   │   │   ├── create_db_1719424521947950600.sql
│   │   │   └── create_db_1719424522725851300.sql
│   │   ├── data.db
│   │   └── db.go
│   ├── events
│   │   └── events.go
│   ├── handlers
│   │   ├── api
│   │   │   └── api.go
│   │   ├── backend
│   │   │   ├── about.go
│   │   │   ├── auth.go
│   │   │   ├── landing.go
│   │   │   └── views.go
│   │   └── frontend
│   │       ├── about.go
│   │       ├── auth.go
│   │       ├── landing.go
│   │       └── views.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── notifier.go
│   │   ├── rateLimiter.go
│   │   ├── recoverer.go
│   │   └── requestId.go
│   ├── router
│   │   └── router.go
│   ├── services
│   │   ├── db.go
│   │   ├── logger.go
│   │   └── rateLimiter.go
│   ├── types
│   │   ├── core
│   │   │   └── gost.go
│   │   └── models
│   └── ui
│       ├── backend
│       │   ├── assets
│       │   │   ├── css
│       │   │   └── js
│       │   ├── components
│       │   │   └── index.js
│       │   ├── pages
│       │   │   └── index.js
│       │   ├── store
│       │   │   └── index.js
│       │   ├── README.md
│       │   ├── index.html
│       │   ├── package.json
│       │   ├── robots.txt
│       │   ├── signin.html
│       │   ├── signup.html
│       │   └── vite.config.js
│       ├── components
│       │   ├── footer
│       │   │   └── footer.templ
│       │   ├── header
│       │   │   └── header.templ
│       │   ├── navigation
│       │   │   └── sidebar.templ
│       │   └── head.templ
│       ├── errors
│       │   ├── 404.templ
│       │   └── 500.templ
│       ├── frontend
│       │   ├── assets
│       │   │   ├── css
│       │   │   └── js
│       │   ├── components
│       │   │   └── index.js
│       │   ├── pages
│       │   │   ├── index.js
│       │   │   ├── signin.templ
│       │   │   └── signup.templ
│       │   ├── store
│       │   │   └── index.js
│       │   ├── README.md
│       │   ├── index.html
│       │   ├── package.json
│       │   ├── robots.txt
│       │   └── vite.config.js
│       ├── layouts
│       │   ├── app.templ
│       │   └── base.templ
│       ├── pages
│       │   ├── about.templ
│       │   └── home.templ
│       ├── public
│       │   └── index.html
│       ├── shared
│       ├── README.md
│       ├── embed.go
│       └── views.go
├── cmd
│   ├── app
│   │   └── main.go
│   └── scripts
├── log
├── plugins
│   ├── auth
│   ├── core
│   │   ├── config.go
│   │   └── core.go
│   └── db
│       ├── dialects
│       │   ├── db2.go
│       │   ├── dialects.go
│       │   ├── firebird.go
│       │   ├── mariadb.go
│       │   ├── mysql.go
│       │   ├── oracle.go
│       │   ├── postgresql.go
│       │   ├── sqlite.go
│       │   └── sqlserver.go
│       └── db.go
├── public
│   └── assets
├── storage
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── package-lock.json
└── package.json
` 
```

### Running Tests

To run tests, use:

```sh
gost test
```

## Configuration

List any configuration settings for your project:

- `DATABASE_URL`: The URL of your database.
- `PORT`: The port on which the server will run.

## Contributing

We welcome contributions! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch with your feature or bug fix.
3. Commit your changes.
4. Push the branch to your fork.
5. Create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgements

Thanks to the contributors and the open-source community for their valuable input and support.
