import React, { useState } from "react";

import { TextCompositor } from "components/molecules/TextCompositor";

interface IProps { }
export const NewTopic: React.FC<IProps> = () => {
  const [content, setContent] = useState("");
  const handleSubmit = () => {

  }

  return (
    <TextCompositor
      content={content}
      header={"Criar um novo tópico"}
      onChangeContent={setContent}
      onSubmit={handleSubmit}
      maxLength={2048}
    />
  )
}