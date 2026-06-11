# date_windows_shim

Simple shim to mimic the Linux `date` command on Windows.

## Usage

```bash
go run .              # default date-style output
go run . -u           # UTC output
go run . +%F\ %T      # GNU date-style format subset
go build -o datew.exe .
```
