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

go >= 1.18

## Installation

```
go install github.com/tetsu040e/zippia@latest
```

郵便番号のデータは概ね1ヶ月ごとに更新されます。
zippia は GitHub Actions を使って定期的に日本郵政のホームページをチェックし、データを更新しています。
`@latest` を使って再インストールすることで、最新のデータを内包したバイナリに更新できます。　　
インストール済みのバイナリのデータがいつ更新されたかを確認するには、 `zippia -vv` を実行してください。

## License

MIT

## Author

[tetsu040e](https://github.com/tetsu040e)
