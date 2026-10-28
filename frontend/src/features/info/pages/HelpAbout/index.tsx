import { Box } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Disclaimer } from "components/Disclaimer";
import React from "react";

// The "sobre" tab of the help.
export const HelpAbout: React.FC = () => (
  <BoxWithTitle title={"sobre o Yorgut"}>
    <Box sx={{ p: 2 }}>
      <Disclaimer variant={"full"} />
    </Box>
  </BoxWithTitle>
);
