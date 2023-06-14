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
