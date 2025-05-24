use std::env;
use std::fs;
use std::io::{self, BufRead};

fn main() -> Result<(), anyhow::Error> {
    let args = env::args().collect::<Vec<_>>();
    if args.len() < 2 {
        eprintln!("Usage: {} <name>", args[0]);
        panic!("Invalid number of arguments");
    }

    let target_file = &args[1];
    let target_method = &args[2];

    println!(
        "Searching for method '{}' in file '{}'",
        target_method, target_file
    );

    let finder = MethodFinder::new(target_file, target_method);
    let method = finder.find_method()?;

    println!("Method found: {}", method.name);
    println!("Arguments: {:?}", method.args);
    println!("File path: {}", method.file_path);
    println!("Method body: {}", method.body);
    println!("Line number: {}", method.line_number);

    Ok(())
}

struct Method {
    name: String,
    args: Vec<String>,
    file_path: String,
    body: String,
    line_number: usize,
}

impl Method {
    fn new(
        name: String,
        args: Vec<String>,
        file_path: String,
        body: String,
        line_number: usize,
    ) -> Self {
        Method {
            name,
            args,
            file_path,
            body,
            line_number,
        }
    }
}

struct MethodFinder {
    file_path: String,
    method_name: String,
}

impl MethodFinder {
    fn new(file_path: &str, method_name: &str) -> Self {
        MethodFinder {
            file_path: file_path.to_string(),
            method_name: method_name.to_string(),
        }
    }

    fn find_method(&self) -> Result<Method, anyhow::Error> {
        let file = fs::File::open(&self.file_path)?;
        let reader = io::BufReader::new(file);
        let mut in_method = false;
        let mut method_body = String::new();
        let mut args = Vec::new();
        let mut line_number = 0;

        for line in reader.lines() {
            line_number += 1;
            let line = line?;
            if line.contains(&self.method_name) {
                in_method = true;
                method_body.push_str(&line);
                // Extract arguments from the method signature
                if let Some(args_str) = line.split('(').nth(1).and_then(|s| s.split(')').next()) {
                    args.extend(args_str.split(',').map(|s| s.trim().to_string()));
                }
            } else if in_method && line.contains("}") {
                method_body.push_str(&line);
                break;
            } else if in_method {
                method_body.push_str(&line);
            }
        }

        if in_method {
            Ok(Method::new(
                self.method_name.clone(),
                args,
                self.file_path.clone(),
                method_body,
                line_number,
            ))
        } else {
            Err(anyhow::anyhow!("Method not found"))
        }
    }
}
