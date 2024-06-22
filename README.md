# GoSt - Your Go Starter Tool

`GoSt` is a highly flexible and user-friendly command-line tool designed to streamline the creation and management of `Go` projects.By providing both boilerplate code & routine task automation `GoSt` provides a seamless & efficient experience.

## Features

- **Flexible Project Creation**: Easily create new `Go` web projects with various configurations using both unnamed and named parameters or conversational style configuration builder.
- **Configurable IDE Support**: Supports a wide range of popular IDEs and editors for `Go` development.
- **Database Management**: Provides commands to create, migrate, seed, and fake database data.
- **Run Projects Easily**: Start your `Go` projects effortlessly with built-in commands.

## Installation

To install `GoSt`, ensure you have `Go` installed on your system, then run:

```sh
go get -u github.com/theHamdiz/gost
```

## Usage

### Project Creation

You can create a new `Go` project using various parameter styles:

- **Unnamed Parameters**:

  ```sh
  gost create myApp tailwindcss none gin
  ```

- **Shorthand Named Parameters**:

  ```sh
  gost create -n myApp -ui tailwindcss -c none -b echo
  ```

- **Longhand Named Parameters**:

  ```sh
  gost create --name myApp --uiFramework tailwindcss --componentFramework none --backendFramework echo
  ```

### Running Commands

- **Run Project**:

  ```sh
  gost run
  # or
  gost r
  ```

- **Database Migrations**:

  ```sh
  gost db migrate
  ```

- **Seed Database**:

  ```sh
  gost db seed
  ```

- **Fake Database Data**:

  ```sh
  gost db fake
  ```

## Configuration

`GoSt` will prompt you for any missing configuration settings the first time you run it. You can also set your preferences globally. Configuration is saved in the home directory as `.gost`, `.gost.json`, or `.gost.toml`.

### Supported IDEs and Editors

`GoSt` supports the following IDEs and editors:

- VSCode
- Goland
- IDEA
- Cursor
- Zed
- Sublime
- Vim
- Nvim
- Nano
- Notepad++
- Zeus
- LiteIDE
- Emacs
- Eclipse

## Contributing

I welcome contributions from the community! If you have suggestions, bug reports, or pull requests, please feel free to submit them, *you would be surprised to see how flexible I am to new requests.*

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes.
4. Push your branch and create a pull request.

## License

`GoSt` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgements

I would like to thank the `Go` community and contributors for their valuable feedback and support in making `GoSt` a powerful tool for `Go` developers.

---

Enjoy using `GoSt` and happy coding!
