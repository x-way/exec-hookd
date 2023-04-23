# exec-hookd - execute commands from webhooks

exec-hookd listens for HTTP POST requests and executes pre-defined commands when a request for a matching path is received.

## Usage
Run the go binary from your local path
```
# exec-hookd -f exec-hookd.json > /dev/null 2>&1 &
```

## Configuration

exec-hookd reads its configuration from the config file `exec-hookd.json` (default location, can be changed with the `-f` flag).

Sample config:
```
{
  "Port": 8059,
  "HookList": [
    {
      "Path": "/myhook",
      "Exec": [
        {
          "Cmd": "/usr/bin/somecmd",
          "Args": [
            "--some",
            "parameter"
          ],
          "Timeout": "5s"
        }
      ]
    }
  ]
}
```
