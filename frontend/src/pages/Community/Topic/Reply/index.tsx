import { Box } from "@mui/material"
import { useReplyTopic } from "api/topic/queries";
import { TextCompositor } from "components/molecules/TextCompositor"
import React, { useState } from "react"
import { useParams } from "react-router-dom";
import { MembersOnly } from "../../components/MembersOnly";

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
      title: "",
      text: content,
    })
  }

  return (
    <MembersOnly communityId={communityId} action={"responder tópicos"}>
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
    </MembersOnly>
  )
}
