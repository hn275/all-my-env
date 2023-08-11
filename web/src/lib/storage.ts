import type { User } from "./auth";

export class UserStorage {
	private static entry: string = "envhub:user";
	public static set(u: User): void {
		window.localStorage.setItem(this.entry, JSON.stringify(u));
	}

	public static get(): User | undefined {
		const u = window.localStorage.getItem(this.entry);
		if (!u) return;
		return JSON.parse(u) as User;
	}

	public static remove(): void {
		window.localStorage.removeItem(this.entry);
	}
}

export class TokenStorage {
	private static entry: string = "envhub:refreshed";

	public static refreshed(): boolean {
		const r = window.sessionStorage.getItem(this.entry);
		return r !== null;
	}

	public static setRefreshed(): void {
		window.sessionStorage.setItem(this.entry, JSON.stringify("true"));
	}
}
