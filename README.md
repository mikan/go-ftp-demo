go-ftp-demo
===========

Simple Go implementation of FTP client for education use.

### Limitation

- Local file operations are unavailable
- Only supported passive mode

### Usage

```
go-ftp-demo -h <HOST> -u <USER> -p <PASS>
```

All parameters and default values:

| Arg | Default value         | Description                   |
| --- | --------------------- | ----------------------------- |
| -h  | localhost             | server hostname or IP address |
| -P  | 21                    | server TCP port number        |
| -u  | anonymous             | username                      |
| -p  | anonymous@example.com | password                      |
| -d  | false                 | print debug log               |

Some user-friendly commands are supported and automatically convert to actual FTP command.
Supported commands:

| User-friendly command | FTP command |
| --------------------- | ----------- |
| cd                    | CWD         |
| ls                    | NLST        |
| dir                   | LIST        |
| cat                   | RETR        |
| rm                    | DELE        |
| pwd                   | PWD         |

## License

[BSD 3-Clause](LICENSE)

### Contact

- [mikan](https://github.com/mikan)
