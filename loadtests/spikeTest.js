import * as helpers from './helpers.js';
import { getToken } from './utils.js';

export const options = {
	stages: [
		{ duration: '5s', target: 1000 },
		{ duration: '5s', target: 1000 },
		{ duration: '20s', target: 20000 },
		{ duration: '2m', target: 20000 },
		{ duration: '5s', target: 0 },
	],
	thresholds: {
		http_req_failed: ['rate<0.15']
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
	let sleepMaxTime = 0.5;
	if (!helpers.listRandomResource(authToken, sleepMinTime, sleepMaxTime)) {
		console.warn('Skipping - list random resource failed');
		return;
	}

	sleepMinTime = 0;
	sleepMaxTime = 1;
	if (!helpers.listBooks(authToken, sleepMinTime, sleepMaxTime)) {
		console.warn('Skipping - list books failed');
		return;
	}

	sleepMinTime = 0;
	sleepMaxTime = 2;
	const {isOk} = helpers.postBook(authToken, sleepMinTime, sleepMaxTime);
	if (!isOk) {
		console.warn('Skipping - POST book failed');
		return;
	}
}
