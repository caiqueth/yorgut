import React from "react";
import { Box, Grid2 as Grid, Stack } from "@mui/material";

interface IPageWithSidebarProps {
  // The big picture (user or community) and the menu under it.
  picture: React.ReactNode;
  menu: React.ReactNode;
  children: React.ReactNode;
}

// The frame of the home, profile, community and scrapbook pages: picture and
// menu on the left, the page on the right. Below the md breakpoint (phones,
// tablets) the column moves on top, with the picture and the menu side by side.
export const PageWithSidebar: React.FC<IPageWithSidebarProps> = ({ picture, menu, children }) => (
  <Box sx={{ width: "100%" }}>
    <Grid container spacing={2}>
      <Grid size={{ xs: 12, md: 2 }}>
        <Stack direction={{ xs: "row", md: "column" }} gap={2} alignItems={{ xs: "flex-start", md: "stretch" }}>
          <Box sx={{ width: { xs: 140, md: "100%" }, flexShrink: 0 }}>{picture}</Box>
          <Box sx={{ flex: 1, minWidth: 0 }}>{menu}</Box>
        </Stack>
      </Grid>
      <Grid size={{ xs: 12, md: 10 }}>{children}</Grid>
    </Grid>
  </Box>
);
