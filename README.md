# Joi-Yo-Executor
- JOI予選でのテストケース実行を容易にするために作ったツールです

# How to install
- go get github.com/cs3238-tsuzu/joi-yo-executor

# Usage
- フォルダ名と問題番号を一致させる
- joi-yo-executor [the path of execution file(default: ./a.out)]

# Environmental Variables
- JOI\_YO\_IN\_TEMPLATE: 入力ファイル名の形式を指定。{{.Prob}}と{{.Case}}で問題番号/入力ケース番号を補完
- JOI\_YO\_OUT\_TEMPLATE: 出力ファイル名の形式を指定。上と同じ
- JOI\_YO\_EXECUTABLE\_PATH: 実行ファイルのパス(引数で指定があればそれを優先)
- JOI\_YO\_PROB: 問題番号の指定。フォルダ名と一致しない場合に使用
- 上3つは予め~/.bashrc等に書き込んでおくと良い

# License
- Under the MIT License
- Copyright (c) 2016 Tsuzu
