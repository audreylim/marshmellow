# Marshmellow

Marshmellow is a simple markdown processor implemented in Go.

## Installation

```
go get github.com/audreylim/marshmellow
```

## Usage

Processing can be done on multiple files. Just run

```
./marshmellow file1.md file2.md file3.md
```

This will turn the markdown files into the corresponding HTML files, ie. `file1.html`, `file2.html`, `file3.html`.

## Syntax

Marshmellow currently supports the following syntax:

##### Headers

```
# Header 1
```

becomes

```
<h1>Header 1</h1>
```

Marshmellow supports up to `h6` headers, ie. `###### Header 6`.

##### Bold

```
**Bold Text**
```

becomes

```
<b>Bold Text</b>
```

##### Italics

```
*Italic Text*
```

becomes

```
<i>Italic Text</i>
```

###### Bullets

```
* Bullet 1
* Bullet 2
* Bullet 3
```

becomes

```
<ul>
<li>Bullet 1</li>
<li>Bullet 2</li>
<li>Bullet 3</li>
</ul>
```
