import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'],
        http_req_failed: ['rate<0.01'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8082';

export default function () {
    // Mocking auth token for now
    const params = {
        headers: {
            'Authorization': 'Bearer sample-token',
        },
    };

    const res = http.get(`${BASE_URL}/api/v1/stock`, params);

    check(res, {
        'is status 200': (r) => r.status === 200,
    });

    sleep(1);
}
