import http from 'k6/http';
import { check, sleep } from 'k6';

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

const BASE_URL = 'http://localhost:8080/api/v1';

export function setup() {
	const token = getToken();
	if (!token) {
		throw new Error('Authorization token is null');
	}

	return token;
}

export default function(authToken) {
	if (!listBooks(authToken)) {
		return;
	}

	const {isOk, id} = postBook(authToken);
	if (!isOk) {
		return;
	}

	if (!patchBook(authToken, id)) {
		return;
	}

	if (!listRandomResource(authToken)) {
		return;
	}
}

function listBooks(authToken) {
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

	sleepRandom(0, 2);
	return true;
}

function postBook(authToken) {
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

	sleepRandom(2, 3);
	return {
		isOk: true,
		id: locID
	};
}

function patchBook(authToken, newestBookId) {
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

	sleepRandom(3, 5);
	return true;
}

function listRandomResource(authToken) {
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

	sleepRandom(2, 4);
	return true;
}

function getRandomBook() {
	const randomYear     = Math.floor(Math.random() * 1525 + 500);
	const randomPages    = Math.floor(Math.random() * 1820 + 20);
	const randomAuthor   = Math.floor(Math.random() * 9 + 1);
	const randomGenre    = Math.floor(Math.random() * 9 + 1);
	const randomLanguage = Math.floor(Math.random() * 11 + 1);

	return JSON.stringify({
		title: 'Book title',
		year: randomYear,
		pages: randomPages,
		author: randomAuthor,
		genre: randomGenre,
		language: randomLanguage
	});
}

function sleepRandom(min = 0, max = 0) {
	sleep(Math.random() * (max - min) + min);
}

function authParams(token) {
	return {
		headers: {
			'Authorization': `Bearer ${token}`,
		},
	};
}

function getToken() {
	const url = `${BASE_URL}/login`;
	const payload = JSON.stringify({
		return_admin_token: true,
	});
	const headers = {
		'Content-Type': 'application/json',
	};

	const res = http.post(url, payload, headers);
	const success = check(res, {
		'Token response status is 200': (r) => r.status === 200,
		'Token is not null': (r) => r.json().token !== null
	});

	if (!success) {
		console.error(`Failed to retrieve token ${res.status} ${res.body}`);
		return null;
	}

	return res.json().token;
}
