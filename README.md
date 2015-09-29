
Sanat
=====

Translates a string resource translation file that looks something like this:

    LoginView.Title
        en = Log in
        fi = Kirjaudu sisään

…into string resource files appropriate for use on several different software platforms.

Run the main program with the `--help` argument to see “usage” information.


Translation File Syntax
------------------------

### Section titles

Section titles are prepended by at least three `=` symbols:

    === Login view

Both leading and trailing `=` symbols are okay:

    ======== Login view ========

Sections are optional — translation files don't have to contain sections.


### Translations

__Translation keys__ (e.g. `LoginView.Title` below) are not indented.
__Translation values__ are indented by at least two spaces:

    LoginView.Title
        en = Log in
        fi = Kirjaudu sisään

Each __translation value__ line must begin with an _ISO 639-1 language name_, followed by a `=` sign, followed by the actual text content of the translation (for the specified language.)

The translation text content may be double quoted:

        fi = "Kirjaudu sisään "

#### Platform limits

Translations can be limited to certain platforms like so:

    LoginView.Title
        en = Log in
        fi = Kirjaudu sisään
        platforms = apple, android

Translations that specify platforms will only be rendered in the translation output files for those platforms (and not for others.)

The currently supported values are:

- `apple` (Apple platforms; iOS and OS X)
- `android`


### Format specifiers

Translation text can contain format specifiers like this:

    Hello {s}, today it’s {f.2} degrees celsius.

Format specifiers are delineated by the `{}` signs, and they can contain the following:

    { 3: f .2 }
      ^^ ^ ^^____ (optional) decimal count
       |  \___ data type
        \___ (optional) order index

The __data type__ can be one of the following:

- `@`: Object
- `s`: String
- `d`: Integer
- `f`: Floating-point number

(These are mapped to the closest corresponding platform-specific format specifiers.)

The __decimal count__ specifies the number of decimals to show for floating-point numbers, and is only relevant if the data type is `f`.

The __order index__ specifies the 1-based index of the “printf argument” to apply for this format specifier. (This is necessary for cases where the word order for the same sentence differs between languages.)


Preprocessors
-------------

Translations can be preprocessed by specifying one of the following preprocessors:

### `markdown`

Translates the text in each translation from Markdown to HTML.



License
-------

Please see the `LICENSE` file.





