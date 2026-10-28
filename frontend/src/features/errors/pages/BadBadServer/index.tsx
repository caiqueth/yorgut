import { Paper, Stack, Typography } from "@mui/material";
import React from "react";

export const BadBadServer: React.FC = () => {
  return (
    <Paper sx={{ width: 500, m: 5, p: 3 }}>
      <Stack direction="column" gap={1}>
        <Typography variant={"h6"}>
          Bad, bad server. No donut for you.
        </Typography>
        <Typography>
          Unfortunately, the yorgut.net server has acted out in unexpected way.
          Hopefully, it will return to its helpful self if you try again in a
          few minutes.
        </Typography>
        <Typography>
          It's likely that the server will behave this way on occasion during
          the coming months. We apologize for the inconvenience and for our
          server's lack of consideration for others.
        </Typography>
      </Stack>
    </Paper>
  );
};
