import React, { ActionDispatch } from "react";
import { createContext, useReducer, useEffect } from "react";

export const AuthContext = createContext({} as IAuthContextState);

export const authReducer = (state: IAuthContextState, action: IAction) => {
  switch (action.type) {
    case "LOGIN":
      return { token: action.payload }
    case "LOGOUT":
      return { token: null }
    default:
      return state
  }
}

export const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = useReducer(authReducer, {
    token: null
  });

  useEffect(() => {
    const token = JSON.parse(localStorage.getItem("token"));

    if (token) {
      dispatch({ type: "LOGIN", payload: token })
    }
  }, [])

  console.log("AuthContext state: ", state);

  return (
    <AuthContext.Provider value={{ ...state, dispatch }}>
      {children}
    </AuthContext.Provider>
  )
}

interface IAction {
  type: string;
  payload: any;
}

interface IAuthContextState {
  token: string;
  dispatch: ActionDispatch<[action: IAction]>;
}
