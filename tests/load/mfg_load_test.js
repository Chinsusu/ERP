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

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8083';

export default function () {
    const params = {
        headers: {
            'Authorization': 'Bearer sample-token',
        },
    };

    // Replace with a valid BOM ID if available
    const bomId = '32495dc-1234-5678-90ab-cdef12345678';
    const res = http.get(`${BASE_URL}/api/v1/boms/${bomId}`, params);

    check(res, {
        'is status 200 or 404': (r) => r.status === 200 || r.status === 404,
    });

    sleep(1);
}
