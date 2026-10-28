import { useIdParam } from "hooks/useIdParam";
import { Box } from "@mui/material"
import { useReplyTopic } from "api/topic/queries";
import { TextCompositor } from "components/TextCompositor"
import React, { useState } from "react"
import { MembersOnly } from "features/community/components/MembersOnly";

interface IProps {
  topicId: number;
  // Called after the reply is saved, e.g. to show the page it landed on.
  onReplied?: () => void;
}
export const Reply: React.FC<IProps> = ({ topicId, onReplied }) => {
  const communityId = useIdParam("communityId");

  const [content, setContent] = useState("");
  const { mutate: reply, isPending } = useReplyTopic(topicId, communityId);

  const handleSubmit = () => {
    reply(
      {
        title: "",
        text: content,
      },
      {
        onSuccess: () => {
          setContent("");
          onReplied?.();
        },
      },
    );
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
