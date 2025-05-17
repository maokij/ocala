# Ocala 概説

Ocala は 8 bit CPU/MPU をターゲットとしたアセンブラです。
各種の構造化機構と対象文脈式を特徴とします。
現在 Z80 および 6502 のバイナリ出力に対応しています。

## 手続きとコード

Ocala ではコードは手続きの中に記述します。

```
  proc <proc-name>(<signature>) { ... }
```

コードは以下の 2 種類の記法で記述できます。

- 機械語命令記法(直接記法)
- 対象文脈記法(中置記法)

機械語命令記法は、CPU 依存の機械語命令を直接記述する形式です。
ただし、`(...)`は定数の記述となるため、メモリ参照には`[...]`を利用します。
また、オペランドの区切りの`,`は省略できます。

```
  LD A 10
  LD A [HL]
  INC A
```

対象文脈記法は、対象となるオペランドを指定し、その文脈での操作を列記する記法です。

```
  A             // Specify the target operand.
    <- 10       // LD  A 10
    +  20       // ADD A 20
    -> [0x1234] // LD [0x1234] A

  A <- 10 + 20 -> [0x1234]  // Same above.
```

手続きにはシグネチャを指定できます。
シグネチャは手続き内でのレジスタの用法に関する注記です。
手続きの呼び出し時はこのシグネチャも記述する必要があります。
意図を示す注記であり、実際の用法を検査するものではありません。

```
  proc f(A B HL => A ! BC E) { // In:     A, B and HL
      ...                      // Out:    A
  }                            // Modify: BC and E.

  proc g(=> A L !) {           // In:     none
      ...                      // Out:    A and L
  }                            // Modify: all others

  proc main() {
      f(A B HL => A ! BC E)    // OK
      f()                      // ?? proc signature mismatch: f.
      g(=> A L !)              // OK
      ...
  }
```

手続きの末尾は `return` / `goto` もしくはその派生である必要があります。

```
  // Valid
  proc valid-a() {
    return          // OK
  }
  proc valid-b() {
    goto valid-c    // OK
  }
  proc valid-c() {
    fallthrough     // OK: same as "goto valid-d"
  }
  proc valid-d() {
    recur           // OK: same as "goto valid-d"
  }
  proc valid-e() {
    never-return loop { NOP } // OK: never return
  }

  // Invalid
  proc invalid-a() {
    NOP              // ?? the last instruction must be a return/fallthrough within the proc
  }
  proc invalid-b() {
    fallthrough      // ?? it is the last proc
  }
```

手続き名など、プログラム上で扱う名前(識別子)には、
英数字(`a-zA-Z0-9`)と一部の記号(`-^!$%&*+/<=>?|~_`)を使用できます。

`+` や `-` などの演算子も名前(識別子)として扱われるので、空白などを用いずに詰めて
記述すると別の識別子として扱われます。
例えば、単項マイナス演算子(`-`)は括弧を用いる必要があります。

```
proc f-g() { return } // the procedure named `f-g'
const f+g = 1 // OK: the constant named `f+g'

A + 1 // register `A', operator `+', number `1'
A+1 // ?? unknown form name A+1

L001:
HL <- -(L001) // OK
HL <- -L001 // ?? undefined name -L001
```

プログラム中では、`,` 記号は通常空白文字と同様に扱われます。
ただし、2 個以上連続して記述することはできません。
また、行末に記述すると行の継続を示すことができます。

```
LD A 1  // OK
LD A, 1 // OK
LD A,, 1 // ??
LD A,
   1 // OK: LD A 1
```

## メモリ参照

メモリの参照は `[...]` を使用します。
6502 であってもメモリ参照には括弧を使用する必要があります。

```
  LD A 12
  LD HL [0xfffe]
  C <- [HL@1024] // LD HL, 1024; LD C, (HL)

  LDA 1          // LDA #1
  LDA [8]        // LDA 8
  LDA [[8] Y]    // LDA (8), Y
