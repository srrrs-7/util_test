
use algo::add;
use algo::partial_sum::partial_sum;

fn main() {
    let sum = add(1,2);
    println!("{}", sum);

    let s = partial_sum(1,2);
    println!("{}", s);
}
