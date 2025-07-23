import * as helpers from './helpers.js';
import { getToken } from './utils.js';

export const options = {
	stages: [
		{ duration: '10s', target: 1000 },
		{ duration: '10m', target: 1000 },
		{ duration: '10s', target: 0 },
	],
	thresholds: {
		http_req_failed: ['rate<0.01']
	}
};

export function setup() {
	const token = getToken();
	if (!token) {
		throw new Error('Authorization token is null');
	}

	return token;
}

export default function(authToken) {
	let sleepMinTime = 0;
	let sleepMaxTime = 2;
	if (!helpers.listBooks(authToken, sleepMinTime, sleepMaxTime)) {
		console.warn('Skipping - list books failed');
		return;
	}

	sleepMinTime = 2;
	sleepMaxTime = 3;
	const {isOk, id} = helpers.postBook(authToken, sleepMinTime, sleepMaxTime);
	if (!isOk) {
		console.warn('Skipping - POST book failed');
		return;
	}

	sleepMinTime = 3;
	sleepMaxTime = 5;
	if (!helpers.patchRandomBook(authToken, id, sleepMinTime, sleepMaxTime)) {
		console.warn('Skipping - PATCH book failed');
		return;
	}

	sleepMinTime = 2;
	sleepMaxTime = 4;
	if (!helpers.listRandomResource(authToken, sleepMinTime, sleepMaxTime)) {
		console.warn('Skipping - list random resource failed');
		return;
	}
}
