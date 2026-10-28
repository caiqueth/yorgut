import { Box } from "@mui/material";
import React, { useId } from "react"

interface IProps {
  element: React.ReactNode | string;
  times: number;
  // How far each copy slides over the previous one, in em so it follows the
  // font size.
  overlap?: number;
  // The color behind the row (#rrggbb), which each copy's cutout is painted
  // with.
  background?: string;
  // Shown by the browser on hover (title), e.g. "90%".
  title?: string;
  // Read by screen readers instead of the icons, e.g. "sexy: 90%".
  label?: string;
}

// "#ffffff" → "1 1 1", the 0–1 channels an feColorMatrix row takes.
const channels = (hex: string) =>
  [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255);

// The element repeated in a row, each copy slightly stacked over the previous
// one. Emoji fonts (Windows' especially) draw see-through parts, so each copy
// after the first sits on a solid cutout of its own shape, in the background
// color, which hides the copy under it.
export const Repeater: React.FC<IProps> = ({
  element,
  times,
  overlap = 0.85,
  background = "#ffffff",
  title,
  label,
}) => {
  // useId gives ":r0:"; keep it to characters a url(#...) reference takes.
  const filterId = `cutout${useId().replace(/[^a-zA-Z0-9]/g, "")}`;
  const [r, g, b] = channels(background);

  return (
    <Box sx={{ display: "flex" }} title={title} role={label ? "img" : undefined} aria-label={label}>
      {/* Paints a shape in the background color and makes any pixel it covers
          at all fully opaque: a solid silhouette. */}
      <svg width={0} height={0} style={{ position: "absolute" }} aria-hidden>
        <filter id={filterId}>
          <feColorMatrix type={"matrix"} values={`0 0 0 0 ${r}  0 0 0 0 ${g}  0 0 0 0 ${b}  0 0 0 255 0`} />
        </filter>
      </svg>
      {Array.from({ length: times }, (_, idx) => (
        <Box key={idx} sx={{ position: "relative", marginLeft: idx !== 0 ? `-${overlap}em` : 0 }}>
          {idx !== 0 && (
            <Box aria-hidden sx={{ position: "absolute", inset: 0, filter: `url(#${filterId})` }}>
              {element}
            </Box>
          )}
          <Box sx={{ position: "relative" }}>{element}</Box>
        </Box>
      ))}
    </Box>
  )
}
