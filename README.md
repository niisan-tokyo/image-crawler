# What is this?

画像を抽出するWebアプリです。
シングルバイナリでローカルで動かせるので、windowsでも使えます。

# 環境構築

VS Codeを利用することを前提に作られています。
VS Code Remote Container で開くことで、ターミナル上でgolangの実行やビルドを実施できます。

リポジトリをクローンしたら、ターミナル上で

```
go mod download
```

で必要なモジュールを落としてください

# ビルド

ビルドは以下のよう実施できます。

```
statik -src=templates -f
GOOS=windows go build -v
```
作られたexeファイルはwindows上で実行可能です。