const go2dts = require("../src/index");
const rimraf = require("rimraf");
const { join } = require("path");
const { readFileSync, readdirSync } = require("fs");

beforeAll(next => {
  rimraf(join(__dirname, "./outputs"), () => {
    go2dts(
      [
        join(__dirname, "./inputs/client"),
        join(__dirname, "./inputs/types"),
        join(__dirname, "./inputs/labsserver/httputils")
      ],
      join(__dirname, "./outputs/labs.types.d.ts")
    );
    next();
  });
});

describe("go2dts", () => {
  it("should match the snapshot", () => {
    expect(
      readFileSync(join(__dirname, "./outputs/labs.types.d.ts"), "utf-8")
    ).toMatchSnapshot();
  });
});
