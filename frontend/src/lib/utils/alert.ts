import Swal from 'sweetalert2';

import 'sweetalert2/dist/sweetalert2.min.css';

export async function confirmDelete(
	title = 'Hapus data?',
	message = 'Data yang sudah dihapus tidak dapat dikembalikan.'
): Promise<boolean> {
	const result = await Swal.fire({
		title,
		text: message,
		icon: 'warning',
		showCancelButton: true,
		confirmButtonText: 'Hapus',
		cancelButtonText: 'Batal',
		reverseButtons: true,
		focusCancel: true,
		buttonsStyling: false,
		customClass: {
			popup: 'gowatch-alert',
			title: 'gowatch-alert-title',
			htmlContainer: 'gowatch-alert-text',
			actions: 'gowatch-alert-actions',
			confirmButton: 'gowatch-alert-danger',
			cancelButton: 'gowatch-alert-cancel'
		}
	});

	return result.isConfirmed;
}

export async function showSuccess(title: string, message: string): Promise<void> {
	await Swal.fire({
		title,
		text: message,
		icon: 'success',
		confirmButtonText: 'OK',
		buttonsStyling: false,
		customClass: {
			popup: 'gowatch-alert',
			title: 'gowatch-alert-title',
			htmlContainer: 'gowatch-alert-text',
			confirmButton: 'gowatch-alert-primary'
		}
	});
}

export async function showError(title: string, message: string): Promise<void> {
	await Swal.fire({
		title,
		text: message,
		icon: 'error',
		confirmButtonText: 'Tutup',
		buttonsStyling: false,
		customClass: {
			popup: 'gowatch-alert',
			title: 'gowatch-alert-title',
			htmlContainer: 'gowatch-alert-text',
			confirmButton: 'gowatch-alert-primary'
		}
	});
}

export async function showInfo(title: string, message: string): Promise<void> {
	await Swal.fire({
		title,
		text: message,
		icon: 'info',
		confirmButtonText: 'Mengerti',
		buttonsStyling: false,
		customClass: {
			popup: 'gowatch-alert',
			title: 'gowatch-alert-title',
			htmlContainer: 'gowatch-alert-text',
			confirmButton: 'gowatch-alert-primary'
		}
	});
}
