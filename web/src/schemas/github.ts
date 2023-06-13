export interface GithubAccount {
  login: string;
  id: number;
  avatarUrl: string;
  name: string;
  email: string;
  bio: string;
  publicRepos: number;
  followers: number;
  following: number;
  createdAt: Date;
  updatedAt: Date;
}
