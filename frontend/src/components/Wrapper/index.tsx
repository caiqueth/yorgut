import React from "react";
import { Navbar } from "components/Navbar";
import { UserBar } from "components/UserBar";

import { Stack } from "@mui/material";

interface IWrapperProps {
  children: React.ReactNode;
}
export const Wrapper: React.FC<IWrapperProps> = ({ children }) => {
  return (
    <Stack
      direction={"column"}
      width={1200}
      alignItems="center"
      sx={{
        display: "flex",
        alignContent: "center",
        m: "0 auto",
        pt: 3,
      }}
    >
      <UserBar />
      <Navbar />
      {children}
    </Stack>
  );
};
