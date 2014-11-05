# The Tavor format

The [Tavor](/) format is an [EBNF-like notation](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_Form) which allows the definition of data (e.g. file formats and protocols) without the need of programming. It is the default format of the [Tavor framework](/) and supports every feature which the framework currently provides.

The format is Unicode text encoded in UTF-8 and consists of terminal and non-terminal symbols which are called <code>tokens</code> throughout the Tavor framework. An explanation of the general meaning can be found in the [What are tokens?](/#token) section.

Every example throughout this page is a complete Tavor format file. The content of each example can be for instance saved into a file called <code>file.tavor</code> and then fuzzed with the Tavor binary.

```bash
tavor --format-file file.tavor fuzz
```

Since some examples have more than one permutation, meaning there is more than one possible fuzzing generation, it is advisable to use the <code>AllPermutations</code> fuzzing strategy to print out every possible permutation of the fuzzed format.

```bash
tavor --format-file file.tavor fuzz --strategy AllPermutations
```

## <a name="table-of-content"></a>Table of content

- [Token definition](#token-definition)
- [Terminal tokens](#terminal-tokens)
	+ [Numbers](#terminal-tokens-numbers)
	+ [Strings](#terminal-tokens-strings)
- [Concatenation](#concatenation)

TODO update this

## <a name="token-definition"></a>Token definition

Every token in the format belongs to a non-terminal token definition which consists of a unique case-sensitive name and its definition part. Both are separated by exactly one equal sign. Syntactical white spaces are ignored. Every token definition must be declared by default in one line. A line ends with a new line character.

To give an example, the following format declares the token <code>START</code> with the constant string token "Hello World" as its definition.

```tavor
START = "Hello World"
```

Token names have the following rules:
- Token names have to start with a letter.
- Token names can only consist of letters, digits and the underscore sign "_".
- Token names have to be unique in the format definition scope.

Additional to these rules it is not allowed to declare a token without any usage in the format definition scope except if it is the <code>START</code> token which is used as the entry point of the format, meaning it defines the beginning of the format. Hence, it is required for every format definition.

## <a name="terminal-tokens"></a>Terminal tokens

Terminal tokens are the constants of the Tavor format.

### <a name="terminal-tokens-numbers"></a>Numbers

Currently only positive decimal integers are allowed. They are written as a sequence of digits.

```tavor
START = 123
```

### <a name="terminal-tokens-strings"></a>Strings

Strings are character sequences between double quotes and can consist of any UTF8 encoded character except new lines, the double quote and the backslash which have to be escaped with a backslash.

```tavor
START = "The next word is \"quoted\" and here is a new line\n"
```

Since Tavor is using Go's text parser as foundation of its format parsing, the same rules for <code>interpreted string literals</code> apply. These rules can be looked up in [Go's language specification](https://golang.org/ref/spec#String_literals).

## <a name="concatenation"></a>Concatenation

Tokens in the definition part are automatically concatenated.

```tavor
START = "This is a string token and this " 123 " was a number token"
```

This example will be concatenated to the string "This is a string token and this 123 was a number token".

## <a name="multi-line"></a>Multi line token definitions

A token definition can be sometimes too long or poorly readable. It can be therefore split into multiple lines by using a comma before the newline character.

```tavor
START = "This",
        "is",
        "a",
        "multi line",
        "definition"
```

The token definition ends at the string "definition" since there is no comma before the new line character. This example also underlines that syntactical white spaces are ignored and can be used to make the format definition more human readable.

## <a name="comments"></a>Comments

The comments of the Tavor format follow the same rules as Go's comments which are specified in [Go's language specification](https://golang.org/ref/spec#Comments).

There are two types of comments:
- **Line comment** which starts with the character sequence <code>//</code> and ends at the next new line character.
- **General comment** which starts with the character sequence <code>/\*</code> and ends at the character sequence <code>\*/</code>. A general comment can contain new line characters.

```tavor
/*

This is a general comment
which can have
multiple lines

*/

START = "This is a string" // this is a line comment

// this is also a line comment
```

General comments can be used, like white space characters, between token definitions and tokens.

```tavor
START /* this is */ = "an" /* extreme */ "example" /* but
it should make it clear how general comments */ "work"
```

## <a name="embedding"></a>Token embedding

Non-terminal tokens can be embedded in the definition part by using the name of the referenced token. The following example embeds the token <code>String</code> into the <code>START</code> token.

```tavor
START = String

String = "this is a string"
```

Token names declared in the global scope of a format definition can be used throughout the format regardless of their declaration position.

Terminal and non-terminal tokens can be mixed.

```tavor
Dot = "."

First  = 1 Dot
Second = 2 Dot
Third  = 3 Dot

START = First ", " Second " and " Third
```

## <a name="alternation"></a>Alternation

Alternations are defined with the pipe character <code>|</code>. The following example defines that the token <code>START</code> can either hold 1, 2 or 3.

```tavor
START = 1 | 2 | 3
```

An alternation term has its own scope which means that a sequence of tokens can be used.

```tavor
START = 1 "green apple" | 2 "orange oranges" | 3 "yellow bananas"
```

Alternation terms can be empty which allows more advanced definitions of formats. For example the next definition defines the possibility of a loop.

```tavor
A = "a" A | B |
B = "b"

START = A
```

This example can hold for example the strings "", "a", "b", "ab", "aab" or any amount of "a" characters ending with one or no "b" character.

## <a name="grouping"></a>Grouping

Tokens can be grouped using parenthesis. A group starts with <code>(</code> and ends with <code>)</code> and is a token on its own. This means that it can be mixed with other tokens. Additionally, a group starts a new scope between its parenthesis and can therefore hold a sequence of tokens. The tokens between the parenthesis is called the <code>group body</code>.

The following example declares that the token <code>START</code> either holds the string "old news" or "new news".

```tavor
START = ("old" | "new") " news"
```

Groups can be nested too. For example the following can be used to define that the <code>START</code> token can either hold "a", "b", "1" or "2".

```tavor
START = (("a" | "b") | (1 | 2))
```

An even more complicated example is the definition of an one to three digits integer.

```tavor
Digit = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

START = Digit (Digit (Digit | ) | )
```

This could be also written with the following format definition.

```tavor
Digit = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

START = Digit | Digit Digit | Digit Digit Digit
```

Group parenthesis can have modifiers which give the group additional abilities. The following sections will introduce these modifiers.

### <a name="grouping-optional"></a>Optional group

The optional group allows the whole group token to be optional. In the next example the <code>START</code> token can hold the string "funny" or "very funny".

```tavor
START = ?("very ") "funny"
```

### <a name="grouping-repeats"></a>Repeat groups

The default modifier for the repeat group is the plus character <code>+</code>. The repetition is executed by default at least once. In the next example the string "a" is repeated and the <code>START</code> token can therefore hold the strings "a", "aa", "aaa" or any amount of "a" characters.

```tavor
START = +("a")
```

Although the format definition allows the repetition to go on forever there are bounds since there is only a finite amount of memory available. The Tavor framework does also set a maximum repetition by default which can be altered by the <code>--max-repeat</code> option or the <code>MaxRepeat</code> variable in the <code>github.com/zimmski/tavor</code> package.

By default the repetition modifier repeats from one to infinite which can be altered by arguments to the modifier. The next example repeats the string "a" exactly twice meaning the <code>START</code> token does only hold the string "aa".

```tavor
START = +2("a")
```

It is also possible to define a repetition range. The next example repeats the string "a" at least twice but at most 4 times. This means that the <code>START</code> token can either hold the string "aa", "aaa" or "aaaa".

```tavor
START = +2,4("a")
```

The <code>from</code> and <code>to</code> arguments can be empty too which sets them to their default values. For example the next definition repeats the string "a" at most 4 times.

```tavor
START = +,4("a")
```

And the next example repeats the string "a" at least twice.

```tavor
START = +2,("a")
```

Since the repetition zero, once or more is very common the modifier <code>\*</code> exists. In the next example the token <code>START</code> can either hold the string "a", "ab", "abb" or any amount of "b" characters prepended by an "a" character.

```tavor
START = "a" *("b")
```

### <a name="grouping-permutation"></a>Permutation group

The <code>@</code> is the permutation modifier which is combined with an alternation in the group body. Each alternation term will be executed exactly once but the order of execution is non-relevant. In the next example the <code>START</code> token can either hold 123, 132, 213, 231, 312 or 321.

```tavor
START = @(1 | 2 | 3)
```

## <a name="character-classes"></a>Character classes

Character classes are a special kind of token and can be directly compared to character classes of regular expressions used in most programming languages such as Perl's implementation which is documented [here](http://perldoc.perl.org/perlre.html#Character-Classes-and-other-Special-Escapes). They behave like terminal tokens meaning that they cannot include others tokens but they are, unlike integers and strings, not single but multiple constants at once. A character class starts with the left bracket <code>[</code> and ends with the right bracket <code>]</code>. Character classes are like terminal tokens in that they are tokens on their own and can be therefore mixed with other tokens. The content between the brackets is called a pattern and can consists of almost any UTF8 encoded character, escape character, special escape and range. In general the character class token can be seen as a shortcut for a string alternation.

For example the following definition lets the <code>START</code> token hold the strings "a", "b" or "c".

```tavor
START = "a" | "b" | "c"
```

With a character class this can be written as the following.

```tavor
START = [abc]
```

### <a name="character-classes-escapes"></a>Escape characters

The following table holds UTF8 encoded characters which are not directly allowed within a character class pattern. Their equivalent escape sequence has to be used instead.

| Character       | Escape sequence   |
| :-------------- | :---------------- |
| <code>-</code>  | <code>\\-</code>  |
| <code>\\</code> | <code>\\\\</code> |
| form feed       | <code>\\f</code>  |
| newline         | <code>\\n</code>  |
| return          | <code>\\r</code>  |
| tab             | <code>\\t</code>  |

For example the following defines that the <code>START</code> token holds only white space characters.

```tavor
START = +([ \n\t\n\r])
```

Since some characters can be hard to type and read the <code>\x</code> escape sequence can be used to define them with their hexadecimal code points. There are two options to do this. Either only two hexadecimal characters are used in the form of <code>\x0A</code> or more than two hexadecimal digits are needed which have to use the form <code>\x{0AF}</code>. The second form allows up to 8 digits and is therefore fully Unicode ready.

To give an example the following definition holds either the Unicode character "/" or "😃".

```tavor
START = [\x2F\x{1F603}]
```

### <a name="character-classes-ranges"></a>Ranges

Ranges can be defined using the <code>-</code> character. A range holds all characters starting at the character before the <code>-</code> and ending at the character after the <code>-</code>. Both characters have to be either an UTF8 encoded or an escaped character. The starting character must have a lower value than the ending character.

For example the following defines a decimal digit.

```tavor
START = [0123456789]
```

This can be easier defined using a range.

```tavor
START = [0-9]
```

It is also possible to use hexadecimal code points, since either range characters can be escape characters.

```
START = [\x23-\x5B]
```

### <a name="character-classes-special-escapes"></a>Special escape characters

Special escape characters combine many characters into one escape character and can also hold additional functionality. The following table is an overview of all currently implemented special escape characters.

| Special escape character | Character class           | Description                     |
| :----------------------- | :------------------------ | :------------------------------ |
| <code>\d</code>          | <code>[0-9]</code>        | Holds a decimal digit character |
| <code>\s</code>          | <code>[ \f\n\r\t]</code>  | Holds the white space character |
| <code>\w</code>          | <code>[a-zA-Z0-9_]</code> | Holds a word character          |

-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------
-------------

# TODO rewrite everything down below

### Token attributes

Token attributes can be used in token definitions by prepending a dollar sign to their name and separate the token name from the attribute by a dot.

```
Letters = *(Letter)
Letter = "a" | "b" | "c"
LetterCount = $Letters.Count // LetterCount then holds the count of the repeater Letters
```

Possible token attributes are:
* Count - Holds the count of this token. Must be a repeater.
* Index - Holds the index of a token. Must be a token of a repeater.
* Unique - Chooses at random a token of a repeater.

### Special tokens

Special tokens can be defined by prepending a dollar sign to their name. Special tokens do not have a format on their right side like regular tokens, instead arguments written as key-value pairs, which are separated by a colon, define the token. At least the "type" argument must be defined.

```
$Number = type: Int
Arithmetic = Number "+" Number
```

Possible arguments are:
* type - Defines the type of the token. Can be "Int" or "Sequence"

Additional (optional) arguments for each type are:
* "Int"
    * from - First integer value
    * to - Last integer value
* "Sequence"
    * start - First sequence value. Default is 1.
    * step - Increment of the sequence. Default is 1.

Possible attributes for each type are:
* "Int"
    * Value - The value of the Int
* "Sequence"
    * Next - Indicates the next value of the sequence.
    * Existing - Indicates an available value of the sequence in the whole data.
    * Reset - The sequence is reseted when this token is reached.

```
$Id = type: Sequence,
      start: 0,
      step: 2
NextId = $Id.Next
```

### Expressions

Expressions can be used on the right side of a token definition.

```
Sum = ${1 + 2 + 3} // Sum will be interpreted as 6

SomeIdOrMore = $Id.Existing | ${Id.Existing + 1}

DoubleTheCount = ${Letter.Count + Letter.Count}
```

### Variables

Every token on the right side of a definition can be saved into a variable.


```tavor
START = "text"<var> Print

Print = <var>.Value
```

This will save the string "text" into the variable "var" without preventing the relay of the string to the output stream.

Since there are circumstances where a token should be just saved into a variable but not relayed to the output stream a second syntax can be used.

```tavor
START = "text"<=var> Print

Print = <var>.Value
```

### Set operators

Some attributes can be combined with set operators. For example

```
$Id = type: Sequence

Pair = $Id.Next<id> " " ${Id.Existing not in (id)}
```

This will search through the existing sequenced IDs without the one saved in the variable "id".

### If, If else and else

```tavor
START = Choose<var> Print

Choose = 1 | 2 | 3

Print = {if var.Value == 1} "var is one" {else if var.Value == 2} "var is two" {else} "var is three" {endif}
```

### Condition operators

* "=="

  ```tavor
  Print = (1 | 2 | 3)<var> {if var.Value == 1} "var is 1" {else} "var is not 1" {endif}
  ```

* "defined"

  ```tavor
  START = Print "save this text"<var> Print

  Print = {if defined var} "var is: " $var.Value {else} "var is not defined" {endif}
  ```
