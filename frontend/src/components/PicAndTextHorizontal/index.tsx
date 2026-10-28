import React from "react";
import { Box, Skeleton, Stack, Typography } from "@mui/material";
import { Link } from "react-router-dom";

interface IPicAndTextHorizontalProps {
  // Null or missing: no picture uploaded.
  imageSrc?: string | null;
  // Where the title links: a profile, a community, a topic.
  url: string;
  // Where the picture links, when it shows someone else than the title, e.g.
  // a topic's author next to the topic title. Defaults to url.
  imageUrl?: string;
  // Only this part is the link: a name or a title.
  title: string;
  // Plain text right after the link, e.g. " em 12 de outubro de 2026:".
  titleSuffix?: string;
  // Plain text keeps its line breaks; anything else (e.g. FormattedText)
  // renders as it is.
  text?: React.ReactNode;
  imageSize?: number;
}
export const PicAndTextHorizontal: React.FC<IPicAndTextHorizontalProps> = ({
  url,
  imageUrl = url,
  title,
  titleSuffix,
  text,
  imageSrc,
  imageSize = 40
}) => {
  return (
    <Stack direction={"row"} gap={1} alignItems={"center"}>
      <Box
        component={Link}
        to={imageUrl}
        sx={{
          display: "block",
          minWidth: imageSize,
          width: imageSize,
          height: imageSize,
          minHeight: imageSize,
          backgroundImage: imageSrc ? `url("${imageSrc}")` : undefined,
          backgroundSize: "cover",
          backgroundRepeat: "no-repeat",
          alignSelf: "start",
        }}
      />
      <Stack direction={"column"}>
        <Typography><Link to={url}>{title}</Link>{titleSuffix}</Typography>
        {text && (typeof text === "string"
          ? <Typography sx={{ whiteSpace: "pre-line" }}>{text}</Typography>
          // A div: formatted text has paragraphs and lists, which can't go in a <p>.
          : <Typography component={"div"}>{text}</Typography>)}
      </Stack>
    </Stack>
  );
};

interface IPicAndTextHorizontalSkeletonProps {
  imageSize?: number;
  // Widths of the text lines, in pixels: the first is the title, the second
  // (if any) the text under it.
  lineWidths?: number[];
}

// The loading state of a PicAndTextHorizontal: the same layout, with the
// picture and each line replaced by a skeleton of the same height.
export const PicAndTextHorizontalSkeleton: React.FC<IPicAndTextHorizontalSkeletonProps> = ({
  imageSize = 40,
  lineWidths = [120, 180],
}) => (
  <Stack direction={"row"} gap={1} alignItems={"center"}>
    <Skeleton
      variant={"rectangular"}
      width={imageSize}
      height={imageSize}
      sx={{ minWidth: imageSize, alignSelf: "start" }}
    />
    <Stack direction={"column"}>
      {lineWidths.map((width, idx) => (
        <Typography key={idx}><Skeleton variant={"text"} width={width} /></Typography>
      ))}
    </Stack>
  </Stack>
);
