import React, { useState } from "react";

import { TextCompositor } from "components/TextCompositor";
import { useCreateTopic } from "api/topic/queries";
import { Box, TextField, Typography } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import { MembersOnly } from "features/community/components/MembersOnly";

interface IProps { }
export const NewTopic: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const [content, setContent] = useState("");
  const [title, setTitle] = useState("");

  const { mutate: createTopic, isPending, isError } = useCreateTopic();

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
        <TextField variant={"outlined"} value={title} onChange={(e) => setTitle(e.target.value)} />
        <TextCompositor
          content={content}
          header={"Criar um novo tópico"}
          onChangeContent={setContent}
          onSubmit={handleSubmit}
          isPending={isPending}
          maxLength={2048}
        />
        {isError && (
          <Typography variant={"body2"} color={"error"}>
            Não foi possível criar o tópico. Confira o título e o texto.
          </Typography>
        )}
      </Box>
    </MembersOnly>
  )
}
