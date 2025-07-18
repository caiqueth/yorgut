import { Box, Typography } from "@mui/material"
import { useReplyTopic } from "api/topic/queries";
import { TextCompositor } from "components/molecules/TextCompositor"
import React, { useState } from "react"
import { useParams } from "react-router-dom";

interface IProps {
  topicId: number;
}
export const Reply: React.FC<IProps> = ({ topicId }) => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const [content, setContent] = useState("");
  const { mutate: reply, isPending } = useReplyTopic(topicId);

  const handleSubmit = () => {
    reply({
      communityId: communityId,
      title: "",
      text: content,
    })
  }

  return (
    <Box sx={{ p: 1 }}>
      <TextCompositor
        content={content}
        header={"Responder:"}
        onChangeContent={setContent}
        onSubmit={handleSubmit}
        isPending={isPending}
        maxLength={2048}
      />
    </Box>
  )
}