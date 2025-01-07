import { add, sub } from "../core/lib/wasm_pkg/lib.js";

const sum = add(5, 3);
console.log(`5 + 3 = ${sum}`);

const difference = sub(10, 4);
console.log(`10 - 4 = ${difference}`);