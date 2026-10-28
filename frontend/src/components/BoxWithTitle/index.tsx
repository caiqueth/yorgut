import { Box, LinearProgress, Paper, Typography, useTheme } from "@mui/material";
import { InlineSkeleton } from "components/InlineSkeleton";
import React from "react";

// A count in a box title, "recados (12)", with a skeleton in place of the
// number until it loads (instead of a misleading 0).
export const Count: React.FC<{ value?: number }> = ({ value }) => (
  <>({value ?? <InlineSkeleton width={"2ch"} />})</>
);

interface IBoxWithTitleProps {
  // Text, or text with a <Count />.
  title: React.ReactNode;
  children?: React.ReactNode;
  // Shows a thin progress bar under the title. The bar is always there, only
  // hidden, so showing it doesn't move anything below. Lists pass
  // isFetching && !isLoading: the first load shows skeletons instead, and the
  // bar covers every reload after it (after a change, or a page change).
  showProgress?: boolean;
}
export const BoxWithTitle: React.FC<IBoxWithTitleProps> = ({
  title,
  children,
  showProgress,
}) => {
  const align = useTheme().yorgut.boxTitleAlign;
  return (
    // overflow keeps the title strip and striped rows inside rounded corners.
    <Paper sx={{ overflow: "hidden" }}>
      <Box
        sx={{
          textAlign: align,
          width: "100%",
          boxSizing: "border-box",
          backgroundColor: "yorgut.boxTitle",
          ...(align === "left" && { px: 1.5, pt: 1 }),
        }}
      >
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {title}
        </Typography>
      </Box>
      <LinearProgress sx={{ visibility: showProgress ? "visible" : "hidden" }} />
      {children}
    </Paper>
  );
};
