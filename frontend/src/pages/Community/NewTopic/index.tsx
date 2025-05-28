import React, { useState } from "react";

import { TextCompositor } from "components/molecules/TextCompositor";
import { useCreateTopic } from "api/topic/queries";
import { Box, TextField, Typography } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";

interface IProps { }
export const NewTopic: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");

  const [content, setContent] = useState("");
  const [title, setTitle] = useState("");

  const { mutate: createTopic } = useCreateTopic();

  const navigate = useNavigate();

  const handleSubmit = () => {
    createTopic({
      communityId,
      title: title,
      text: content,
    })
    navigate(`/community/${communityId}/forum`)
  }

  return (
    <Box>
      <Typography>Título</Typography>
      <TextField variant={"outlined"} value={title} onChange={(e) => setTitle(e.target.value)} />
      <TextCompositor
        content={content}
        header={"Criar um novo tópico"}
        onChangeContent={setContent}
        onSubmit={handleSubmit}
        maxLength={2048}
      />
    </Box>
  )
}