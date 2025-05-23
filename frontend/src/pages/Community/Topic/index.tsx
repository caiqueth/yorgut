import { Box, Stack, Typography } from "@mui/material";
import { useTopicReplies } from "api/topic/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Link, useParams } from "react-router-dom";

interface IProps { }
export const Topic: React.FC<IProps> = () => {
  const params = useParams();
  const topicId = parseInt(params.topicId ?? "");

  const { data: replies } = useTopicReplies(topicId);

  const items = replies ?
    [replies.topic, ...(replies.replies ?? [])].map((r) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={`/${r.authorProfilePicture}`} height={40} />
        <Stack direction={"column"}>
          <Typography><Link to={`/profile/${r.authorId}`}>{r.authorFirstName}</Link> em {dayjsFromDate(r.createdAt, "D [de] MMMM [de] YYYY")}:</Typography>
          <Typography>{r.text}</Typography>
        </Stack>
      </Stack>
    )) : [];

  return (
    <Box>
      <BoxWithTitle title="respostas">
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  )
}