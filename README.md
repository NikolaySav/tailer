Golang utility program for connecting via ssh and streaming log file.

Only password authentication can be used for now.

## Usage example
According to config.example.yaml

Create config.yaml with server and project configurations:
```yaml
servers:
  -
    name: "develop"
    host: "111.222.333.444"
    port: "22"
    username: "username"
    password: "password"
```
Create project config for your target server:
```yaml
projects: 
  - 
    name: "projectName"
    filePath: "/var/www/projectName/current/storage/logs/laravel.log"
    server: "develop"
```

Build 
```
go build tailer.go
```

Read default number of rows
```
./tailer projectName
```

Read N number of rows
```
./tailer projectName 100
```
## Bastion host example
As usual add bastion and target server configs to servers in config.yaml:
```yaml
servers:
  -
    name: "bastionServer"
    host: "111.222.333.555"
    port: "22"
    username: "bastionUsername"
    password: "bastionPassword"
  -
    name: "target"
    host: "10.1.1.1"
    port: "22"
    username: "username"
    password: "password"
```
Create project config for your target server:
```yaml
projects: 
  -
    name: "bastionProject"
    filePath: "/var/www/bastionProject/current/storage/logs/laravel.log"
    server: "target"
    bastionServer: "bastionServer"
```