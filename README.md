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
gost create
```

or

```sh
gost c
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
.
├── cmd             # Main applications of the project
├── app             # Private application and library code
├── pkg             # Public library code
├── web             # Web server-related files
│   ├── static      # Static files
│   └── templates   # HTML templates
├── go.mod          # Go module file
├── main.go         # Main entry point of the application
└── README.md       # This file
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
