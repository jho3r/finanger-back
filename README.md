# finanger-back

This is the backend for the Finanger project. It is a REST API built with Golang

## Requirements

- Golang 1.21.1 or higher
- Docker
- Postgres

## Installation

1. Clone the repository
2. Run `go mod download` to download all the dependencies

## Debug with VSCode

1. Install the Go extension for VSCode
2. Create a `.envs` folder in the root of the project
3. Create a `local.env` file in the `.envs` folder following the `.env.example` file
4. Add the following configuration to your `launch.json` file

    ``` json
    {
        "version": "0.2.0",
        "configurations": [
            {
                "name": "Debug",
                "type": "go",
                "request": "launch",
                "mode": "debug",
                "program": "${workspaceFolder}/cmd/main.go",
                "envFile": "${workspaceFolder}/.envs/local.env",
            }
        ]
    }
    ```

5. Run the debugger
6. You can now add breakpoints and debug the application

## Migrations

1. Install the migrate CLI tool:

   ``` bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

2. Use the makefile to create a new migration:

   ``` bash
   make migrate-new name=<migration_name>
   ```

3. Add the SQL queries to the migration file
4. Use the makefile to export the environment variables: `make envs` *(this will export the environment variables from the `.envs/local.env` file)*
5. Use the makefile to run the migrations: `make migrate-up` or `make migrate-down`

## Run the application locally

1. Create a `.envs` folder in the root of the project
2. Create a `local.env` file in the `.envs` folder following the `.env.example` file
3. Use the makefile to export the environment variables: `make envs` *(this will export the environment variables from the `.envs/local.env` file)*
4. Make sure you have a Postgres database running
5. Use the makefile to run the application: `make run`
