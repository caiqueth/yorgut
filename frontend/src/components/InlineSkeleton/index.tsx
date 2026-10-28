import React from "react";
import { Skeleton, SxProps, Theme } from "@mui/material";

interface IInlineSkeletonProps {
  // Roughly the width of the text it stands for, e.g. "6ch" for a first name.
  width: string | number;
}

// A word or number that hasn't loaded, inside a sentence. It covers where the
// letters will be: from just under the baseline to about the top of capital
// letters (0.8em tall, 0.1em under the baseline), so it lines up with the
// words around it and never makes the line taller. (MUI's "text" variant
// instead squashes a line-high box, which floats above the letters.)
export const InlineSkeleton: React.FC<IInlineSkeletonProps> = ({ width }) => (
  <Skeleton
    variant={"rounded"}
    width={width}
    height={"0.8em"}
    sx={{ display: "inline-block", verticalAlign: "-0.1em" }}
  />
);

// For a heading whose height depends on its letters, like the drop-cap names
// of the profile and community headers: while loading, the heading gets
// loadingTextSx and renders <TextPlaceholder />, i.e. sample text in a
// transparent color with a skeleton over it, so the line has exactly the
// height it will have once the name arrives.
export const loadingTextSx: SxProps<Theme> = { position: "relative", color: "transparent" };

export const TextPlaceholder: React.FC<{ sample?: string }> = ({ sample = "Carregando" }) => (
  <>
    {sample}
    <Skeleton variant={"rounded"} sx={{ position: "absolute", inset: 0, height: "auto" }} />
  </>
);
