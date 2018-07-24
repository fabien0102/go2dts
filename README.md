# Go2dts

A simple cli-tools to transform golang `struct` to typescript `interface`

### Installation

```bash
npm i -g go2dts
```

### Usage

```bash
go2dts <goLangDir> <typescriptDir>
```

### Testing

Just put your golang file into `__tests__/inputs` and it will be parse each time you execute `npm test`.

By default, every output files have a jest snapshot associate, so you can update the broken snapshot to what you want and start developing ;)

### Know issues

- The type mapping is incomplete (I'm not a golang developper, so I add types when I discover them)
- This library will not follow the `import` dependencies, so the output types can be broken
