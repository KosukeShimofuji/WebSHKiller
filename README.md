# WebSHKiller - Break through obfusucation and kill webshell

WebSHkillerはWebShellを検知するためのソフトウェアです。
WebSHkillerの特徴は動的解析を行うことによって難読化を施したwebshellを検知することができる点です。
動的解析はOpenStack上に作成されたdocker host(Jail)のdockerコンテナ上で実行され、解析が終わればdocker host(jail)ごと削除されるので、解析環境がmalwareに感染したとしても解析対象のホストには影響を与えません。

## 概要図

```
+--------------------------+      +--------------------------+     +--------------------------+     
| Client                   |      | Client                   |     | Client                   |
| +----------------------+ |      | +----------------------+ |     | +----------------------+ |
| | OpenStack creds info | |      | | OpenStack creds info | |     | | OpenStack creds info | |
| +----------------------+ |      | +----------------------+ |     | +----------------------+ |
+------------+-------------+      +-------------+------------+     +-------------+------------+
             |                                  |                                |
             +----------------------------------+--------------------------------+
                                                |              
                           +--------------------|----------------+
                           | Control            v                |
                           |            +------------+           |
                           |            | RestFulAPI |           |
                           |            +------------+           |
                           |            | DataBase   |           |
                           |            +----+-------+           |
                           |                 ^                   |
                           +-----------------|-------------------+
                                             |
                            +----------------+-----------------+
                            |                |                 |
                   +--------|----------------|-----------------|----------+
                   | Jail   |                |                 |          |
                   | +-------------+   +-------------+   +-------------+  |
                   | | SANDBOX     |   | SANDBOX     |   | SANDBOX     |  |
                   | | +---------+ |   | +---------+ |   | +---------+ |  |
                   | | | HTTPD   | |   | | HTTPD   | |   | | HTTPD   | |  |
                   | | +---------+ |   | +---------+ |   | +---------+ |  |
                   | | | MOD_PHP | |   | | MOD_PHP | |   | | MOD_PHP | |  |
                   | | +---------| |   | +---------| |   | +---------| |  |
                   | +-------------+   +-------------+   +-------------+  |
                   +------------------------------------------------------+
```

## 初期設定

###  Controlサーバの作成

Controlサーバは動的解析処理のセッション管理、Jail環境の構築、動的解析処理の効率化のためのキャッシングなどを担います。
ControlサーバをOpenStack上に作成します。私はOpenStackProviderはConohaを利用しています。

 * OpenStackProviderにComputeNodeをデプロイする

ansibleでデプロイを行いますので、事前にansibleをインストールしておいてください。

```
$ sudo apt-get install -y libbz2-dev zlib1g-dev libssl-dev libreadline-dev libsqlite3-dev
$ git clone https://github.com/yyuu/pyenv.git $HOME/.pyenv
$ git clone https://github.com/yyuu/pyenv-virtualenv.git $HOME/.pyenv/plugins/pyenv-virtualenv
$ echo 'export PATH="$HOME/.pyenv/bin:$PATH"' >> $HOME/.bashrc
$ echo 'eval "$(pyenv pyenv init -)"' >>  $HOME/.bashrc
$ echo 'eval "$(pyenv virtualenv-init -)"' >>  $HOME/.bashrc
$ source $HOME/.bashrc
$ pyenv install 2.7.11
$ pyenv global 2.7.11
$ pip install --upgrade pip
$ pip install ansible
```

compute nodeのdistroはdebian、architectureはx86/64で作成してください。

```
$ sudo -c "echo '163.44.172.244 webshkiller-control.openstack' >> /etc/hosts"
$ ssh-keygen -t rsa -C webshkiller -f ./login_user
$ ansible-playbook -i inventory site.yml --private-key=./webshkiller.pem
```

sshでログインします。

```
ssh webshkiller@webshkiller-control.openstack -i login_user
```

## 処理手順

 * webshdd-cliはwebshdd-coreが存在するか否かの確認を行う
   * webshdd-cliはwebshdd-coreが存在しない場合、webshdd-coreを作成する

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

## How to use

```
export OPENSTACK_IDENTITY_URI="https://identity.tyo1.conoha.io/v2.0"
export OPENSTACK_USERNAME="XXX"
export OPENSTACK_PASSWORD="XXX"
export OPENSTACK_TENANT_ID="XXX"
```

# 参考文献

 * [WebSHArk 1.0: A Benchmark Collection for Malicious Web Shell Detection](https://pdfs.semanticscholar.org/d2de/06d1e4e07890c9b27bdb4baa07c1922b3c16.pdf)
 * https://github.com/b374k/b374k
 * https://github.com/openstack/golang-client
 * https://hirokiaramaki.com/2016/01/17/go-openstack/
 * http://gophercloud.io/

