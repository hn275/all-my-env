import { AiFillGithub } from "react-icons/ai";
import { Repository } from "./types";
import { Link } from "react-router-dom";
import { WEB } from "lib/routes";

export function Repo(repo: Repository) {
  const params: string = new URLSearchParams({
    repo_name: encodeURIComponent(repo.full_name),
  }).toString();

  return (
    <li className="bg-dark-100">
      <div>
        <img
          src={repo.owner.avatar_url}
          alt={repo.owner.login}
          className="aspect-auto w-8 rounded-full"
        />

        <Link to={`${WEB.repo}/${repo.id}`}>{repo.name}</Link>
        {!repo.is_owner ? (
          <p>Collaborator</p>
        ) : repo.linked ? (
          <p>Linked</p>
        ) : (
          <Link
            to={`${WEB.repo}/${repo.id}/connect?${params}`}
            className="bg-accent-blue block w-max p-1"
          >
            Link Repo
          </Link>
        )}
        {repo.fork ? <>forked</> : <></>}
      </div>

      <a href={`https://www.github.com/${repo.full_name}`} target="_blank">
        <AiFillGithub />
      </a>
    </li>
  );
}
