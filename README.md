Golang utility program for connecting via ssh and streaming log file.

Only password authentication can be used for now.

## Usage example
Accoring to config.example.yaml

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