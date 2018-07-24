const go2dts = require("../src/index");
const rimraf = require("rimraf");
const { join } = require("path");
const { readFileSync, readdirSync } = require("fs");

beforeAll(next => {
  rimraf(join(__dirname, "./outputs"), () => {
    go2dts(join(__dirname, "./inputs"), join(__dirname, "./outputs"));
    next();
  });
});

describe("go2dts", () => {
  readdirSync(join(__dirname, "./outputs"))
    .filter(fileName => /^[a-z]+(?!test)\.d.ts$/.test(fileName))
    .forEach(fileName => {
      describe(fileName, () => {
        it("should match the snapshot", () => {
          expect(
            readFileSync(join(__dirname, "./outputs", fileName), "utf-8")
          ).toMatchSnapshot();
        });
      });
    });
});
