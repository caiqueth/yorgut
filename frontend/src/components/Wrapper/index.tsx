import React from "react";
import { Navbar } from "components/Navbar";
import { UserBar } from "components/UserBar";

import { Box, Stack, useTheme } from "@mui/material";
import { YorgutTitle } from "components/YorgutTitle";
import { Disclaimer } from "components/Disclaimer";
import { PHONE_BAR_HEIGHT } from "components/Navbar";

interface IWrapperProps {
  children: React.ReactNode;
}
export const Wrapper: React.FC<IWrapperProps> = ({ children }) => {
  const theme = useTheme();
  const { border: lineColor, navbar, navbarEdge } = theme.palette.yorgut;
  // The navbar holds the logo instead (Navbar): no diagonal and title here.
  const logoInNavbar = theme.yorgut.logoInNavbar;
  // Or the user's links (Navbar).
  const userBarInNavbar = theme.yorgut.userBarInNavbar;
  return (
    <Stack
      direction={"column"}
      spacing={2}
      alignItems="center"
      sx={{
        display: "flex",
        alignContent: "center",
        boxSizing: "border-box",
        width: "100%",
        maxWidth: 1200,
        m: "0 auto",
        // On phones a header filled here (below) touches the top too.
        pt: theme.yorgut.headerAtTop ? 0 : { xs: logoInNavbar ? 1 : 0, sm: 3 },
        px: { xs: 1, sm: 2 },
      }}
    >
      {/* On phones the title goes at the left of the navbar, at the top so it
          stays put when the navbar's menu opens.
          From sm up, one line runs under the links, up the diagonal and over
          the title. The three parts stretch to the row's height and the
          diagonal goes corner to corner, so its ends always meet the other
          two lines, however tall the links are. */}
      <Stack
        direction={{ xs: "row-reverse", sm: "row" }}
        alignItems={"stretch"}
        sx={{
          width: "100%",
          // On phones the navbar's fill runs under the title too (Navbar
          // leaves its own out), out to both edges of the screen over the
          // page's side padding.
          ...(!logoInNavbar && {
            background: { xs: navbar, sm: "none" },
            boxShadow: { xs: `inset 0 1px ${navbarEdge}, inset 0 -1px ${navbarEdge}`, sm: "none" },
            boxSizing: "border-box",
            width: { xs: `calc(100% + ${theme.spacing(2)})`, sm: "100%" },
            mx: { xs: -1, sm: 0 },
            px: { xs: 1, sm: 0 },
            py: { xs: 0.5, sm: 0 },
          }),
        }}
      >
        <Box
          width={"100%"}
          sx={{
            p: 0, m: 0,
            display: "flex",
            flexDirection: "column",
            justifyContent: "flex-end",
            // On phones the header's fill (above) has its own edges.
            ...(!logoInNavbar && { borderBottom: { sm: `1px solid ${lineColor}` } }),
          }}>
          {/* On phones the navbar's "Mais" menu holds the user's links. */}
          {!userBarInNavbar && <Box sx={{ display: { xs: "none", sm: "block" } }}><UserBar /></Box>}
          <Navbar />
        </Box>
        {!logoInNavbar && (
          <>
            <Box sx={{
              display: { xs: "none", sm: "block" },
              flexShrink: 0,
              width: 43,
              minHeight: 43,
              // "to bottom right" puts the 50% line through the other two
              // corners: from bottom left (the links' line) to top right (the
              // title's).
              background: `linear-gradient(to bottom right, transparent calc(50% - 1px), ${lineColor} calc(50% - 0.5px), ${lineColor} calc(50% + 0.5px), transparent calc(50% + 1px))`,
            }}>
            </Box>

            <Box sx={{
              borderTop: { sm: `1px solid ${lineColor}` },
              flexShrink: 0,
              // As tall as the navbar's phone row (Navbar), the logo centered.
              alignSelf: { xs: "flex-start", sm: "stretch" },
              height: { xs: PHONE_BAR_HEIGHT, sm: "auto" },
              display: { xs: "flex", sm: "block" },
              alignItems: "center",
              pr: { xs: 1, sm: 0 },
            }}>
              <YorgutTitle size={{ xs: 22, sm: 36 }} />
            </Box>
          </>
        )}
      </Stack>
      {children}
      <Box component={"footer"} sx={{ width: "100%", pt: 2, pb: 3, borderTop: `1px solid ${lineColor}`, textAlign: "center" }}>
        <Disclaimer variant={"short"} />
      </Box>
    </Stack>
  );
};
