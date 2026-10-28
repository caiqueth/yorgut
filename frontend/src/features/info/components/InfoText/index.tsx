import { Box, Typography } from "@mui/material";
import React from "react";

// The building blocks of the info texts (terms of use, community rules): a
// titled section and a bulleted list.

export const Section: React.FC<{ title: string; children: React.ReactNode }> = ({ title, children }) => (
  <Box>
    <Typography variant={"subtitle2"} fontWeight={"bold"} gutterBottom>
      {title}
    </Typography>
    {children}
  </Box>
);

export const List: React.FC<{ items: React.ReactNode[] }> = ({ items }) => (
  <Box component={"ul"} sx={{ m: 0, pl: 3 }}>
    {items.map((item, idx) => (
      <Typography key={idx} component={"li"} variant={"body2"} sx={{ mb: 0.5 }}>
        {item}
      </Typography>
    ))}
  </Box>
);
