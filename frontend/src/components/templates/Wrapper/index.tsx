import React from "react";
import { Navbar } from "components/organisms/Navbar";
import { UserBar } from "components/organisms/UserBar";

import { Box, Stack } from "@mui/material";
import { YorgutTitle } from "components/atoms/YorgutTitle";
import { useMe } from "api/me/queries";

const lineColor = "#bac6dd";

interface IWrapperProps {
  children: React.ReactNode;
}
export const Wrapper: React.FC<IWrapperProps> = ({ children }) => {
  return (
    <Stack
      direction={"column"}
      width={1200}
      spacing={2}
      alignItems="center"
      sx={{
        display: "flex",
        alignContent: "center",
        m: "0 auto",
        pt: 3,
        px: 2,
      }}
    >
      <Stack direction={"row"} width={"100%"} alignItems={"center"}>
        <Box
          width={"100%"}
          sx={{
            p: 0, m: 0,
            display: "flex",
            flexDirection: "column",
            borderBottom: `1px solid ${lineColor}`
          }}>
          <UserBar />
          <Navbar />
        </Box>
        <Box sx={{
          display: "block",
          width: 43,
          height: 43,
          background: `linear-gradient(135deg, transparent 49%, ${lineColor} 50%, transparent 51%)`,
        }}>
        </Box>

        <Box sx={{
          borderTop: `1px solid ${lineColor}`
        }}>
          <YorgutTitle size={36} />
        </Box>
      </Stack>
      {children}
    </Stack>
  );
};
