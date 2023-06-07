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

  public async encrypt() {
    if (this._key === "") throw new Error("secret key not set");
    const encoder = new TextEncoder();

    const iv = crypto.getRandomValues(new Uint8Array(12));
    const encoded = encoder.encode(this.value);
    const key = await crypto.subtle.importKey(
      "raw",
      encoder.encode(this._key),
      { name: "AES-GCM" },
      false,
      ["encrypt"],
    );
    const cipherText = await crypto.subtle.encrypt(
      { name: "AES-GCM", iv },
      key,
      encoded,
    );

    const decoder = new TextDecoder();
    return {
      encrypted: decoder.decode(cipherText),
      iv,
    };

    //return aes.encrypt(this.value, this._key).toString();
  }

  public async decrypt(iv: Uint8Array) {
    if (this._key === "") throw new Error("secret key not set");

    const encoder = new TextEncoder();
    const key = await crypto.subtle.importKey(
      "raw",
      encoder.encode(this._key),
      "AES-GCM",
      false,
      ["encrypt"],
    );

    try {
      const pt = await crypto.subtle.decrypt(
        { name: "AES-GCM", iv },
        key,
        encoder.encode(this.value),
      );

      const decoder = new TextDecoder();
      return decoder.decode(pt);
      //return aes.decrypt(this.value, this._key).toString(CryptoJS.enc.Utf8);
    } catch (e) {
      console.error(e);
    }
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
