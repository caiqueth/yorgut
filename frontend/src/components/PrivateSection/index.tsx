import React from "react";
import { Box, Typography } from "@mui/material";

interface IPrivateSectionProps {
  // Whose profile it is; undefined while it loads.
  firstName?: string;
  // What is kept to friends, with its article: "o álbum", "os recados".
  what: string;
}

// Shown instead of a part of a profile its user keeps to friends.
export const PrivateSection: React.FC<IPrivateSectionProps> = ({ firstName, what }) => (
  <Box sx={{ p: 1 }}>
    <Typography variant={"body2"}>
      🔒 {firstName ?? "Esta pessoa"} mostra {what} só para amigos.
    </Typography>
  </Box>
);
