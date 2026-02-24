# passgen

安全な乱数を使用してパスワードを生成する
シンプルな Go 製 CLI ツールです。

`crypto/rand` を使用しているため、
暗号学的に安全なランダム値で生成されます。

---

## 特徴

* `crypto/rand` による安全な乱数生成
* 長さ指定可能
* 複数生成可能
* 文字セットの追加・削除対応
* 標準入力から文字セットを受け取り可能
* `--long-option` 形式にも対応

---

## インストール

### ビルド

```bash
make build
```

または

```bash
go build -o bin/passgen ./cmd/passgen
```

---

## 使い方

### デフォルト（8文字・1個）

```bash
./bin/passgen
```

例:

```
MdVRWYbc
```

---

## オプション

| オプション            | 説明                  | デフォルト |
| ---------------- | ------------------- | ----- |
| `-l`, `--length` | パスワード長              | 8     |
| `-n`, `--number` | 生成数                 | 1     |
| `-a`, `--add`    | 指定文字をデフォルト文字セットに追加  | false |
| `-d`, `--delete` | 指定文字をデフォルト文字セットから削除 | false |

---

## 使用例

### 長さを指定

```bash
./bin/passgen -l 10
```

---

### 10個生成

```bash
./bin/passgen -n 10
```

---

### 数字だけで生成（標準入力）

```bash
echo "1234567890" | ./bin/passgen -n 5
```

出力例:

```
81283838
84370005
56763370
```

---

### デフォルト文字 + 数字を追加

```bash
echo "1234567890" | ./bin/passgen -a -n 5
```

---

### 記号を追加

```bash
echo ",.;+:*" | ./bin/passgen -a -n 5
```

---

### 文字を削除

```bash
echo "abcABC" | ./bin/passgen -d
```

---

### 文字セットを完全指定

```bash
./bin/passgen 0123456789
```

---

## 注意

`-a` と `-d` は同時に指定できません。

```
error: -a/--add と -d/--delete は同時に指定できません
```

---

## セキュリティについて

本ツールは `crypto/rand` を使用しています。

`math/rand` は使用していません。

そのため、予測困難なパスワードを生成できます。

---

## 仕様

デフォルト文字セット:

```
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
```

* 重複文字は自動的に除去されます
* 標準入力と引数の両方に対応しています

---

## ディレクトリ構成

```
passgen/
 ├── cmd/passgen/main.go
 ├── Makefile
 └── README.md
```
