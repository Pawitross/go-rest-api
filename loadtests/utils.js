import http from 'k6/http';
import { check, sleep } from 'k6';

export const BASE_URL = 'http://localhost:8080/api/v1';

export function getRandomBook() {
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

export function sleepRandom(min = 0, max = 0) {
	if (min === 0 && max === 0) {
		return;
	}

	if (min < 0 || max < 0) {
		throw new Error('Value/s should not be negative');
	}

	if (min > max) {
		throw new Error('Minimal sleep value is bigger than maximum');
	}

	sleep(Math.random() * (max - min) + min);
}

export function authParams(token) {
	return {
		headers: {
			'Authorization': `Bearer ${token}`,
		},
	};
}

export function getToken() {
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
