# first build myprocess.go
```
    go build -o myprocess myprocess.go
```

# start 2 servers listening at 4040 and 8080
```
    # Terminal 1
    nc -l -k 4040

    # Terminal 2
    nc -l -k 8080
```

# In another terminal run the myprocess
```
    ./myprocess
```

# check the process there or not
```
    ps aux | grep myprocess
```