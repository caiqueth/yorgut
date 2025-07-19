import { Backdrop, CircularProgress } from "@mui/material";
import { useMe } from "api/me/queries";
import React from "react";
import { createContext } from "react";

export const AuthContext = createContext({} as IAuthContextState);

export const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const { isFetching } = useMe();

  return (
    <AuthContext.Provider value={{}}>
      <Backdrop open={isFetching}>
        <CircularProgress />
      </Backdrop>
      {children}
    </AuthContext.Provider>
  )
}


interface IAuthContextState { }
