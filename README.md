# zippia

zippia はシンプルな日本の郵便番号検索APIサーバーです。
データベースを梱包しているため、バイナリひとつで動作します。

## Usage

### シンプルな使い方
```
$ zippia
```
127.0.0.1:5000 をバインドして HTTP サーバーが立ち上がります

### バインドする host, port を指定

```
$ zippia --host 0.0.0.0 --port 8080
```

詳細は `zippia -h` を確認してください。

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
