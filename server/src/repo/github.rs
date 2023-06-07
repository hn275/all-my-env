pub struct AuthClient<'a> {
    pub client_id: &'a str,
    pub client_secret: &'a str,
    pub code: &'a str,
    pub redirect_uri: &'a str,
    pub scope: Vec<&'a str>,
}

pub struct GithubClient;

impl<'a> GithubClient {
    fn get_access_token(auth: &'a AuthClient) {}
}
