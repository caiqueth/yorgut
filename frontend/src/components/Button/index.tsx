import React from "react";

import { Button as MuiButton, ButtonProps } from "@mui/material";

type IButtonProps = ButtonProps;
export const Button: React.FC<IButtonProps> = ({ children, ...props }) => {
  return (
    <MuiButton
      {...props}
      sx={{
        ...props.sx,
        boxShadow: "inset 1px 1px 8px 0px rgba(255,255,255,1);",
        borderRadius: 2,
      }}
    >
      {children}
    </MuiButton>
  );
};
