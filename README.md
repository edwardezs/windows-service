### Windows Service
Runs child process of a Windows service.

Logs service events to a rotating file `service.log`:
```json5
Process started
{"level":"info","time":"2024-05-26T09:58:09+03:00","message":"Starting server"}
{"level":"info","time":"2024-05-26T09:58:26+03:00","message":"Shutting down server"}
{"level":"info","time":"2024-05-26T09:58:26+03:00","message":"Server stopped"}
Process stopped
```

If the child process crashes, it attempts to restart it reporting the reason of crash:
```json5
Process started
{"level":"info","time":"2024-05-26T13:35:03+03:00","message":"Starting server"}
Process exited with error: exit status 1, attempting restart
Process restarted
{"level":"info","time":"2024-05-26T13:35:29+03:00","message":"Starting server"}
```

Operations (only in Administrator mode):
- `make build` - builds the Windows service and test child process binaries
- `make install` - installs the Windows service (without registry entry)
- `make start` - starts the Windows service process in the background
- `make stop` - stops the Windows service process
- `make delete` - deletes the Windows service. If the service is running, it will be stopped first

You can also manage the service through Task Manager.

### Configuration File

```json5
{
  // name of the registered Windows service (required)
  "name": "service",
  // description of the service (required)
  "description": "Windows service",
  // absolute path to the parent process binary (required)
  "parentExecPath": "C:/Users/user/service.exe",
  // absolute path to the child process binary (required)
  "childExecPath": "C:/Users/user/server.exe",
  // arguments for launching the child process (optional)
  "childExecArgs": ["-config", "C:/Users/user/config.json"],
  // path to the log file for the child process
  // if only a file name is provided, the file will be created in the child process binary's directory
  "logFilePath": "service.log",
  // maximum size of the log file before rotating (optional)
  // when the size is exceeded, a new file will be created with the specified name, and the old log file will be renamed
  "logFileMaxSizeMB": 50,
  // maximum number of log file rotations (optional)
  "logFileMaxBackups": 3,
  // maximum retention period for the log file in days (optional)
  "logFileMaxAgeDays": 28,
  // compress the log file using gzip (optional)
  "logFileCompress": false
}
```

If the child process of the service has its own configuration file, the paths in it must also be absolute!

Example in service.config.json.example.

### Tests

To run Windows service tests:

- Run `make test` as an Administrator