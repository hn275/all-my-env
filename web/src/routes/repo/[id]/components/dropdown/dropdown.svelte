<script lang="ts">
	import { store, type RepositoryEnv } from "../../store";

	export let repoName: string;
	let state: RepositoryEnv;
	$: state = $store;
</script>

<div class="dropdown dropdown-hover dropdown-end">
	<div class="btn btn-ghost bg-dark-100 hover:cursor-auto">
		<i class="fa-solid fa-gear"></i>
	</div>
	<ul
		class="dropdown-content menu bg-dark-100 rounded-box z-[1] w-52 p-2 shadow"
	>
		<li>
			<button>
				<i class="fa-solid fa-user w-5"></i>
				Contributors
			</button>
		</li>
		<li>
			<a
				href={`https://github.com/${repoName}`}
				target="_blank"
				class=""
			>
				<i class="fa-brands fa-github w-5"></i>
				Git Repository
			</a>
		</li>
	</ul>
</div>
