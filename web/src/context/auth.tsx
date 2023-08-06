import { Dispatch, ReactNode, useContext, useReducer } from "react";
import { createContext } from "react";

// state
type UserLogin = {
	id: number;
	github_token: string;
	access_token: string;
	avatar_url: string;
	name: string;
	login: string;
};
type User = UserLogin | undefined;

// dispatch
type AuthCtxLogIn = {
	type: "login";
	payload: User;
};

type AuthCtxLogOut = {
	type: "logout";
};

type AuthCtxDispatch = AuthCtxLogIn | AuthCtxLogOut;

interface AuthCtxProps {
	state: User;
	dispatch: Dispatch<AuthCtxDispatch>;
}

const AuthCtx = createContext<AuthCtxProps>({} as AuthCtxProps);

function reducer(state: User, action: AuthCtxDispatch) {
	switch (action.type) {
		case "login":
			return action.payload;
		case "logout":
			return undefined;
		default:
			return state;
	}
}

type AuthProviderProps = {
	children: ReactNode;
};

export function AuthProvider({ children }: AuthProviderProps) {
	const [state, dispatch] = useReducer(reducer, undefined);
	return (
		<AuthCtx.Provider value={{ state, dispatch }}>{children}</AuthCtx.Provider>
	);
}

export function useAuth(): AuthCtxProps {
	return useContext(AuthCtx);
}
