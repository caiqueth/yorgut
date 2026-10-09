import React from "react";
import { palette } from "theme";

import { Button as MuiButton, ButtonProps } from "@mui/material";

type IButtonProps = ButtonProps;
export const Button: React.FC<IButtonProps> = ({ children, ...props }) => {
  return (
    <MuiButton
      {...props}
      sx={{
        m: 1,
        // backgroundColor: "rgba(195,208,238,1)",
        backgroundColor: palette.mediumBlue,
        boxShadow: `inset 0px 0px 3px 1px rgba(146, 163, 195,1);`,
        borderRadius: 2,
        px: 2,
        color: "black",
        lineHeight: 1,
        textTransform: "none",
        typographyProps: {
          variant: "h1",
        },
        ...props.sx,
      }}
    >
      {children}
    </MuiButton>
  );
};
