# ini-merge

Overlay one INI file onto another and write the merged document (CLI).

## Usage

```
ini-merge -base base.ini -over overlay.ini [-out -]
```

- `-base`  base INI (required)
- `-over`  overlay INI (required); matching keys replace, new keys/sections append
- `-out`   output path, `-` means stdout

UTF-8 BOM is stripped. Lines starting with `#` or `;` are comments.
A line that is not `[section]` or `key=value` is an error.
