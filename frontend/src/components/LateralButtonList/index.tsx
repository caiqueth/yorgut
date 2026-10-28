import React from "react";

import { Button, Skeleton, Stack, useTheme } from "@mui/material";

export interface ILateralButtonItem {
  text: string;
  avatar: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  // While an action started by this item runs: disabled, with a spinner in
  // place of the avatar.
  loading?: boolean;
}

// Holds the place of an item that depends on data still loading (e.g. "adicionar"
// on a profile, until the friendship status arrives), so the menu doesn't grow
// a line at the top when it loads. Same height as any item.
export const placeholderItem: ILateralButtonItem = {
  text: "",
  avatar: <Skeleton variant={"circular"} width={12} height={12} />,
};

interface ILateralButtonList {
  items: ILateralButtonItem[];
}
export const LateralButtonList: React.FC<ILateralButtonList> = ({ items }) => {
  const theme = useTheme();
  return (
    // The gaps show the divider color; the ring around the items is a shadow,
    // so it adds no size.
    <Stack
      spacing={`${theme.yorgut.menuGap}px`}
      sx={{
        backgroundColor: theme.palette.yorgut.menuDivider,
        boxShadow: `0 0 0 1px ${theme.palette.yorgut.menuDivider}`,
      }}
    >
      {items.map(({ text, avatar, onClick, disabled, loading }, idx) => (
        <Button
          disableElevation
          key={idx}
          fullWidth
          variant="contained"
          onClick={onClick}
          startIcon={avatar}
          loading={loading}
          loadingPosition={"start"}
          disabled={!onClick || disabled}
          sx={{
            backgroundColor: theme.palette.yorgut.highlight,
            borderRadius: 0,
            color: theme.palette.yorgut.menuText,
            fontSize: "12px",
            justifyContent: "left",
            lineHeight: 1,
            overflow: "hidden",
            textOverflow: "ellipsis",
            textTransform: "none",
            whiteSpace: "nowrap",
          }}
        >
          {text || <Skeleton variant={"text"} width={"70%"} />}
        </Button>
      ))}
    </Stack>
  );
};
