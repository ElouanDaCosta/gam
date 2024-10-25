# Golang Application Manager
A CLI in Golang to interact and manage a single or multiple applications in Golang.

## Prerequisites

- Golang 1.16 or higher

To see the version of Golang installed on your machine, run the following command.

```bash
go version
```

If you don't have Golang installed on your machine, you can download it from the [official website](https://golang.org/dl/).

Or using apt or brew.

```bash
sudo apt search golang-go
sudo apt install golang-go
```

```bash
brew install go
```

## Installation

Clone the repository and run the following command at the root of the CLI.

You can execute the installation script using the following command.
```bash
./install.sh
```

It will execute all the process to install the CLI.
If you do not have Golang installed on your machine, refer to the prerequisites section.

If you do not want to use the installation script, you can follow the steps below to install the CLI manually.

```bash
echo export GAC=$(pwd) >> ~/.zshrc
echo export PATH="$PATH:$(go env GOPATH)/bin"
```

```bash
mkdir storage
```

You can build the CLI using the following command.
Then you can run the CLI using the binary file.

```bash
go build -o ./bin/go-app-cli
```

If you want to use the CLI from anywhere, you can add the binary to the PATH. To do this, run the following command.
The go install command will install the binary in the $GOPATH/bin directory.
The binary name will be the same as the directory name.
Now you can run the CLI from anywhere.

```bash
go install
```

## Usage

### Create a new application

Create a new application using the following command.

```bash
./bin/go-app-cli init --name <YOUR_APP_NAME>
```

You can create a dockerfile for the application using the following command.

```bash
./bin/go-app-cli init --name <YOUR_APP_NAME> --dockerfile
```

The dockerfile will be created in the root of the application and the Golang version will be the same as the application.

### List applications

List all the applications using the following command.

```bash
./bin/go-app-cli list
```

List a specific application using the following command.

```bash
./bin/go-app-cli list --name <YOUR_APP_NAME>
```

### Remove applications

Delete a specific application from the storage using the following command.

```bash
./bin/go-app-cli prune --name <YOUR_APP_NAME>
```

Delete completely the applications using the following command.

```bash
./bin/go-app-cli prune --remove-app --name <YOUR_APP_NAME>
```

### Update applications

Update the go version of a specific application using the following command.

```bash
./bin/go-app-cli update --name <YOUR_APP_NAME> --version <YOUR_GO_VERSION>
```

Update the go version of all the applications using the following command.

```bash
./bin/go-app-cli update --all --version <YOUR_GO_VERSION>
```

## Project Structure

```bash
.
├── bin
├── cmd
├── configs
├── storage
├── templates
├── utils
├── .gitignore
├── go.mod
├── go.sum
├── install.sh
├── main.go
└── README.md

```

1. **bin**: Contains the binary file of the CLI.
2. **cmd**: Contains the commands for the CLI.
3. **configs**: Contains the configuration files for application structure that the CLI generate.
4. **storage**: Contains the storage files for the CLI.
5. **templates**: Contains the templates for main.go and dockerfiles that the CLI generate.
6. **utils**: Contains the utility functions for the CLI.
