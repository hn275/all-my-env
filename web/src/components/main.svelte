<script lang="ts">
	import { onMount } from "svelte";
	import { authStore } from "@lib/auth";
	import "../index.css";
	import { UserStorage } from "@lib/storage";

	onMount(() => {
		authStore.set({ ...authStore.get(), user: UserStorage.get() });
		authStore.listen((store) => {
			if (store.user) {
				UserStorage.set(store.user);
			} else {
				UserStorage.remove();
			}
		});
	});
</script>

<div class="max-w-screen-xl w-full mx-auto">
	<slot />
</div>
