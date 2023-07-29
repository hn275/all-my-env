import { AiFillGithub } from "react-icons/ai";
import { Repository } from './types'
import { Link } from "react-router-dom";
import { WEB } from "lib/routes";

interface Props {
  repo: Repository
  isOwner: boolean
}

export function Repo({ repo, isOwner }: Props) {
  const params: string = new URLSearchParams({
    "repo_id": String(repo.id),
    "repo_name": encodeURIComponent(repo.full_name),
  }).toString()

  return (
    <>
      <div>
        <img
          src={repo.owner.avatar_url}
          alt={repo.owner.login}
          className="aspect-auto w-8 rounded-full"
        />

        <Link to={`/repos/${repo.id}`}>{repo.name}</Link>
        {
          !isOwner ?
            (
              <p>Collaborator</p>
            ) : repo.linked ? (
              <p>Linked</p>
            ) : (
              <Link
                to={`${WEB.repo}/connect?${params}`}
                className="block bg-accent-blue w-max p-1"
              >
                Link Repo
              </Link>
            )
        }
      </div>

      <a href={`https://www.github.com/${repo.full_name}`} target="_blank">
        <AiFillGithub />
      </a>
    </>
  )
}
