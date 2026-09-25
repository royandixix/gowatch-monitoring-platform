export function formatDate(value: string): string {
	const date = new Date(value);

	if (Number.isNaN(date.getTime())) {
		return '-';
	}

	return new Intl.DateTimeFormat('id-ID', {
		dateStyle: 'medium',
		timeStyle: 'short'
	}).format(date);
}

export function formatInterval(seconds: number): string {
	if (seconds < 60) {
		return `${seconds} detik`;
	}

	if (seconds % 60 === 0) {
		return `${seconds / 60} menit`;
	}

	return `${seconds} detik`;
}
