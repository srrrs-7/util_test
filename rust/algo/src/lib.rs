mod dp;

pub fn add(left: usize, right: usize) -> usize {
    left + right
}

pub fn sum(seek_weight: i32, weights: Vec<i32>) -> i32 {
    dp::sum::partial_sum(seek_weight, weights)
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
