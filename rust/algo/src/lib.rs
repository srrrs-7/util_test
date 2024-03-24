
mod dp;

pub fn add(left: usize, right: usize) -> usize {
    left + right
}

pub fn sum(left: i32, right: i32) -> i32 {
    dp::sum::partial_sum(left, right)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}
