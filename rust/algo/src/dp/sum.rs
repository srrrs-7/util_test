
pub fn partial_sum(seek_weight: i32, weights: Vec<i32>) -> i32 {
    // define DP table
    let mut dp: Vec<Vec<i32>> = Vec::new();

    // create DP table
    for _ in 0..weights.len() {
        let mut cols: Vec<i32> = vec![0; (seek_weight+1) as usize];
        cols[0] = 1;
        dp.push(cols);
    }

    // dp logic
    for row in 1..weights.len() {
        for col in 0..seek_weight+1 {
            if weights[row] <= col {
                dp[row][col as usize] += dp[row-1][(col-weights[row]) as usize] + dp[row-1][col as usize]
            } else {
                dp[row][col as usize] = dp[row-1][col as usize]
            }
        }
    }

    println!("{:?}", dp);

    dp[(weights.len()-1) as usize][seek_weight as usize]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_partial_sum() {
        assert_eq!(partial_sum(1, vec!(2,4,7,3,12)), 3);
    }
}