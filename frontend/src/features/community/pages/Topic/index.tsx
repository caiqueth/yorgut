import { Box, Typography } from "@mui/material";
import { useTopicReplies } from "api/topic/queries";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { useParams } from "react-router-dom";
import { Reply } from "features/community/components/TopicReply";

interface IProps { }
export const Topic: React.FC<IProps> = () => {
  const params = useParams();
  const topicId = parseInt(params.topicId ?? "");

  const { data: replies, isLoading } = useTopicReplies(topicId);

  const items = replies ?
    [replies.topic, ...(replies.replies ?? [])].map((r) => (
      <PicAndTextHorizontal
        title={`${r.authorFirstName} em ${dayjsFromDate(r.createdAt, "D [de] MMMM [de] YYYY")}:`}
        text={r.text}
        url={`/profile/${r.authorId}`}
        imageSrc={r.authorProfilePicture}
        imageSize={56}
      />
    )) : [];

  return (
    <Box>
      <BoxWithTitle title="respostas">
        <Box sx={{ px: 1 }}>
          <Typography variant={"h6"}>{replies?.topic.title}</Typography>
        </Box>
        <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
        <Reply topicId={topicId} />
      </BoxWithTitle>
    </Box>
  )
}