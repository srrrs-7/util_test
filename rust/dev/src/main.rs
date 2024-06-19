
use algo::add;
use algo::sum;

fn main() {
    let total = add(1,2);
    println!("{}", total);

    let s = sum(16,vec![2,4,7,3,12]);
    println!("{}", s);

    let ss = sum(10,vec![7,3,4,3,2]);
    println!("{}", ss);

    life_time();
}

fn life_time() {
    let n: i32;
    {
        let a = 5;
        n = a;
    }
    println!("{}", n);
}