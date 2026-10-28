import React from "react";
import { Box } from "@mui/material";
import Markdown, { Components } from "react-markdown";
import rehypeRaw from "rehype-raw";
import rehypeSanitize, { Options as SanitizeOptions } from "rehype-sanitize";
import remarkBreaks from "remark-breaks";
import remarkGfm from "remark-gfm";

// What user text may use, as markdown or as HTML: emphasis, links, lists,
// quotes and code. Everything else is dropped by rehype-sanitize before React
// renders it: scripts and styles with their content, other tags (images,
// headings, tables, iframes...) keeping only their text, and every attribute
// but a link's http(s)/mailto href. Nothing is ever set as raw HTML: the
// markdown and the HTML in it become a syntax tree, which is cleaned and then
// rendered as React elements.
const schema: SanitizeOptions = {
  tagNames: ["p", "br", "strong", "b", "em", "i", "u", "s", "del", "a", "ul", "ol", "li", "blockquote", "code", "pre"],
  attributes: { a: ["href"] },
  protocols: { href: ["http", "https", "mailto"] },
  strip: ["script", "style"],
};

// Links leave the app in a new tab, without telling the other site where the
// reader came from, and give it no ranking credit.
const components: Components = {
  a: ({ href, children }) => (
    <a href={href} target={"_blank"} rel={"noopener noreferrer nofollow ugc"}>
      {children}
    </a>
  ),
};

interface IFormattedTextProps {
  // Text a user wrote, in markdown with optional simple HTML.
  text: string;
}

// Renders user-written text with its formatting, safely. Single line breaks
// are kept (remark-breaks), as people type scraps like chat messages.
export const FormattedText: React.FC<IFormattedTextProps> = ({ text }) => (
  <Box
    sx={{
      overflowWrap: "anywhere",
      "& p, & ul, & ol, & blockquote, & pre": { my: 0.5 },
      "& > :first-child": { mt: 0 },
      "& > :last-child": { mb: 0 },
      "& ul, & ol": { pl: 3 },
      "& blockquote": { mx: 0, pl: 1, borderLeft: "3px solid", borderColor: "divider" },
      "& pre": { whiteSpace: "pre-wrap" },
    }}
  >
    <Markdown
      remarkPlugins={[remarkGfm, remarkBreaks]}
      rehypePlugins={[rehypeRaw, [rehypeSanitize, schema]]}
      components={components}
    >
      {text}
    </Markdown>
  </Box>
);
