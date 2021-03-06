<img src="public/logo.png" width=100 height=100>

# Go 言語でつくるインタプリタ

主に Go 言語の学習をメインとした写経

## Go

### interface

メソッドの型だけを定義した型の一種であり、任意の型がどのようなメソッドを実装するべきかを規定

### struct

Class という概念がない代わりとしての役割を果たす

## 構文解析器

_重要トピック_

- 2.6.9 中置演算子

2 つ戦略がある

- トップダウン構文解析
  - 再帰下降構文解析 // こんかいはこれ Pratt 構文解析器
  - アーリー法
  - 予測的構文解析
- ボトムアップ構文解析

### 文(statement)の構文解析

文が一番基本的な解析
=> 式
=> その他のサブセット

let, return

文と式の違いについて
式は値を生成し、文はしない

let x = 10; は式は 10 で、 10 は値を生成し、`let x = 10` は値を生成しない
という区別を今回は認識する(言語等でそれぞれ違う解釈が存在する)

<img src="public/simple-ast.png" width=250 height=200>

### LetStatement インターフェイス

構成としては `let x = 5;` で考えると楽

- 定義する変数の名前(x) => Name(Identifier 構造体)
- 定義した変数に挿入される値(5) => Value

### ReturnStatement インターフェース

`return <expression>` という構造になる

## 式の構文解析 (ヘビー)

Monkey では let/return 文以外はすべて式となる

- 四則演算では優先順位が必要であり、()等でも優先順位の上昇等がある
- add(1,2) では関数呼び出しだが、 (-1 - 10)では計算のグループ化となる
- 前置演算子: -5/!true/!false
- 二項演算子: 5 + 5/ 5 - 5/5 \* 5
- 呼び出し式: add(2, 3)/max(5, add(5, (5 \* 5)))
- 識別子: foo \* bar/ foobar
- 関数リテラル: fn(x, y) { return x + y };
- if 式: let result = if (10 > 5) { true } else { false };

## Monkey の構文解析器

Vaughan Pratt のトップダウン演算子順位解析(Pratt 構文解析)を利用する

> とてもシンプルで理解しやすく、実装も容易で、使いやすく、理論上はともかく
> 実用上は非常に効率的であり、それでいて、利用者の殆どの合理的な構文の要求に耐えられるほど柔軟である

Pratt は構文解析関数(parseLetStatement メソッド等)を単一のトークンタイプに関連付ける
この方式の肝となるのは、それぞれのトークンタイプに対して 2 つの構文解析関数を関連付ける点
これはトークンの配置、つまり中置か前置かによる

### 用語

- 前置演算子(prefix operator) => --5(デクリメント)
- 後置演算子(postfix operator) => foobar++
- 中置演算子(infix operator) = 5 \* 8

### AST

文では let/return 文があったが、これに 3 つ目として 式文を追加する
式文 => 1 つの式だけからなる文のこと(式のラッパーにすぎない)

```
// 例
let x = 5; // let文
x + 10; // 式文(Monkeyではこれが有効)
```

### 中置演算子

構造

```
<expression> <infix operator> <expression>
```
