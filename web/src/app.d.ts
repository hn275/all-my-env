// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface Platform {}
		interface ApiResponseError {
			message: string;
		}
		type ApiResponse<T> = T | ApiResponse;
	}
}

export {};
