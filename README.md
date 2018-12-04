# ypipe

ypipe is similar to the tee command (POSIX IEEE Std 1003.2). It copy stdin to stdout and also to a file.

## Example

```console
echo test| ypipe -o file.log
```