# Marshmellow

Marshmellow is a simple markdown processor implemented in Go.

## Installation

```
go get github.com/audreylim/marshmellow/...
```

## Usage

Conversion can be done on multiple files. Just run

```
./marshmellow file1.md file2.md file3.md
```

This turns the markdown files into corresponding HTML files, ie. `file1.html`, `file2.html`, `file3.html`.

## Syntax

Marshmellow currently supports the following syntax:

##### Headers

```markdown
# Header 1
```

becomes

```html
<h1>Header 1</h1>
```

Marshmellow supports up to `h6` headers, ie. `###### Header 6`.

##### Paragraph

```
Paragraph 1
```

becomes

```html
<p>Paragraph 1</p>
```

##### Bold

```markdown
**Bold Text**
```

becomes

```html
<b>Bold Text</b>
```

##### Italics

```
*Italic Text*
```

becomes

```html
<i>Italic Text</i>
```

##### Bullets

```markdown
* Bullet 1
* Bullet 2
* Bullet 3
```

becomes

```html
<ul>
<li>Bullet 1</li>
<li>Bullet 2</li>
<li>Bullet 3</li>
</ul>
```

## License

MIT
