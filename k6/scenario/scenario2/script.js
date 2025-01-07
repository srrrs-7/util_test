import { browser } from 'k6/browser';
import http from 'k6/http';
import { sleep, check } from 'k6';
import { jUnit, textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';
import papaparse from 'https://jslib.k6.io/papaparse/5.1.1/index.js';
import { SharedArray } from 'k6/data';

// 環境変数からURLを取得、設定されていない場合はデフォルト値を設定
const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:80/employee';
const API_BASE_URL = __ENV.API_BASE_URL || 'http://host.docker.internal:80/api';

// CSVファイルからデータを読み込む
const csvData = new SharedArray('users', function () {
  const f = papaparse.parse(open('./var.csv'), { header: true }).data;
  return f;
});

export const options = {
  scenarios: {
    ui_chromium: {
      executor: 'ramping-vus',
      exec: 'uiTestChromium',
      options: {
        browser: {
          type: 'chromium',
          headless: false,
          slowMo: '50ms',
        },
      },
      stages: [
        { duration: '30s', target: 10 },
        { duration: '1m', target: 20 },
        { duration: '30s', target: 0 },
      ],
      gracefulStop: '10s',
    },
    ui_firefox: {
      executor: 'shared-iterations',
      exec: 'uiTestFirefox',
      options: {
        browser: {
          type: 'firefox',
          headless: true,
        },
      },
      vus: 5,
      iterations: 20,
      maxDuration: '30s',
    },
    api_load_test: {
      executor: 'constant-vus',
      exec: 'apiTest',
      vus: 50,
      duration: '30s',
    },
  },
  thresholds: {
    'browser_web_vital_lcp{browser:chromium}': ['p(95) < 1500'],
    'browser_web_vital_fid{browser:chromium}': ['p(90) < 100'],
    'browser_web_vital_cls{browser:chromium}': ['avg < 0.1'],
    'browser_web_vital_lcp{browser:firefox}': ['p(90) < 2000'],
    'http_req_duration{scenario:api_load_test}': ['p(99) < 500'],
    'http_req_failed{scenario:api_load_test}': ['rate < 0.01'],
    'checks{api_check:true}': ['rate>0.99'],
  },
  http: { // HTTP設定を共通化
    userAgent: 'MyK6Test/1.0',
    insecureSkipTLSVerify: true,
  },
  systemTags: ['iter', 'metric', 'scenario', 'status', 'subproto', 'tls_version', 'url', 'error', 'expected_response', 'group', 'method', 'name', 'proto', 'vu'],
};

// UIテスト (Chromium)
export async function uiTestChromium() {
  const page = browser.newPage();
  await runUITest(page, 'chromium');
}

// UIテスト (Firefox)
export async function uiTestFirefox() {
  const page = browser.newPage();
  await runUITest(page, 'firefox');
}

// 共通UIテストロジック
async function runUITest(page, browserType) {
  // CSVデータからランダムにユーザー情報を取得
  const userData = csvData[Math.floor(Math.random() * csvData.length)];

  try {
    await page.goto(BASE_URL);
    await page.screenshot({ path: `screenshots/screenshot-${browserType}.png` });
    
    // ログを標準出力に出力
    console.log(`[${browserType}] - Page loaded successfully - User: ${userData.username}`);

    // 例：ユーザー名とパスワードを入力 (実際のアプリケーションに合わせて変更)
    // await page.type('#username', userData.username);
    // await page.type('#password', userData.password);
    // await page.click('#login');

  } catch (error) {
    console.error(`[${browserType}] - Error during UI test: ${error}`);
  } finally {
    page.close();
  }
}

// APIテスト
export function apiTest() {
  // CSVデータからランダムにユーザー情報を取得
  const userData = csvData[Math.floor(Math.random() * csvData.length)];

  // 例：ユーザー名を使ってAPIを叩く (実際のアプリケーションに合わせて変更)
  const res = http.get(`${API_BASE_URL}/employees?username=${userData.username}`);

  check(res, {
    'is status 200': (r) => r.status === 200,
  },{
    api_check: true,
  });

  // ログを標準出力に出力
  console.log(`[API] - Status: ${res.status}, URL: ${res.url}, User: ${userData.username}`);

  sleep(1);
}

// テスト結果のハンドリング
export function handleSummary(data) {
  // テスト結果のログファイル名を生成
  const logFileName = `log/summary-${new Date().toISOString()}.log`;
  
  // console.logの出力をファイルにリダイレクトするための設定
  const stdout = console.log;
  console.log = function (...args) {
    const logLine = args.map(arg => (typeof arg === 'object' ? JSON.stringify(arg) : arg)).join(' ');
    stdout(logLine); // 元のconsole.logの挙動を維持
    http.post('file:///' + logFileName, logLine + '\n'); // ログをファイルに追記
  };
  
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }), // 標準出力に概要を出力
    'log/summary.log': textSummary(data, { indent: ' ', enableColors: false }), // /k6/log/summary.logにテキスト形式で結果を出力
    'log/junit.xml': jUnit(data), // /k6/log/junit.xmlにJUnit形式で結果を出力
    [logFileName]: http.file(console.log, 'text/plain'), // console.logの内容をファイルに書き込み
  };
}