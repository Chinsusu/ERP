import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 20 },  // Ramp up
        { duration: '1m', target: 50 },   // Stay at 50
        { duration: '1m', target: 100 },  // Peak load
        { duration: '30s', target: 0 },   // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<200', 'p(99)<500'],
        http_req_failed: ['rate<0.01'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';

export default function () {
    const payload = JSON.stringify({
        email: 'admin@company.vn',
        password: 'Admin@123',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(`${BASE_URL}/api/v1/auth/login`, payload, params);

    check(res, {
        'is status 200': (r) => r.status === 200,
        'has token': (r) => r.json().data.access_token !== undefined,
    });

    sleep(1);
}
