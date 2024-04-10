# goappbase
Go projects application base

## Add as git submodule
```
git submodule add -b main https://github.com/mitoteam/goappbase.git
```

Disable commit hash tracking for submodule in `.gitmodules` :
```
[submodule "goappbase"]
    ignore = all # ignore hash changes
```

Add `use ./goappbase` to main project `go.work`

Do `go mod tidy`

## Useful commands

Pull main project with submodules:

```
git pull --recurse-submodules
```
