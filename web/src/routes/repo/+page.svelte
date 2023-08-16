<script lang="ts">
  import Main from "@components/main.svelte";
  import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";

  // breadcrumbs
  let breadcrumbs: Array<Breadcrumbs> | undefined;
  onMount(() => {
    const url: URL = new URL(window.location.href);
    let repoName: string = url.searchParams.get("name")!;
    repoName = decodeURIComponent(repoName);
    breadcrumbs = [{ text: repoName, href: url.toString() }];
  })

</script>

<Main {breadcrumbs}>

</Main>
