# Server

## Doc

- [Github endpoints](https://docs.github.com/en/rest/overview/endpoints-available-for-github-apps?apiVersion=2022-11-28)
- [Firebase](https://docs.rs/firebase-rs/latest/firebase_rs/struct.Firebase.html)
- [Reqwest](https://docs.rs/reqwest/latest/reqwest/)
- [Actix Web](https://docs.rs/crate/actix-web/latest)
- [Aes GCM](https://docs.rs/aes-gcm/latest/aes_gcm/)

## Docker

```sh
docker compose up # start up a server
                  # `--build` flag is needed first time running
```

## .env

> yeah, pain, maybe we can be the users of our own product soon.

```txt
FIREBASE_DB=""
GITHUB_CLIENT_ID=""
GITHUB_CLIENT_SECRET=""
```

## Examples

<details>
    <summary>Ciphering/Deciphering</summary>

```rust
#[derive(Clone, Debug)]
struct Foo {
    id: String,
    key: String,
    value: String,
}

impl KeyGen for Foo {
    // returns part of the key to use in key derivation function
    // it is the only method required.
    fn key(&self) -> Vec<u8> {
        [self.id.as_bytes(), self.key.as_bytes()].concat()
    }
}

[#test]
fn cipher_decipher() {
    let foo = Foo {
        id: "123".to_owned(),
        key: "foo".to_owned(),
        value: "bar".to_owned(),
    };

    // cipher
    let nonce = None;
    // `None` value here would be the nonce
    // seal function will generate a new nonce.
    // if nonce is available to reuse _for the same secret_
    // (ie, PUT/PATCH endpoint to update the variable)
    // pass it in as `Some([u8; 12])` or `Some(cipher::Nonce)`
    let ciphered = Key::new(KeyType::RowKey)
        .generate_key(&foo)
        .unwrap()
        .seal(foo.value.as_bytes(), nonce)
        .unwrap();

    assert_ne!(foo.value.as_bytes(), ciphered.ciphertext);

    // decipher
    let opened = Key::new(KeyType::RowKey)
        .generate_key(&foo)
        .unwrap()
        .open(ciphered.ciphertext.as_slice(), ciphered.nonce)
        .unwrap();

    assert_eq!(foo.value.as_bytes(), opened);
}
```

</details>
