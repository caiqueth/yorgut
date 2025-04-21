import React from "react";
import { isStripedRow } from "helpers/isStripedRow";
import { palette } from "theme";

import { Grid2 as Grid } from "@mui/material";

interface IStripedGridProps {
  items: React.ReactNode[];
  elementsPerRow?: number;
}
export const StripedGrid: React.FC<IStripedGridProps> = ({
  items,
  elementsPerRow = 3,
}) => {
  return (
    <Grid container direction={"row"} width={"100%"}>
      {items.map((item, idx) => (
        <Grid
          size={Math.ceil(12 / elementsPerRow)}
          key={idx}
          sx={{
            p: 1,
            backgroundColor: isStripedRow(idx, elementsPerRow)
              ? palette.darkBlue
              : palette.mediumBlue,
          }}
        >
          {item}
        </Grid>
      ))}
    </Grid>
  );
};
