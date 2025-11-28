# SnapAttend
Welcome to the SnapAttend Project. This guide will help you to set up, build and run the SnapAttend application.

## 🔧 Prerequisites

Before you start, make sure you have the appropriate tools installed on your system:

| 🛠 Tools        | 🔢 Version | 🔗 Link                                                      |
|-----------------|------------|--------------------------------------------------------------|
| ☕ **GO**        | `1.25.4`   | [GO Download](https://go.dev/dl/)                            |
| **GoLand**      | `Latest`   | [GoLand Download](https://www.jetbrains.com/go/download/)    |
| **PostGre SQL** | `18`       | [PostGre SQL Download](https://www.postgresql.org/download/) |
| **PG Admin 4**  | `9.8`      | [PG Admin 4 Download](https://www.pgadmin.org/download/)     |

🔍 **You can verify your GO Version using the below command:**
```sh
go version
```
## GO Dependencies

Please find below the list of important maven dependencies used in this project.

| 🏷 Dependency | 🔢 Version |
|---------------|------------|
| **Viper**     | `1.21.0`   |
| **Gorm**      | `1.31.0`   |
| **Postgres**  | `1.6.0`    |

## Environment Variables

Make sure to set the following environment variables in your system before running the application.:

| Key                 | Use of the Key                                                             | Example Value                 |
|---------------------|----------------------------------------------------------------------------|-------------------------------|
| **GOBIN**           | The path where you want your GO external development tool to be installed. | %USERPROFILE%\go\bin          |
| **GOPATH**          | The path where you want your GO dependencies to be installed.              | %USERPROFILE%\go              |
| **GOROOT**          | The path where you want your GO software is installed.                     | %USERPROFILE%\Applications\GO |
| **SNAP_ATTEND_ENV** | The configuration file be used for the project.                            | dev                           |

## 📦️Build
If you make any code changes, follow the below commands to rebuild the project and generate a new package.
```sh
go clean
go build
```

## Running the Application

- Update the DB details in the config file. 
- After updating the config file with DB details, make sure the PostGre DB running.
- Start the GO Application using the below command:
```sh
go run .
```
If you are running the application in local then use air. This helps the application to restart automatically on code changes.
```sh
air
```

## 📝 Notes
- The tables will be created automatically in the database when you run the application for the first time.
- Make sure to run the PostGre SQL server script provided in the `DB_Scripts` folder to create the necessary records after starting the application.
