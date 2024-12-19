import http from 'k6/http';
 import { sleep, check } from 'k6';
 
 // each env url
 const URL = 'http://host.docker.internal:80/employee';
 
 export const options = {
     vus: 20,
     duration: '600s',
 };
 
 export default function () {
     const params = {
         headers: {
             'Cookie': 'cookie'
         },
     }
     const res = http.get(URL, params);
     check(res, {
         'status is 200': (r) => r.status === 200,
     });
 
     sleep(1);
 }