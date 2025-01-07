use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub extern "C" fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[wasm_bindgen]
pub extern "C" fn sub(a: i32, b: i32) -> i32 {
    a - b
}