```

## 対象文脈式演算子

対象文脈では下記の操作を記述できます(_ は文脈、% は引数)。

| 操作        | Z80                           | 6502                            |
|-------------|-------------------------------|---------------------------------|
| `<- %`      | `LD _ %`                      | `T%_ / LD_ %`                   |
| `-> %`      | `LD % _`                      | `T_% / ST_ %`                   |
| `<-> %`     | `EX _ %`                      | -                               |
| `+ %`       | `ADD _ %`                     | `CLC; ADC %`                    |
| `+$ %`      | `ADC _ %`                     | `ADC %`                         |
| `- %`       | `SUB %` /<br> `OR A; SBC _ %` | `SEC; SBC %`                    |
| `-$ %`      | `SBC _ %`                     | `SBC %`                         |
| `-? %`      | `CP _ %`                      | `CM_ %`                         |
| `& %`       | `AND %`                       | `AND %`                         |
| `\| %`      | `OR %`                        | `ORA %`                         |
| `^ %`       | `XOR %`                       | `EOR %`                         |
| `<* %`      | (`RLCA` / `RLC _`) * %        | (`CMP 0x80; ROL A`) * %         |
| `<*$ %`     | (`RLA` / `RL _`) * %          | (`ROL _`) * %                   |
| `>* %`      | (`RRCA` / `RRC _`) * %        | (`LSR A; BCC +2; ORA 0x80`) * % |
| `>*$ %`     | (`RRA` / `RR _`) * %          | (`ROR _`) * %                   |
| `<< %`      | (`SLA _`) * %                 | (`ASL _`) * %                   |
| `>> %`      | (`SRA _`) * %                 | (`CMP 0x80; ROR A`) * %         |
| `>>> %`     | (`SRL _`) * %                 | (`LSR _`) * %                   |
| `-set %`    | `SET % _`                     | -                               |
| `-reset %`  | `RES % _`                     | -                               |
| `-bit? %`   | `BIT % _`                     | `BIT`                           |
| `-in %`     | `IN _ %`                      | -                               |
| `-out %`    | `OUT % _`                     | -                               |
| `++`        | `INC _`                       | `CLC; ADC 1` / `IN_` / `INC _`  |
| `--`        | `DEC _`                       | `SEC; SBC 1` / `DE_` / `DEC _`  |
| `-push`     | `PUSH _`                      | `PH_`                           |
| `-pop`      | `POP _`                       | `PL_`                           |
| `-not`      | `CPL`                         | `EOR 0xff`                      |
| `-neg`      | `NEG`                         | `EOR 0xff; CLC; ADC 1`          |
| `-zero?`    | `AND A` /<br>`INC _; DEC _`   | -                               |
|             |                               |                                 |
| `@%`        | alias of `<-`                 |                                 |
| `. { ... }` | inline sub context            |                                 |

## 制御構造

手続き内では、以下の制御構造を使用できます。

| 制御構造                         | 処理内容                     |
|----------------------------------|------------------------------|
| **基本ブロック**                 |                              |
| `do { ... }`                     | ブロック評価                 |
| `<label>:`                       | ラベル定義                   |
| **分岐**                         |                              |
| `if <cond> { ... } else { ... }` | 条件分岐                     |
| `goto <label>`                   | ジャンプ                     |
| `goto-if <cond> <label>`         | 条件ジャンプ                 |
| `goto-rel <label>`               | (Z80) 相対ジャンプ           |
| `goto-rel-if <cond> <label>`     | (Z80) 相対条件ジャンプ       |
| **ループ制御**                   |                              |
| `loop { ... }`                   | 無限ループ                   |
| `loop { ... } while <cond>`      | 条件ループ                   |
| `once { ... }`                   | 1度のみのループ              |
| `redo`                           | ループ先頭へジャンプ         |
| `redo-if <cond>`                 | ループ先頭へ条件ジャンプ     |
| `continue`                       | ループ条件部へジャンプ       |
| `continue-if <cond>`             | ループ条件部へ条件ジャンプ   |
| `break`                          | ループ中断                   |
| `break-if <cond>`                | 条件付ループ中断             |
| **手続き制御**                   |                              |
| `return`                         | リターン                     |
| `return-if <cond>`               | (Z80) 条件リターン           |
| `recur`                          | 手続き先頭へジャンプ         |
| `fallthrough`                    | リターンせず次の手続きを評価 |
| `never-return loop { ... }`      | 手続き末で無限ループ         |

条件 `<cond>` として、以下を使用できます。

| Z80   | Z80 別名           | 6502  | 6502 別名                    |
|-------|--------------------|-------|------------------------------|
| `NZ?` | `!=?` `not-zero?`  | `NE?` | `!=?` `not-zero?`            |
| `Z?`  | `==?` `zero?`      | `EQ?` | `==?` `zero?`                |
| `NC?` | `>=?` `not-carry?` | `CC?` | `<?` `not-carry?` `borrow?`  |
| `C?`  | `<?` `carry?`      | `CS?` | `>=?` `carry?` `not-borrow?` |
| `PO?` | `odd?` `not-over?` | `VC?` | `not-over?`                  |
| `PE?` | `even?` `over?`    | `VS?` | `over?`                      |
| `P?`  | `plus?`            | `PL?` | `plus?`                      |
| `M?`  | `minus?`           | `MI?` | `minus?`                     |

## インライン手続き

手続きシグネチャの先頭に `-*` を指定すると、インライン手続き定義となります。

```
  proc i(-* A => !) { // Inline
    ...
  }

  proc main() {
      i(-* A => !)    // OK
      i(A => !)       // ?? signature mismatch
  }
```

## 定数

Ocala では、定数を定義することができます。

```
  const <const-name> = <constexpr>
```

定数式では他の定数やラベルを参照できます。また、各種の演算子や括弧も利用可能です。

```
  data ROM_ADDR = 0x4000
  data RAM_ADDR = ROM_ADDR + 0x8000
