import React from "react";
import { Box, Typography } from "@mui/material";

// The width of the "|" before each item. In pixels: the negative margin and
// the bar must be exactly as wide, and they sit in different font sizes.
const SEPARATOR = "18px";

interface ILinkBarProps {
  items: React.ReactNode[];
  justifyContent?: React.CSSProperties["justifyContent"] | Record<string, React.CSSProperties["justifyContent"]>;
  // False leaves a plain gap where each "|" would be.
  separators?: boolean;
}

// A row of links separated by "|" (the navbar, the account links). It wraps on
// narrow screens, and a line never starts or ends with a "|": every item draws
// its own bar before it, and the bar of the first item on each wrapped line
// falls outside the box, which hides it (the very first one is hidden too).
export const LinkBar: React.FC<ILinkBarProps> = ({ items, justifyContent, separators = true }) => (
  <Box sx={{ overflow: "hidden" }}>
    <Box
      sx={{
        display: "flex",
        flexWrap: "wrap",
        justifyContent,
        ml: `-${SEPARATOR}`,
      }}
    >
      {items.map((item, idx) => (
        <Typography
          key={idx}
          variant={"body2"}
          component={"span"}
          sx={{
            whiteSpace: "nowrap",
            "&::before": {
              content: separators ? '"|"' : '""',
              display: "inline-block",
              width: SEPARATOR,
              textAlign: "center",
            },
            // The first bar would show when the row is right-aligned; hidden,
            // not removed, so it keeps the width the negative margin expects.
            "&:first-of-type::before": { visibility: "hidden" },
          }}
        >
          {item}
        </Typography>
      ))}
    </Box>
  </Box>
);
