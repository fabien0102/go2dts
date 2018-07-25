#!/usr/bin/env node

const program = require("commander");
const go2dts = require("../src/index");
const { readFileSync } = require("fs");
const { join } = require("path");

const package = JSON.parse(
  readFileSync(join(__dirname, "../package.json"), "utf-8")
);

program
  .version(package.version)
  .usage("<inputDirs ...> <outputDirOrFile>")
  .action((...args) => {
    const currentDir = process.cwd();
    const inputDirs = args.slice(0, -2).map(i => join(currentDir, i));
    const outputDirOrFile = join(currentDir, args[args.length - 2]);

    go2dts(inputDirs, outputDirOrFile);
    console.log(`Types definition created into ${outputDirOrFile}`);
  })
  .parse(process.argv);
