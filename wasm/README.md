# wasm projects

## build command

```sh
cd /wasm/core/lib
rustup target add wasm32-unknown-unknown
# build
cargo build --target wasm32-unknown-unknown --release
# output
/wasm/core/lib/target/wasm32-unknown-unknown/release/lib.wasm
# wasm to deno
cargo install wasm-bindgen-cli
wasm-bindgen target/wasm32-unknown-unknown/release/lib.wasm --out-dir wasm-pkg --typescript --target deno
# deno run
deno run --allow-read --allow-net main.ts
```
