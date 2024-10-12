use algo::add;
use algo::sum;
use tokio::sync::mpsc;
use tokio::time::{self, Duration};

// 異なる型のメッセージを扱うための列挙型を定義
enum Message {
    Text(String),
    Number(i32),
}

#[tokio::main]
async fn main() {
    // 単一のチャネルを作成し、送信用と受信用のハンドルを取得
    let (tx, mut rx) = mpsc::channel(32);

    // 非同期タスクで文字列メッセージを送信
    let tx_clone = tx.clone();
    tokio::spawn(async move {
        loop {
            // 文字列メッセージを送信。送信に失敗した場合はループを終了
            if tx_clone.send(Message::Text("Hello from string channel".to_string())).await.is_err() {
                break;
            }
            // 1秒間待機
            time::sleep(Duration::from_secs(1)).await;
        }
    });

    // 非同期タスクで整数メッセージを送信
    let tx_clone2 = tx.clone();
    tokio::spawn(async move {
        loop {
            // 整数メッセージを送信。送信に失敗した場合はループを終了
            if tx_clone2.send(Message::Number(42)).await.is_err() {
                break;
            }
            // 2秒間待機
            time::sleep(Duration::from_secs(2)).await;
        }
    });

    // no relation function
    let total = add(1, 2);
    println!("{}", total);

    let s = sum(16, vec![2, 4, 7, 3, 12]);
    println!("{}", s);

    let ss = sum(10, vec![7, 3, 4, 3, 2]);
    println!("{}", ss);

    life_time();

    // メインループでチャネルからのメッセージを受け取る
    while let Some(msg) = rx.recv().await {
        match msg {
            Message::Text(text) => {
                println!("Received string: {}", text);
            }
            Message::Number(num) => {
                println!("Received int: {}", num);
            }
        }
    }
}

fn life_time() {
    let n: i32;
    {
        let a = 5;
        n = a;
    }
    println!("{}", n);
}

