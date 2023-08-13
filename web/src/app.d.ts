// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	declare namespace EnvHub {
		interface Error {
			message: string;
		}
		type Response<T> = T | EnvHub.Error;
	}
}

export {};
