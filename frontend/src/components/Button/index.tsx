import React from "react";

import { Button as MuiButton, ButtonProps, useTheme } from "@mui/material";

type IButtonProps = ButtonProps;
// The app's button. Its look comes from the theme: colors from
// palette.yorgut.button*, size and text from yorgut.buttonShape.
export const Button: React.FC<IButtonProps> = ({ children, ...props }) => {
  const theme = useTheme();
  const { button, buttonShadow, buttonBorder, buttonText } = theme.palette.yorgut;
  const shape = theme.yorgut.buttonShape;

  return (
    <MuiButton
      {...props}
      sx={{
        m: 1,
        // background, not backgroundColor: the theme's fill may be a gradient.
        // A caller that changes it sets background too.
        background: button,
        // The outline is a shadow ring, not a border, so it adds no size.
        boxShadow: `inset 0px 0px 3px 1px ${buttonShadow}, 0 0 0 1px ${buttonBorder}`,
        borderRadius: 2,
        color: buttonText,
        lineHeight: 1,
        textTransform: "none",
        ...(shape ? { ...shape, minWidth: 0 } : { px: 2 }),
        ...props.sx,
      }}
    >
      {children}
    </MuiButton>
  );
};
