import * as helpers from './helpers.js';
import { getToken } from './utils.js';

export const options = {
	stages: [
		{ duration: '10s', target: 1000 },
		{ duration: '10s', target: 1000 },
		{ duration: '20s', target: 13000 },
		{ duration: '10m', target: 13000 },
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
	if (!helpers.listBooks(authToken, 0, 1)) {
		console.warn('Skipping - list books fail');
		return;
	}

	const {isOk, id} = helpers.postBook(authToken, 0, 2);
	if (!isOk) {
		console.warn('Skipping - POST book fail');
		return;
	}

	if (!helpers.patchBook(authToken, id, 3, 6)) {
		console.warn('Skipping - PATCH book fail');
		return;
	}

	if (!helpers.listRandomResource(authToken, 0, 0.5)) {
		console.warn('Skipping - list random resource fail');
		return;
	}
}
