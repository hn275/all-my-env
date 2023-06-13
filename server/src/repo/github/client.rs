use reqwest;

pub struct GithubClient {
    queries: Vec<String>,
    token: String,
}

impl GithubClient {
    pub fn new_with_token(token: &str) -> GithubClient {
        return GithubClient {
            queries: Vec::<String>::new(),
            token: token.to_owned(),
        };
    }

    pub fn query<'a>(&mut self, key: &'a str, value: &'a str) -> &mut Self {
        let q = format!("{}={}", key, value);
        self.queries.push(q);
        return self;
    }

    pub fn get<'a>(&self, path: &'a str) -> reqwest::RequestBuilder {
        let mut url = format!("https://api.github.com{}", path);

        if self.queries.len() > 0 {
            let queries = self.queries.join("&");
            url.push_str("?");
            url.push_str(queries.as_str());
        }

        return reqwest::Client::new()
            .get(url)
            .header("User-Agent", "All My ENV")
            .header("Authorization", format!("Bearer {}", self.token))
            .header("Accept", "application/vnd.github+json");
    }
}
