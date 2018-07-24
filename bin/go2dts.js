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
  .usage("<inputDir> <outputDirOrFile>")
  .action((inputDir, outputDirOrFile) => {
    const currentDir = process.cwd();
    go2dts(join(currentDir, inputDir), join(currentDir, outputDirOrFile));
    console.log(`Types definition created into ${outputDirOrFile}`);
  })
  .parse(process.argv);
