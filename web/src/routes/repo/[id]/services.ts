import { store, type Variable } from "./store";

// trigger the delete modal to be open
export function deleteVariable(v: Variable): void {
	store.update((s) => ({ ...s, state: { deleteVariable: v } }));
}

export function cancelDelete(): void {
	store.update((s) => ({ ...s, state: { deleteVariable: undefined } }));
}
