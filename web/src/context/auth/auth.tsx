import { Dispatch, ReactNode, useEffect, useReducer } from "react";
import { createContext } from "react";
import { Entry } from "lib/storage";

// state
export type User = {
	id: number;
	access_token: string;
	avatar_url: string;
	name: string;
	login: string;
};

// dispatch
export type AuthCtxLogIn = {
	type: "login";
	payload: User;
};

export type AuthCtxLogOut = {
	type: "logout";
};

export type AuthCtxDispatch = AuthCtxLogIn | AuthCtxLogOut;

export type AuthCtxState = {
	user?: User;
	auth: boolean;
};

export interface AuthCtxProps {
	state: AuthCtxState;
	dispatch: Dispatch<AuthCtxDispatch>;
}

export const AuthCtx = createContext<AuthCtxProps>({} as AuthCtxProps);

const initAuth = false;

function reducer(state: AuthCtxState, action: AuthCtxDispatch) {
	switch (action.type) {
		case "login":
			window.localStorage.setItem(Entry.user, JSON.stringify(action.payload));
			return { auth: true, user: action.payload };
		case "logout":
			window.localStorage.removeItem(Entry.user);
			return { auth: initAuth };
		default:
			return state;
	}
}

type AuthProviderProps = {
	children: ReactNode;
};

export function AuthProvider({ children }: AuthProviderProps) {
	const init: AuthCtxState = {
		auth: initAuth,
	};
	const [state, dispatch] = useReducer(reducer, init);

	useEffect(() => {
		const u = window.localStorage.getItem(Entry.user);
		if (!u) return;
		const user: User = JSON.parse(u);
		dispatch({ type: "login", payload: user });
	}, []);

	return (
		<AuthCtx.Provider value={{ state, dispatch }}>{children}</AuthCtx.Provider>
	);
}
