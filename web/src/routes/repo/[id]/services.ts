import { store } from "./store";

function deleteVariable(varID: string): void {
	store.update((v) => ({ ...v, state: { deleteVariable: varID } }));
}
