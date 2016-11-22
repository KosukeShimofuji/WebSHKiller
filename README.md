# webshdd - WebShell Dynamic Detector

## 概要

## webshdd-cli

WebShellが存在する可能性があるホストで実行されるクライアントプログラムであり、以下の手順に基づき処理を行う。

 * Indexファイルの作成
 * Indexファイルを考慮しながらscan対象を作成

### database schema on client side

 * index table

Indexテーブルはwebshdd-cliが対象とするスキャン対象ファイル群のstat情報を保持するものである。
Indexテーブル前回のスキャンから変更のないスキャン対象ファイル群を減らすために使用される。

| 項目 | 意味 | 
|------|------|
|File  | filepath |
|Size  | filesize | 
|Perm  | permission | 
|Uid   | user id |
|Gid   | group id |
|Atime | access time |
|Ctime | createtion time|
|Mtime | modification time |
|sha256 | file hash |
|_update| last update |

 * report table

reportテーブルは実行したscanに関するレポートを補完するために用いられる。

| 項目 | 意味 | 
|------|------|
|File  | filepath |
|date | scan開始日時 |
|processing time|scanに要した時間|
|result| スキャン結果 |

# 参考文献

 * [WebSHArk 1.0: A Benchmark Collection for Malicious Web Shell Detection](https://pdfs.semanticscholar.org/d2de/06d1e4e07890c9b27bdb4baa07c1922b3c16.pdf)
 * https://github.com/b374k/b374k
 * https://github.com/openstack/golang-client

