# WebMonitor
Make a request to a URL every X seconds and notify via Email and SMS if it takes more than Y seconds to load.

Tool created with the help of Github Copilot.

# Build
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o builds/webmonitor -trimpath
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o builds/webmonitor_386 -trimpath
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o builds/webmonitor_osx -trimpath
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o builds/webmonitor_win -trimpath
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o builds/webmonitor_win_386 -trimpath
```

# Config Example
```
{
    "services": [
        {
            "enable": true,
            "name": "test",
            "url": "http://webmonitor.example.com",
            "timeout": 3,
            "interval": 60,
            "emails": ["my@email.com"],
            "smsNumbers": ["13137437189"]
        }
    ],
    "coreApp": {
        "host": "https://opalo.com.co:9905",
        "user": "user",
        "password": "pass",
        "client_id" : "client_id",
        "client_secret" : "client_secret",
        "company": "1"
    }
}
```


