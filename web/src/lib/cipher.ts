import aes from "crypto-js/aes";
import CryptoJS from "crypto-js";

// generated from https://generate-random.org/string-generator?count=1&length=32&has_lowercase=0&has_lowercase=1&has_uppercase=0&has_uppercase=1&has_symbols=0&has_numbers=0&has_numbers=1&is_pronounceable=0
// key should be of 32 char
const test_secret = "hSchE34mgOHmvsQokZGSZY4jPQbqD9qY";

export enum SecretKeys {
	Auth = "auth",
}

export class Cipher {
	private _key: string = "";
	private value: string = "";

	constructor(value: string) {
		this.value = value;
	}

	public setKey(k: SecretKeys) {
		this._key = this.getKey(k);
		return this;
	}

	public encrypt(): string {
		if (this._key === "") throw new Error("secret key not set");
		return aes.encrypt(this.value, this._key).toString();
	}

	public decrypt(): string {
		if (this._key === "") throw new Error("secret key not set");
		return aes.decrypt(this.value, this._key).toString(CryptoJS.enc.Utf8);
	}

	private getKey(key: SecretKeys): string {
		switch (key) {
			case SecretKeys.Auth:
				return test_secret; // TODO: load from env instead
			default:
				return "";
		}
	}
}
