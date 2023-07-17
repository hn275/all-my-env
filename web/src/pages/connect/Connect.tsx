import { WEB } from 'lib/routes';
import { useSearchParams, useNavigate } from 'react-router-dom'

export function Connect() {
  const nav = useNavigate();
  const [param] = useSearchParams()
  const repoID = param.get("repo_id")

  if (!repoID) nav(WEB.home)




  return <>Connect</>
}
