# goapp - MiTo Team's Golang projects application base
Go projects application base

## Add as git submodule
```
git submodule add -b main https://github.com/mitoteam/goapp.git internal/goapp
```

Disable commit hash tracking for submodule in `.gitmodules` :
```
[submodule "goapp"]
    ignore = all # ignore hash changes
```

Add `use ./internal/goapp` to main project `go.work`

Do `go mod tidy`

## Useful commands

Pull main project with submodules:

```
git pull --recurse-submodules
```
