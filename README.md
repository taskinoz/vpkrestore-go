# VPK Restore

This program checks the integrity of `.vpk` files against a remote set of hash values. If mismatches are found, it offers the option to download the correct versions of the files.

## Usage

To run the program from source:
```bash
go run your_program.go [flags]
```

To run the program from a binary:
```bash
vpkrestore.exe [flags]
```

## Flags

| Flag | Description | Default |
| ---- | ----------- | ------- |
| `-d` | Enable debug mode | `false` |
| `-a` | Enable auto mode | `false` |

## Auto Mode

Auto mode will automatically download the correct files without prompting the user. This is useful for running the program in a script.

## Debug Mode

Debug mode will print out the name of each `.vpk` file as well as the md5, sha1, and sha256 as it is being checked.

## Example

```bash
go run vpkrestore.go -d -a
```

## Building

To build the program from source:
```bash
go build -o vpkrestore.exe vpkrestore.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
