import { Box } from "@mui/material";
import { useTopicReplies } from "api/topic/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { useParams } from "react-router-dom";
import { Reply } from "./Reply";

interface IProps { }
export const Topic: React.FC<IProps> = () => {
  const params = useParams();
  const topicId = parseInt(params.topicId ?? "");

  const { data: replies } = useTopicReplies(topicId);

  const items = replies ?
    [replies.topic, ...(replies.replies ?? [])].map((r) => (
      <PicAndTextHorizontal
        title={`${r.authorFirstName} em ${dayjsFromDate(r.createdAt, "D [de] MMMM [de] YYYY")}:`}
        text={r.text}
        url={`/profile/${r.authorId}`}
        imageSrc={`/${r.authorProfilePicture}`}
        imageSize={56}
      />
    )) : [];

  return (
    <Box>
      <BoxWithTitle title="respostas">
        <StripedGrid items={items} elementsPerRow={1} />
        <Reply topicId={topicId} />
      </BoxWithTitle>
    </Box>
  )
}