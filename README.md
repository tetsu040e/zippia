# zippia

zippia はシンプルな郵便番号検索 API サーバーです。  
データベースを内包しているため、バイナリひとつで動作します。

## Usage

### most simple usage
```
$ zippia
```
オプションなしで起動すると、 127.0.0.1:5000 をバインドして HTTP サーバーが立ち上がります。

### specify binds host and port

```
$ zippia --host 0.0.0.0 --port 8080
```

詳細は `zippia -h` を確認してください。

## API specification

https://tetsu040e.github.io/zippia/ を参照してください


## Requirements

go >= 1.16

## Installation

```
go install github.com/tetsu040e/zippia@latest
```

## License

MIT

## Author

[tetsu040e](https://github.com/tetsu040e)
