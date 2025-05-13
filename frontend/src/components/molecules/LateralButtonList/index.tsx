import React from "react";
import { palette } from "theme";

import { Button, Stack, useTheme } from "@mui/material";

interface ILateralButtonList {
  items: {
    text: string;
    avatar: React.ReactNode;
    onClick?: () => void;
  }[];
}
export const LateralButtonList: React.FC<ILateralButtonList> = ({ items }) => {
  const theme = useTheme();
  return (
    <Stack spacing={"2px"}>
      {items.map(({ text, avatar, onClick }, idx) => (
        <Button
          disableElevation
          key={idx}
          fullWidth
          variant="contained"
          onClick={onClick}
          startIcon={avatar}
          disabled={!onClick}
          sx={{
            backgroundColor: palette.lightBlue,
            borderRadius: 0,
            color: theme.palette.text.primary,
            fontSize: "12px",
            justifyContent: "left",
            lineHeight: 1,
            overflow: "hidden",
            textOverflow: "ellipsis",
            textTransform: "none",
            whiteSpace: "nowrap",
          }}
        >
          {text}
        </Button>
      ))}
    </Stack>
  );
};
