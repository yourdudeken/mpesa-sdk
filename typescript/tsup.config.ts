import { defineConfig } from "tsup";

export default defineConfig({
  entry: {
    index: "src/index.ts",
    "webhooks/index": "src/webhooks/index.ts",
    "errors/index": "src/errors/index.ts",
    "types/index": "src/types/index.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  splitting: true,
  sourcemap: true,
  clean: true,
  treeshake: true,
  minify: false,
  tsconfig: "./tsconfig.json",
});
