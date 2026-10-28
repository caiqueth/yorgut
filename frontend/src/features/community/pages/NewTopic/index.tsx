import { useIdParam } from "hooks/useIdParam";
import React, { useState } from "react";

import { TextCompositor } from "components/TextCompositor";
import { useCreateTopic } from "api/topic/queries";
import { Box, TextField, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { MembersOnly } from "features/community/components/MembersOnly";

export const NewTopic: React.FC = () => {
  const communityId = useIdParam("communityId");

  const [content, setContent] = useState("");
  const [title, setTitle] = useState("");

  const { mutate: createTopic, isPending, error } = useCreateTopic();

  const navigate = useNavigate();

  const handleSubmit = () => {
    createTopic(
      {
        communityId,
        title: title,
        text: content,
      },
      { onSuccess: () => navigate(`/community/${communityId}/forum`) },
    );
  }

  return (
    <MembersOnly communityId={communityId} action={"criar tópicos"}>
      <Box>
        <Typography>Título</Typography>
        <TextField
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          slotProps={{ htmlInput: { maxLength: 200 } }}
        />
        <TextCompositor
          content={content}
          header={"Criar um novo tópico"}
          onChangeContent={setContent}
          onSubmit={handleSubmit}
          isPending={isPending}
          maxLength={2048}
          canSubmit={title.trim() !== ""}
        />
        {error && (
          <Typography variant={"body2"} color={"error"}>
            {error.message}
          </Typography>
        )}
      </Box>
    </MembersOnly>
  )
}
