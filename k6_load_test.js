import http from "k6/http";
import { check, sleep } from "k6";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

// --- Configuration ---
// Add the base URLs for your three running KGS instances.
const baseUrls = [
  "http://localhost:8080",
  "http://localhost:8081",
  "http://localhost:8082",
];

// --- Test Scenario ---
// This scenario ramps up the number of virtual users (VUs) over time
// to simulate a realistic increase in traffic.
export const options = {
  stages: [
    { duration: "30s", target: 50 }, // Ramp up to 50 users over 30 seconds
    { duration: "30s", target: 1000 }, // Ramp up to 100 users over the next minute
    { duration: "30s", target: 2000 }, // Ramp up to 200 users over the next 2 minutes¬
    { duration: "30s", target: 0 }, // Ramp down to 0 users
  ],
  thresholds: {
    http_req_failed: ["rate<0.01"], // Fail the test if more than 1% of requests fail
    http_req_duration: ["p(95)<200"], // 95% of requests must complete under 200ms
  },
};

// --- Main Test Logic ---
export default function () {
  // Distribute requests evenly across the server instances using the virtual user ID.
  const instanceIndex = __VU % baseUrls.length;
  const baseUrl = baseUrls[instanceIndex];
  const url = `${baseUrl}/get-key`;

  // Make the HTTP GET request to fetch a key.
  const res = http.get(url);

  // Check if the request was successful and the key is not empty.
  const isSuccess = check(res, {
    "status is 200 OK": (r) => r.status === 200,
    "key is not empty": (r) => r.body && r.body.length > 0,
  });

  // If the request was successful, log the key with a special prefix for parsing
  if (isSuccess) {
    console.log(`KEY:${res.body}`);
  }

  // Add a small sleep to simulate a real user pausing between actions.
  sleep(0.1);
}

// --- Custom Summary Handler ---
export function handleSummary(data) {
  return {
    'stdout': textSummary(data),
  };
}
