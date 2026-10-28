import { Box } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { TermsOfUse } from "features/info/components/TermsOfUse";
import React from "react";

// The terms of use tab of the help.
export const HelpTerms: React.FC = () => (
  <BoxWithTitle title={"termos de uso e privacidade"}>
    <Box sx={{ p: 2 }}>
      <TermsOfUse />
    </Box>
  </BoxWithTitle>
);
