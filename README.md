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

Pull main project with submodules:

```
git pull --recurse-submodules
```
