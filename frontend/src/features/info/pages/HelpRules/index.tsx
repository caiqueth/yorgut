import { Box } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { CommunityRules } from "features/info/components/CommunityRules";
import React from "react";

// The rules tab of the help.
export const HelpRules: React.FC = () => (
  <BoxWithTitle title={"regras de convivência"}>
    <Box sx={{ p: 2 }}>
      <CommunityRules />
    </Box>
  </BoxWithTitle>
);
