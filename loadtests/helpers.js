import http from 'k6/http';
import { check } from 'k6';

import { BASE_URL, authParams, getRandomBook, sleepRandom } from './utils.js';

export function listBooks(authToken, minSleep = 0, maxSleep = 0) {
	const params = authParams(authToken);

	const res = http.get(`${BASE_URL}/books?limit=2000`, params);
	const success = check(res, {
		'Book GET status is 200': (r) => r.status === 200,
		'Book GET response is not null': (r) => r.json() !== null,
		'Book GET response has more than zero entries': (r) => r.json().length > 0 ,
		'Book GET title is not empty': (r) => r.status === 200 && r.json()[0].title !== '',
		'Book GET author is not zero': (r) => r.status === 200 && r.json()[0].author !== 0,
	});

	if (!success) {
		console.error(`Failed to GET books ${res.status} ${res.json().error}`);
		return false;
	}

	sleepRandom(minSleep, maxSleep);
	return true;
}

export function postBook(authToken, minSleep = 0, maxSleep = 0) {
	const params = authParams(authToken);
	params.headers['Content-Type'] = 'application/json';
	const payload = getRandomBook();

	const res = http.post(`${BASE_URL}/books`, payload, params);
	const success = check(res, {
		'Book POST status is 201': (r) => r.status === 201,
		'Book POST body is not null': (r) => r.json() !== null,
		'Book POST Location header is not empty': (r) => r.headers.Location !== null
	});

	if (!success) {
		console.error(`Failed to POST book ${res.status} ${res.body}`);
		return { isOk: false, id: 0 };
	}

	const headerLocSplit = res.headers.Location.split('/');
	const locID = headerLocSplit[headerLocSplit.length - 1];

	sleepRandom(minSleep, maxSleep);
	return {
		isOk: true,
		id: locID
	};
}

export function patchRandomBook(authToken, newestBookId, minSleep = 0, maxSleep = 0) {
	let patchID = Math.floor(Math.random() * newestBookId + 1);
	if (patchID > newestBookId) {
		patchID = newestBookId;
	}

	const params = authParams(authToken);
	params.headers['Content-Type'] = 'application/json';
	const payload = getRandomBook();

	const res = http.patch(http.url`${BASE_URL}/books/${patchID}`, payload, params);
	const success = check(res, {
		'Book PATCH status is 204': (r) => r.status === 204,
	});

	if (!success) {
		console.error(`Failed to PATCH book ${res.status} ${res.body}`);
		return false;
	}

	sleepRandom(minSleep, maxSleep);
	return true;
}

export function listRandomResource(authToken, minSleep = 0, maxSleep = 0) {
	const params = authParams(authToken);

	const resources = ['books', 'genres', 'authors', 'languages'];
	const randomResource = resources[Math.floor(Math.random() * resources.length)];

	const res = http.get(`${BASE_URL}/${randomResource}?limit=2000`, params);
	const success = check(res, {
		'Resource GET status is 200': (r) => r.status === 200,
		'Resource response is not null': (r) => r.json() !== null,
		'Resource response has more than zero entries': (r) => r.json().length > 0 ,
	});

	if (!success) {
		console.error(`Failed to GET resource ${res.status} ${res.json().error}`);
		return false;
	}

	sleepRandom(minSleep, maxSleep);
	return true;
}
