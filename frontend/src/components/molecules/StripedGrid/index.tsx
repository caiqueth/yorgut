import React from "react";
import { isStripedRow } from "helpers/isStripedRow";
import { Link } from "react-router-dom";
import { palette } from "theme";

import { Grid2 as Grid, Stack, Typography } from "@mui/material";

interface IStripedGridProps {
  items: {
    imageSrc: string;
    url: string;
    text: string;
  }[];
  elementsPerRow?: number;
}
export const StripedGrid: React.FC<IStripedGridProps> = ({
  items,
  elementsPerRow = 3,
}) => {
  return (
    <Grid container direction={"row"} width={"100%"}>
      {items?.map(({ imageSrc, url, text }, idx) => (
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
          <Stack
            alignItems={"center"}
            flexDirection={"column"}
            sx={{ textAlign: "center" }}
          >
            <Link to={url}>
              <img height={64} src={imageSrc} alt="image not found" />
              <Typography variant={"body2"}>{text}</Typography>
            </Link>
          </Stack>
        </Grid>
      ))}
    </Grid>
  );
};
