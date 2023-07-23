package jsonwebtoken

type jwtDecoderMock struct{}

func Mock() {
	Decoder = &jwtDecoderMock{}
}

func (j *jwtDecoderMock) Decode(t string) (*JwtToken, error) {
	u := GithubUser{Token: t, ID: 1}
	return &JwtToken{GithubUser: u}, nil
}
