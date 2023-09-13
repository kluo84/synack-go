# mBot (A Go-Based Synack Mission Bot)

This currently has bare minimal functionality (onboarding and claiming missions). Future plans for the bot include notifications, auto-relogin, and automated recon.

## Build Instructions
You can build mBot straight from the source directory:
```
go build .
```

## Usage
Example:
```
mBot -t asdfasdfasdfasdfasdfsadfasdfasdfasdf -d 30
```
For help, use the `-h` flag:
```
mBot -h
```

| Flag | Description | Example |
|------|-------------|---------|
| `-d` | Specifies a delay between polls (to please Synack) | `mBot -d 30` |
| `-t` | Passes your Authorization: Bearer token to the bot (this is necessary) | `mBot -t asdfasdfasdfasdfasdfasdfaf` |