import { Box, Stack } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Disclaimer } from "components/Disclaimer";
import { TermsOfUse } from "features/info/components/TermsOfUse";
import React from "react";

// Public: the footer of every page links here, logged in or not, so the terms
// can be read before logging in.
export const About: React.FC = () => (
  <Stack direction={"column"} gap={2} sx={{ width: "100%" }}>
    <BoxWithTitle title={"sobre o Yorgut"}>
      <Box sx={{ p: 2 }}>
        <Disclaimer variant={"full"} />
      </Box>
    </BoxWithTitle>
    <BoxWithTitle title={"termos de uso e privacidade"}>
      <Box sx={{ p: 2 }}>
        <TermsOfUse />
      </Box>
    </BoxWithTitle>
  </Stack>
);
