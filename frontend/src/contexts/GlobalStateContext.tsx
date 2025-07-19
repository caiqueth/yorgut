import { Backdrop, CircularProgress } from "@mui/material";
import React, { createContext, useReducer } from "react";

export const GlobalStateContext = createContext({} as IGlobalContextState);

const globalStateReducer = (state: IGlobalContextState, action: string) => {
  switch (action) {
    case "SET_LOADING":
      return { ...state, isLoading: true };
    case "UNSET_LOADING":
      return { ...state, isLoading: false };
    default:
      return state;
  }
};

export const GlobalContextStateProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = useReducer(globalStateReducer, {
    isLoading: false,
  } as IGlobalContextState);

  return (
    <GlobalStateContext.Provider value={{ ...state, dispatch }}>
      {children}
      <Backdrop open={state.isLoading} style={{ zIndex: 9999 }}>
        <CircularProgress />
      </Backdrop>
    </GlobalStateContext.Provider>
  );
};

interface IGlobalContextState {
  isLoading: boolean;
  dispatch: React.Dispatch<string>;
}
