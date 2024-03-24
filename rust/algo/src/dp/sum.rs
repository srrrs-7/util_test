
pub fn partial_sum(lhs: i32, rhs: i32) -> i32 {
    lhs + rhs
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_partial_sum() {
        assert_eq!(partial_sum(1, 2), 3);
    }
}