```

## 定数式演算子

定数式では、以下の演算子を利用可能です。

| operator          |
|-------------------|
| `*` `/` `%`       |
| `+` `-`           |
| `<<` `>>` `>>>`   |
| `<` `<=` `>` `>=` |
| `==` `!=`         |
| `&`               |
| `\|`              |
| `^`               |
| `&&`              |
| `\|\|`            |

## データ

Ocala では、データを定義することができます。

```
  data <data-name> = <type> [ <constexpr>... ] * <repeat> : <section-name>
```

`<type>` は byte(1 バイト単位定義) または word(2 バイト単位定義) を指定できます。

`<constexpr>` はデータの要素となる定数式です。空の場合は省略できます。

`<repeat>` は繰り返し回数を指定します。省略時のデフォルトは 1 です。

`<section>` はデータの配置セクションをしていします。
モジュールの現在のセクションと同じであれば省略できます。

```
  data str = byte [ "hello!" ]
  data tab = byte [ 0 1 2 4 8 ] : rodata
  data dat = word * 10 : bss
```

## モジュール

Ocala ではプログラムは複数の「モジュール」から構成されます。

```
  module <module-name> { ... }
```

モジュールは以下の 2 つの側面を持ちます。

- 名前空間
- セクション(バイナリの配置区画)の集合

モジュールは名前空間です。
Ocala では定数や手続きなどの名前は字句的に解決されます。ただし、名前空間を指定することで
他のモジュールに属する名前を参照することもできます。

```
  module mod-a {
    const c = 1
    const d = c + 1
  }

  module mod-b {
    const c = mod-a:d // Use the constant 'd' in the module 'mod-a'
  }
```

セクションはコードやデータの集まりです。
出力バイナリ上ではセクション単位でコードとデータが配置されます。
モジュール定義直下では任意の位置でセクションの開始を宣言できます。

```
  section <section-name>
```

標準で各モジュールごとに以下のセクションを持ちます。

- text(コード領域)
- rodata(読込専用データ領域)
- bss(未初期化データ領域)

モジュールや手続きは以下の要素から構成されます。

- データ定義文
- 機械語命令文
- 対象文脈式文
- その他の特殊形式文(疑似命令)

## リンク

各モジュールの各セクションはコード生成前に内部中間表現レベルで
「リンク」処理され結合されます。

リンク処理の内容は link 特殊形式を利用して指定します。

```
  link { <link-directive>... }

  <link-directive>:
    org <address> <size> <mode>
    merge <section-name> <module-name>...
```

例えば、以下の様に text セクションと bss セクションをリンクします。

```
  link {
     org 0x4000 0x8000 2  // Set the current address to 0x4000.
                          // Set the region size to 0x8000 bytes.
                          // Output zero-padded binary(mode 2).
     merge text main _    // Merge all text sections.
                          // Start with the main module.
     merge rodata main _  // Merge all rodata sections.

     org 0xc000 0x4000 0  // Set the current address to 0xc000.
                          // Set the region size to 0x4000 bytes.
                          // Do not output this region(mode 0).
     merge bss main _     // Merge all bss sections.
  }
```

一般にリンク処理の内容は対象機種毎にほぼ定型であるため、
ライブラリ側でリンク方法をマクロ定義できます。
そのため、実際にはプログラム側では link 特殊形式を直接記述せずに
ライブラリに定義されたマクロを利用します。

```
include "msx.oc"
msx:link-as-rom main _
```

## マクロ

マクロを定義することで、複数の処理をまとめることができます。

```
  macro <macro-name> ( <param>... ) [ <var>... ] { ... }
```

マクロではマクロ引数を使用できます。
`<name>: <default>` の形式で省略可能な引数を指定できます。
最後の引数に `...` を指定すると可変長引数となります。

また、ラベル等に利用できるマクロ変数を使用できます。
マクロ変数はマクロ外と重複しない名前に展開されます。

```
  macro myloop(body) [beg] {
    %=beg: do %=body
    goto beg
  }

  myloop { NOP }
```

引数、変数は下記のプレースホルダーを利用して展開できます。

| プレースホルダー | 展開内容                     |
|------------------|------------------------------|
| %=               | そのままの値                 |
| %*               | 可変長引数の埋め込み         |
| %#               | 可変長引数の個数             |
| %&               | 定数式内へ非定数式を埋め込み |

## その他の補助機能

その他以下の補助機能を使用可能です。

| 制御構造                    | 処理内容                             |
|-----------------------------|--------------------------------------|
| **スタック**                |                                      |
| `push* <reg>...`            | 一括プッシュ                         |
| `pop* <reg>...`             | 一括ポップ                           |
| `push/pop <reg>... { ... }` | レジスタを一時退避してブロックを評価 |
| **割り込み**                |                                      |
| `di/ei { ... }`             | (Z80) DI、ブロック評価、EI           |
| **フラグ**                  |                                      |
| `set-carry(*-)`             | キャリーフラグセット                 |
| `clear-carry(*-)`           | キャリーフラグクリア                 |
| `clear-over(*-)`            | (6502) オーバーフローフラグクリア    |
